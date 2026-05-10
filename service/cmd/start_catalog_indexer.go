package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/pipeline"
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

	repo, err := app.Repository(cfg)
	if err != nil {
		return err
	}

	resolver := catalog.NewResolver()
	srcCfg := cfg.CatalogIndexer.Source
	resolver.RegisterStore(srcCfg.CatalogName, repo)

	catalogRegister := app.CatalogRegister(ctx, repo, srcCfg)
	catalogRegister.RegisterCatalog()

	src, err := app.Source(srcCfg)
	if err != nil {
		return err
	}

	adapter, err := resolver.Get(srcCfg.CatalogName)
	if err != nil {
		return err
	}

	db := repo.GetDbInstance()
	pipeline, err := pipeline.New(pipeline.PipelineConfig{
		Context: ctx,
		Config:  cfg,
		DB:      db,
		Source:  src,
		Adapter: adapter,
		Store:   repo,
	})
	if err != nil {
		return fmt.Errorf("error creating pipeline: %w", err)
	}

	pipeline.Run()
	pipeline.Stop()

	if err := pipeline.CloseSource(); err != nil {
		return fmt.Errorf("Error closing reader: %w", err)
	}

	slog.Info("Catalog indexer finished successfully")
	return nil
}
