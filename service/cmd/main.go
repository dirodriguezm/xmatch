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
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	metadata_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	partition_reader "github.com/dirodriguezm/xmatch/service/internal/preprocessor/reader"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/reducer"
	partition_writer "github.com/dirodriguezm/xmatch/service/internal/preprocessor/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// @title			CrossWave HTTP API
// @version		1.0
// @description	API for the CrossWave Xmatch service. This service allows to search for objects in a given region and to retrieve metadata from the catalogs.
// @host			localhost:8080
// @BasePath		/v1
// @contact.name	Diego Rodriguez Mancini
// @contact.email	diegorodriguezmancini@gmail.com
func startHttpServer() {
	ctr := di.BuildServiceContainer()
	var httpServer *httpservice.HttpServer
	ctr.Resolve(&httpServer)
	httpServer.InitServer()
}

func startCatalogIndexer() {
	slog.Debug("Starting catalog indexer")
	ctr := di.BuildIndexerContainer()

	var cfg *config.Config
	ctr.Resolve(&cfg)

	// update catalogs table
	var catalogRegister *indexer.CatalogRegister
	ctr.Resolve(&catalogRegister)
	catalogRegister.RegisterCatalog()

	// initialize w
	var w writer.Writer[any]
	ctr.NamedResolve(&w, "indexer_writer")
	w.Start()

	// initialize metadata writer
	if cfg.CatalogIndexer.MetadataWriter != nil && cfg.CatalogIndexer.Source.Metadata {
		var metadataWriter writer.Writer[any]
		ctr.NamedResolve(&metadataWriter, "metadata_writer")
		metadataWriter.Start()
	}

	// initialize indexer
	var catalogIndexer *mastercat_indexer.IndexerActor
	ctr.Resolve(&catalogIndexer)
	catalogIndexer.Start()

	// initialize metadata indexer
	if cfg.CatalogIndexer.Source.Metadata && cfg.CatalogIndexer.MetadataWriter != nil {
		var actor *metadata_indexer.IndexerActor
		ctr.Resolve(&actor)
		actor.Start()
	}

	// initialize reader
	var reader reader.Reader
	ctr.Resolve(&reader)
	reader.Start()

	w.Done()
}

func startPreprocessor() {
	ctr := di.BuildPreprocessorContainer()

	// initialize partition_writer
	var partition_w *partition_writer.PartitionWriter
	ctr.Resolve(&partition_w)
	partition_w.Start()

	// initialize source reader
	var reader reader.Reader
	ctr.Resolve(&reader)
	reader.Start()

	var cfg *config.Config
	ctr.Resolve(&cfg)

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
	ctr.Resolve(&reducer)
	reducer.Start()

	var partition_r *partition_reader.PartitionReader
	ctr.Resolve(&partition_r)
	partition_r.Start()

	partition_r.Done()
	reducer.Done()
}

func main() {
	slog.Info("Starting xmatch service")
	profile := flag.CommandLine.Bool("profile", false, "Enable profiling")
	flag.Parse()
	command := flag.Arg(0)

	if *profile {
		slog.Info("Profiling enabled")
		cpuFile, _ := os.Create("cpu.prof")
		defer cpuFile.Close()
		pprof.StartCPUProfile(cpuFile)
		defer pprof.StopCPUProfile()
	}

	switch command {
	case "server":
		startHttpServer()
	case "indexer":
		startCatalogIndexer()
	case "preprocessor":
		startPreprocessor()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
