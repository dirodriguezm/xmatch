package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/di"
)

func StartCatalogIndexer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	slog.Info("Starting catalog indexer")
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

	// initialize mastercatWriter
	var mastercatWriter *actor.Actor
	err = ctr.NamedResolve(&mastercatWriter, "mastercat_writer")
	if err != nil {
		return err
	}
	mastercatWriter.Start()

	// initialize metadata writer
	var metadataWriter *actor.Actor
	if cfg.CatalogIndexer.Source.Metadata {
		err := ctr.NamedResolve(&metadataWriter, "metadata_writer")
		if err != nil {
			return err
		}
		metadataWriter.Start()
	}

	// initialize indexer
	var mastercatIndexer *actor.Actor
	err = ctr.NamedResolve(&mastercatIndexer, "mastercat_indexer")
	if err != nil {
		return err
	}
	mastercatIndexer.Start()

	// initialize metadata indexer
	var metadataIndexer *actor.Actor
	if cfg.CatalogIndexer.Source.Metadata {
		err := ctr.NamedResolve(&metadataIndexer, "metadata_indexer")
		if err != nil {
			return err
		}
		metadataIndexer.Start()
	}

	// initialize reader
	var reader *reader.SourceReader
	err = ctr.Resolve(&reader)
	if err != nil {
		return fmt.Errorf("Could not resolve reader: %w", err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			panic(fmt.Errorf("Error closing reader: %w", err))
		}
	}()

	reader.Read()
	mastercatIndexer.Stop()
	mastercatWriter.Stop()
	if cfg.CatalogIndexer.Source.Metadata {
		metadataIndexer.Stop()
		metadataWriter.Stop()
	}

	slog.Info("Catalog indexer finished successfully")
	return nil
}
