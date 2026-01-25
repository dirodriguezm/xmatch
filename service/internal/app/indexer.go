package app

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	reader_factory "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/factory"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

const ALLWISE = "allwise"
const GAIA = "gaia"

func Config(getenv func(string) string) (config.Config, error) {
	return config.Load(getenv)
}

func Repository(cfg config.Config) (conesearch.Repository, error) {
	conn := cfg.CatalogIndexer.Database.Url
	// Ensure write access with proper SQLite parameters
	if !strings.Contains(conn, "?") {
		conn += "?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
	}
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, fmt.Errorf("Could not create sqlite connection: %w", err)
	}
	db.SetMaxOpenConns(1)    // SQLite only supports 1 writer connection
	db.SetMaxIdleConns(1)    // Keep 1 idle connection
	db.SetConnMaxLifetime(0) // Connections don't expire
	db.SetConnMaxIdleTime(0) // Idle connections don't expire
	_, err = db.Exec("select 'test conn'")
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}
	return repository.New(db), nil
}

func CatalogRegister(ctx context.Context, repo conesearch.Repository, srcConfig config.SourceConfig) *indexer.CatalogRegister {
	return indexer.NewCatalogRegister(ctx, repo, srcConfig)
}

func Source(cfg config.SourceConfig) (*source.Source, error) {
	return source.NewSource(cfg)
}

func MastercatWriter(ctx context.Context, cfg config.Config, repo conesearch.Repository, src *source.Source) (*actor.Actor, error) {
	switch cfg.CatalogIndexer.IndexerWriter.Type {
	case "parquet":
		w, err := parquet_writer.New[repository.Mastercat](cfg.CatalogIndexer.IndexerWriter, ctx)
		if err != nil {
			return nil, err
		}
		return actor.New("mastercat writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	case "sqlite":
		w := sqlite_writer.New(repo, ctx, repo.BulkInsertObject)
		return actor.New("mastercat writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	default:
		return nil, fmt.Errorf("Writer type not allowed")
	}
}

func MetadataWriter(ctx context.Context, cfg config.Config, repo conesearch.Repository, src *source.Source) (*actor.Actor, error) {
	switch cfg.CatalogIndexer.MetadataWriter.Type {
	case "parquet":
		var w writer.Writer
		var err error
		switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
		case ALLWISE:
			w, err = parquet_writer.New[repository.Allwise](cfg.CatalogIndexer.MetadataWriter, ctx)
		case GAIA:
			w, err = parquet_writer.New[repository.Gaia](cfg.CatalogIndexer.MetadataWriter, ctx)
		default:
			err = fmt.Errorf("Unknown catalog %s", cfg.CatalogIndexer.Source.CatalogName)
		}
		if err != nil {
			return nil, err
		}
		return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	case "sqlite":
		switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
		case ALLWISE:
			w := sqlite_writer.New(repo, ctx, repo.BulkInsertAllwise)
			return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
		case GAIA:
			w := sqlite_writer.New(repo, ctx, repo.BulkInsertGaia)
			return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
		default:
			return nil, fmt.Errorf("Unknown catalog %s", cfg.CatalogIndexer.Source.CatalogName)
		}
	default:
		return nil, fmt.Errorf("Unknown Metadata Writer Type: %s", cfg.CatalogIndexer.MetadataWriter.Type)
	}
}

func MastercatIndexer(cfg config.CatalogIndexerConfig, writer *actor.Actor, ctx context.Context) (*actor.Actor, error) {
	fillMastercat := func(schema repository.InputSchema, ipix int64) repository.Mastercat {
		switch cfg.Source.CatalogName {
		case ALLWISE:
			return repository.AllwiseInputSchema.FillMastercat(schema.(repository.AllwiseInputSchema), ipix)
		case GAIA:
			return repository.GaiaInputSchema.FillMastercat(schema.(repository.GaiaInputSchema), ipix)
		default:
			panic("Catalog not supported")
		}
	}

	ind, err := mastercat_indexer.New(cfg.Indexer, fillMastercat)
	if err != nil {
		return nil, err
	}
	return actor.New("mastercat indexer", cfg.ChannelSize, ind.Index, nil, []*actor.Actor{writer}, ctx), nil
}

func MetadataIndexer(cfg config.CatalogIndexerConfig, writer *actor.Actor, ctx context.Context) *actor.Actor {
	fillMetadata := func(schema repository.InputSchema) repository.Metadata {
		switch cfg.Source.CatalogName {
		case ALLWISE:
			return repository.AllwiseInputSchema.FillMetadata(schema.(repository.AllwiseInputSchema))
		case GAIA:
			return repository.GaiaInputSchema.FillMetadata(schema.(repository.GaiaInputSchema))
		default:
			panic("Catalog not supported")
		}
	}
	ind := metadata.New(fillMetadata)
	return actor.New("metadata indexer", cfg.ChannelSize, ind.Index, nil, []*actor.Actor{writer}, ctx)
}

func Reader(
	src *source.Source,
	cfg config.ReaderConfig,
	srcConfig config.SourceConfig,
	mastercatIndexer *actor.Actor,
	metadataIndexer *actor.Actor,
) (reader.SourceReader, error) {
	r, err := reader_factory.ReaderFactory(src, cfg)
	if err != nil {
		return reader.SourceReader{}, err
	}

	sourceReader := reader.SourceReader{
		Reader:    r,
		BatchSize: cfg.BatchSize,
		Receivers: []*actor.Actor{mastercatIndexer},
	}
	if srcConfig.Metadata {
		sourceReader.Receivers = append(sourceReader.Receivers, metadataIndexer)
	}
	return sourceReader, nil
}
