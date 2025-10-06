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

package di

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	metadata_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	reader_factory "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/factory"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/golobby/container/v3"
)

func RegisterLogger(ctr container.Container, stdout io.Writer) {
	ctr.Singleton(func() *slog.LevelVar {
		levels := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"info":  slog.LevelInfo,
			"error": slog.LevelError,
			"warn":  slog.LevelWarn,
			"":      slog.LevelInfo,
		}
		var programLevel = new(slog.LevelVar)
		logger := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{Level: programLevel}))
		slog.SetDefault(logger)
		programLevel.Set(levels[os.Getenv("LOG_LEVEL")])
		return programLevel
	})
}

func RegisterConfig(ctr container.Container, getenv func(string) string) {
	cfg, err := config.Load(getenv)
	if err != nil {
		panic(err)
	}
	ctr.Singleton(func() *config.Config {
		return cfg
	})
}

func RegisterDB(ctr container.Container) {
	ctr.Singleton(func(cfg *config.Config) *sql.DB {
		conn := cfg.CatalogIndexer.Database.Url
		// Ensure write access with proper SQLite parameters
		if !strings.Contains(conn, "?") {
			conn += "?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
		}
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite connection", "conn", conn)
			panic(err)
		}
		db.SetMaxOpenConns(1)    // SQLite only supports 1 writer connection
		db.SetMaxIdleConns(1)    // Keep 1 idle connection
		db.SetConnMaxLifetime(0) // Connections don't expire
		db.SetConnMaxIdleTime(0) // Idle connections don't expire
		_, err = db.Exec("select 'test conn'")
		if err != nil {
			slog.Error("Could not connect to database", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})
}

func RegisterRepository(ctr container.Container) {
	ctr.Singleton(func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})
}

func RegisterCatalogRegister(ctr container.Container, ctx context.Context) {
	ctr.Singleton(func(repo conesearch.Repository, cfg *config.Config) *indexer.CatalogRegister {
		return indexer.NewCatalogRegister(ctx, repo, *cfg.CatalogIndexer.Source)
	})
}

func RegisterSource(ctr container.Container) {
	ctr.Singleton(func(cfg *config.Config) *source.Source {
		src, err := source.NewSource(cfg.CatalogIndexer.Source)
		if err != nil {
			slog.Error("Could not register Source")
			panic(err)
		}
		return src
	})
}

func RegisterMastercatWriter(ctr container.Container, ctx context.Context) {
	ctr.NamedSingleton("mastercat_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
	) *actor.Actor {
		if cfg.CatalogIndexer.IndexerWriter == nil {
			panic("Indexer writer not configured")
		}
		switch cfg.CatalogIndexer.IndexerWriter.Type {
		case "parquet":
			cfg.CatalogIndexer.IndexerWriter.Schema = config.MastercatSchema
			w, err := parquet_writer.New[repository.Mastercat](cfg.CatalogIndexer.IndexerWriter, ctx)
			if err != nil {
				panic(err)
			}
			return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
		case "sqlite":
			w := sqlite_writer.New(repo, ctx, "mastercat")
			return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
		default:
			slog.Error("Writer type not allowed", "type", cfg.CatalogIndexer.IndexerWriter.Type)
			panic("Writer type not allowed")
		}
	})
}

func RegisterMastercatIndexer(ctr container.Container, ctx context.Context) {
	ctr.NamedSingleton("mastercat_indexer", func(src *source.Source, cfg *config.Config) *actor.Actor {
		ind, err := mastercat_indexer.New(src, cfg.CatalogIndexer.Indexer)
		if err != nil {
			panic(err)
		}
		var mastercatWriter *actor.Actor
		ctr.NamedResolve(&mastercatWriter, "mastercat_writer")
		return actor.New(cfg.CatalogIndexer.ChannelSize, ind.Index, nil, []*actor.Actor{mastercatWriter}, ctx)
	})
}

func RegisterMetadataWriter(ctr container.Container, ctx context.Context) {
	ctr.NamedSingleton("metadata_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
	) *actor.Actor {
		if cfg.CatalogIndexer.MetadataWriter == nil {
			slog.Info("Skipping registration for metadata writer. MetadataWriter not configured")
			return nil
		}
		switch cfg.CatalogIndexer.MetadataWriter.Type {
		case "parquet":
			var w writer.Writer
			var err error
			switch cfg.CatalogIndexer.MetadataWriter.Schema {
			case config.AllwiseSchema:
				w, err = parquet_writer.New[repository.Allwise](cfg.CatalogIndexer.MetadataWriter, ctx)
			case config.VlassSchema:
				w, err = parquet_writer.New[repository.VlassObjectSchema](cfg.CatalogIndexer.MetadataWriter, ctx)
			case config.GaiaSchema:
				w, err = parquet_writer.New[repository.Gaia](cfg.CatalogIndexer.MetadataWriter, ctx)
			default:
				err = fmt.Errorf("Unknown schema %v", cfg.CatalogIndexer.MetadataWriter.Schema)
			}
			if err != nil {
				panic(err)
			}
			return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
		case "sqlite":
			switch cfg.CatalogIndexer.MetadataWriter.Schema {
			case config.AllwiseSchema:
				w := sqlite_writer.New(repo, ctx, "allwise")
				return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
			case config.VlassSchema:
				w := sqlite_writer.New(repo, ctx, "vlass")
				return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
			case config.GaiaSchema:
				w := sqlite_writer.New(repo, ctx, "gaia")
				return actor.New(cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx)
			default:
				panic(fmt.Errorf("Unknown schema %v", cfg.CatalogIndexer.MetadataWriter.Schema))
			}
		default:
			panic(fmt.Errorf("Unknown Metadata Writer Type: %s", cfg.CatalogIndexer.MetadataWriter.Type))
		}
	})
}

func RegisterMetadataIndexer(ctr container.Container, ctx context.Context) {
	ctr.NamedSingleton("metadata_indexer", func(src *source.Source, cfg *config.Config) *actor.Actor {
		var metadataWriter *actor.Actor
		ctr.NamedResolve(&metadataWriter, "metadata_writer")
		indexer := metadata_indexer.New(cfg.CatalogIndexer.Source.CatalogName)
		return actor.New(cfg.CatalogIndexer.ChannelSize, indexer.Index, nil, []*actor.Actor{metadataWriter}, ctx)
	})
}

func RegisterReader(ctr container.Container) {
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *reader.SourceReader {
		r, err := reader_factory.ReaderFactory(src, cfg.CatalogIndexer.Reader)
		if err != nil {
			slog.Error("Could not register reader", "error", err, "source", src, "config", cfg.CatalogIndexer.Reader)
			panic(err)
		}
		var mastercatIndexer *actor.Actor
		var metadataIndexer *actor.Actor
		err = ctr.NamedResolve(&mastercatIndexer, "mastercat_indexer")
		if err != nil {
			panic(fmt.Errorf("Could not resolve mastercat indexer: %w", err))
		}
		err = ctr.NamedResolve(&metadataIndexer, "metadata_indexer")
		if err != nil {
			panic(fmt.Errorf("Could not resolve metadata indexer: %w", err))
		}

		sourceReader := reader.SourceReader{
			Reader:    r,
			Src:       src,
			BatchSize: cfg.CatalogIndexer.Reader.BatchSize,
			Receivers: []*actor.Actor{mastercatIndexer},
		}
		if cfg.CatalogIndexer.Source.Metadata {
			sourceReader.Receivers = append(sourceReader.Receivers, metadataIndexer)
		}
		return &sourceReader
	})
}

func BuildIndexerContainer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) container.Container {
	ctr := container.New()

	// Order is important
	// Pipeline is:
	// Reader -> MastercatIndexer -> MastercatWriter
	//        -> MetadataIndexer -> MetadataWriter
	// The order is reversed: Writer, Indexer, Reader
	RegisterConfig(ctr, getenv)
	RegisterLogger(ctr, stdout)
	RegisterDB(ctr)
	RegisterRepository(ctr)
	RegisterCatalogRegister(ctr, ctx)
	RegisterSource(ctr)
	RegisterMastercatWriter(ctr, ctx)
	RegisterMetadataWriter(ctr, ctx)
	RegisterMastercatIndexer(ctr, ctx)
	RegisterMetadataIndexer(ctr, ctx)
	RegisterReader(ctr)

	return ctr
}
