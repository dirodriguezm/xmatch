package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/app"
)

func StartCatalogIndexer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	slog.Info("Starting catalog indexer")

	cfg, err := app.Config(getenv)
	if err != nil {
		return err
	}

	// database
	repo, err := app.Repository(cfg)
	if err != nil {
		return err
	}

	// update catalogs table
	catalogRegister := app.CatalogRegister(ctx, repo, cfg.CatalogIndexer.Source)
	catalogRegister.RegisterCatalog()

	src, err := app.Source(cfg.CatalogIndexer.Source)
	if err != nil {
		return err
	}

	// initialize mastercatWriter
	mastercatWriter, err := app.MastercatWriter(ctx, cfg, repo, src)
	if err != nil {
		return err
	}
	mastercatWriter.Start()

	// initialize metadata writer
	metadataWriter, err := app.MetadataWriter(ctx, cfg, repo, src)
	if cfg.CatalogIndexer.Source.Metadata {
		metadataWriter.Start()
	}

	// initialize indexer
	mastercatIndexer, err := app.MastercatIndexer(cfg.CatalogIndexer, mastercatWriter)
	if err != nil {
		return err
	}
	mastercatIndexer.Start()

	// initialize metadata indexer
	var metadataIndexer *actor.Actor
	if cfg.CatalogIndexer.Source.Metadata {
		metadataIndexer = app.MetadataIndexer(cfg.CatalogIndexer, metadataWriter)
		metadataIndexer.Start()
	}

	// initialize reader
	sourceReader, err := app.Reader(src, cfg.CatalogIndexer.Reader, cfg.CatalogIndexer.Source, mastercatIndexer, metadataIndexer)
	defer func() error {
		err := sourceReader.Close()
		if err != nil {
			return fmt.Errorf("Error closing reader: %w", err)
		}
		return err
	}()

	sourceReader.Read()
	mastercatIndexer.Stop()
	mastercatWriter.Stop()
	if cfg.CatalogIndexer.Source.Metadata {
		metadataIndexer.Stop()
		metadataWriter.Stop()
	}

	slog.Info("Catalog indexer finished successfully")
	return nil
}
