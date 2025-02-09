package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func startHttpServer() {
	ctr := di.BuildServiceContainer()
	var httpServer *httpservice.HttpServer
	ctr.Resolve(&httpServer)
	httpServer.InitServer()
}

func startCatalogIndexer() {
	slog.Debug("Starting catalog indexer")
	ctr := di.BuildIndexerContainer()

	// update catalogs table
	var catalogRegister *indexer.CatalogRegister
	ctr.Resolve(&catalogRegister)
	catalogRegister.RegisterCatalog()

	// initialize writer
	var writer indexer.Writer[repository.ParquetMastercat]
	ctr.Resolve(&writer)
	writer.Start()

	// initialize indexer
	var catalogIndexer *indexer.Indexer
	ctr.Resolve(&catalogIndexer)
	catalogIndexer.Start()

	// initialize reader
	var reader indexer.Reader
	ctr.Resolve(&reader)
	reader.Start()

	writer.Done()
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
