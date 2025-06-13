// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	metadata_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	partition_reader "github.com/dirodriguezm/xmatch/service/internal/preprocessor/reader"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/reducer"
	partition_writer "github.com/dirodriguezm/xmatch/service/internal/preprocessor/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/web"
	"github.com/gin-gonic/gin"
)

// @title			CrossWave HTTP API
// @version		1.0
// @description	API for the CrossWave Xmatch service. This service allows to search for objects in a given region and to retrieve metadata from the catalogs.
// @host			localhost:8080
// @BasePath		/v1
// @contact.name	Diego Rodriguez Mancini
// @contact.email	diegorodriguezmancini@gmail.com
func startHttpServer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	ctr := di.BuildServiceContainer(ctx, getenv, stdout)
	var api *api.API
	var web *web.Web
	ctr.Resolve(&api)
	ctr.Resolve(&web)

	r := gin.New()
	r.Use(gin.Recovery())
	if getenv("USE_LOGGER") != "" {
		r.Use(func(c *gin.Context) {
			slog.Info("request", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.Next()
		})
	}
	r.SetTrustedProxies([]string{"localhost"})

	api.SetupRoutes(r)
	web.SetupRoutes(r)

	err := r.Run()
	return err
}

func startCatalogIndexer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	slog.Debug("Starting catalog indexer")
	ctr := di.BuildIndexerContainer(ctx, getenv, stdout)

	var cfg *config.Config
	err := ctr.Resolve(&cfg)
	if err != nil {
		return err
	}

	// update catalogs table
	var catalogRegister *indexer.CatalogRegister
	err = ctr.Resolve(&catalogRegister)
	if err != nil {
		return err
	}
	catalogRegister.RegisterCatalog()

	// initialize w
	var w writer.Writer[any]
	err = ctr.NamedResolve(&w, "indexer_writer")
	if err != nil {
		return err
	}

	w.Start()

	// initialize metadata writer
	if cfg.CatalogIndexer.MetadataWriter != nil && cfg.CatalogIndexer.Source.Metadata {
		var metadataWriter writer.Writer[any]
		err := ctr.NamedResolve(&metadataWriter, "metadata_writer")
		if err != nil {
			return err
		}
		metadataWriter.Start()
	}

	// initialize indexer
	var catalogIndexer *mastercat_indexer.IndexerActor
	err = ctr.Resolve(&catalogIndexer)
	if err != nil {
		return err
	}
	catalogIndexer.Start()

	// initialize metadata indexer
	if cfg.CatalogIndexer.Source.Metadata && cfg.CatalogIndexer.MetadataWriter != nil {
		var actor *metadata_indexer.IndexerActor
		err := ctr.Resolve(&actor)
		if err != nil {
			return err
		}
		actor.Start()
	}

	// initialize reader
	var reader reader.Reader
	ctr.Resolve(&reader)
	reader.Start()

	w.Done()
	return nil
}

func startPreprocessor(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	ctr := di.BuildPreprocessorContainer(ctx, getenv, stdout)

	// initialize partition_writer
	var partition_w *partition_writer.PartitionWriter
	err := ctr.Resolve(&partition_w)
	if err != nil {
		return err
	}
	partition_w.Start()

	// initialize source reader
	var reader reader.Reader
	err = ctr.Resolve(&reader)
	if err != nil {
		return err
	}
	reader.Start()

	var cfg *config.Config
	err = ctr.Resolve(&cfg)
	if err != nil {
		return err
	}

	readerResults := reader.GetOutbox()
	go func() {
		defer close(partition_w.InboxChannel)
		for i := range readerResults {
			for msg := range readerResults[i] {
				wMsg := writer.WriterInput[repository.InputSchema]{
					Error: msg.Error,
					Rows:  msg.Rows,
				}

				partition_w.InboxChannel <- wMsg
			}
		}
	}()
	partition_w.Done()

	// Now the partition reader part
	var reducer *reducer.Reducer
	err = ctr.Resolve(&reducer)
	if err != nil {
		return err
	}
	reducer.Start()

	var partition_r *partition_reader.PartitionReader
	err = ctr.Resolve(&partition_r)
	if err != nil {
		return err
	}
	partition_r.Start()

	partition_r.Done()
	reducer.Done()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
	args []string,
) error {
	slog.Info("Starting xmatch service")

	if len(args) < 2 {
		panic("run: Missing arguments")
	}
	fs := flag.NewFlagSet("xmatch", flag.ContinueOnError)

	var profile bool
	fs.BoolVar(&profile, "profile", false, "Enable profiling")

	if err := fs.Parse(args); err != nil {
		return err
	}

	command := args[1]

	if profile {
		slog.Info("Profiling enabled")
		cpuFile, _ := os.Create("cpu.prof")
		defer cpuFile.Close()
		pprof.StartCPUProfile(cpuFile)
		defer pprof.StopCPUProfile()
	}

	switch command {
	case "server":
		return startHttpServer(ctx, getenv, stdout)
	case "indexer":
		return startCatalogIndexer(ctx, getenv, stdout)
	case "preprocessor":
		return startPreprocessor(ctx, getenv, stdout)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
	return nil
}
