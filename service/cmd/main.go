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
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
