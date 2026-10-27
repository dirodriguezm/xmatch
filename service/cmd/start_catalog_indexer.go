package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/swift"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	metadata "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/pipeline"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
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

	db, err := app.IndexerDatabase(cfg)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			slog.Error("error closing database", "err", closeErr)
		}
	}()

	repo := app.IndexerRepository(db)
	srcCfg := cfg.CatalogIndexer.Source

	catalogRegister := app.CatalogRegister(ctx, repo, srcCfg)
	catalogRegister.RegisterCatalog()

	src, err := app.Source(srcCfg)
	if err != nil {
		return err
	}

	switch strings.ToLower(srcCfg.CatalogName) {
	case "gaia":
		return runCatalogPipeline[repository.Gaia, gaia.InputSchema](ctx, cfg, srcCfg, src, repo)
	case "allwise":
		return runCatalogPipeline[repository.Allwise, allwise.InputSchema](ctx, cfg, srcCfg, src, repo)
	case "erosita":
		return runCatalogPipeline[repository.Erosita, erosita.InputSchema](ctx, cfg, srcCfg, src, repo)
	case "swift":
		return runCatalogPipeline[repository.Swift, swift.InputSchema](ctx, cfg, srcCfg, src, repo)
	default:
		return fmt.Errorf("unknown catalog: %s", srcCfg.CatalogName)
	}
}

// catalogSchema is what the pipeline needs from a catalog: a decoded row that
// knows its coordinates, its mastercat ID and how to build its metadata model.
type catalogSchema[Md any] interface {
	mastercat_indexer.Source
	metadata.Converter[Md]
}

// runCatalogPipeline composes the catalog's typed components and runs them to
// completion. The metadata branch only exists when the source declares it.
func runCatalogPipeline[Md repository.Insertable, Raw catalogSchema[Md]](
	ctx context.Context,
	cfg config.Config,
	srcCfg config.SourceConfig,
	src *source.Source,
	repo *repository.Queries,
) error {
	ciCfg := cfg.CatalogIndexer

	mcWriter, err := pipeline.DefaultMastercatWriter(ctx, ciCfg, repo)
	if err != nil {
		return fmt.Errorf("building mastercat writer: %w", err)
	}

	var mdWriter *actor.Actor[Md]
	if srcCfg.Metadata {
		mdWriter, err = pipeline.DefaultMetadataWriter[Md](ctx, ciCfg, repo)
		if err != nil {
			return fmt.Errorf("building metadata writer: %w", err)
		}
	}

	mcIndexer, err := pipeline.DefaultMastercatIndexer[Raw](ciCfg, mcWriter, ctx)
	if err != nil {
		return fmt.Errorf("building mastercat indexer: %w", err)
	}

	var mdIndexer *actor.Actor[Raw]
	if srcCfg.Metadata {
		mdIndexer = pipeline.DefaultMetadataIndexer[Md, Raw](ciCfg, mdWriter, ctx)
	}

	srcReader, err := pipeline.DefaultSourceReader(src, ciCfg.Reader, srcCfg, mcIndexer, mdIndexer)
	if err != nil {
		return fmt.Errorf("building reader: %w", err)
	}

	p, err := pipeline.New(pipeline.PipelineConfig[Raw, Md]{
		MastercatWriter:  mcWriter,
		MastercatIndexer: mcIndexer,
		MetadataWriter:   mdWriter,
		MetadataIndexer:  mdIndexer,
		SourceReader:     srcReader,
	})
	if err != nil {
		return fmt.Errorf("error creating pipeline: %w", err)
	}

	p.Run()
	p.Stop()

	if err := p.CloseSource(); err != nil {
		return fmt.Errorf("error closing reader: %w", err)
	}

	slog.Info("Catalog indexer finished successfully")
	return nil
}
