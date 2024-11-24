package main

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
)

func startHttpServer() {
	ctr := di.BuildServiceContainer()
	var httpServer httpservice.HttpServer
	ctr.Resolve(&httpServer)
	httpServer.InitServer()
}

func startCatalogIndexer() {
	slog.Debug("Starting catalog indexer")
	ctr := di.BuildIndexerContainer()

	// initialize writer
	var writer indexer.Writer
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
	startCatalogIndexer()
}
