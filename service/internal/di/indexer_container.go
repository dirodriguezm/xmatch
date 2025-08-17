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

func BuildIndexerContainer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) container.Container {
	ctr := container.New()
	// read config
	cfg, err := config.Load(getenv)
	if err != nil {
		panic(err)
	}

	ctr.Singleton(func() *config.Config {
		return cfg
	})

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

	// Register DB
	ctr.Singleton(func(cfg *config.Config) *sql.DB {
		conn := cfg.CatalogIndexer.Database.Url
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite connection", "conn", conn)
			panic(err)
		}
		_, err = db.Exec("select 'test conn'")
		if err != nil {
			slog.Error("Could not connect to database", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})

	// Register Repository
	ctr.Singleton(func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})

	// Register CatalogRegister
	ctr.Singleton(func(repo conesearch.Repository, cfg *config.Config) *indexer.CatalogRegister {
		return indexer.NewCatalogRegister(ctx, repo, *cfg.CatalogIndexer.Source)
	})

	// Register Source
	ctr.Singleton(func(cfg *config.Config) *source.Source {
		src, err := source.NewSource(cfg.CatalogIndexer.Source)
		if err != nil {
			slog.Error("Could not register Source")
			panic(err)
		}
		return src
	})

	// Register reader
	readerResults := make(map[string]chan reader.ReaderResult)
	readerResults["indexer"] = make(chan reader.ReaderResult, cfg.CatalogIndexer.ChannelSize)
	readerResults["metadata"] = make(chan reader.ReaderResult, cfg.CatalogIndexer.ChannelSize)
	ctr.Singleton(func(src *source.Source, cfg *config.Config) reader.Reader {
		outputChannels := []chan reader.ReaderResult{readerResults["indexer"]}
		if cfg.CatalogIndexer.Source.Metadata {
			outputChannels = append(outputChannels, readerResults["metadata"])
		}
		r, err := reader_factory.ReaderFactory(src, outputChannels, cfg.CatalogIndexer.Reader)
		if err != nil {
			slog.Error("Could not register reader", "error", err, "source", src, "config", cfg.CatalogIndexer.Reader)
			panic(err)
		}
		return r
	})

	// Register mastercat indexer
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *mastercat_indexer.IndexerActor {
		writerInput := make(chan writer.WriterInput[repository.Mastercat], cfg.CatalogIndexer.ChannelSize)
		ind, err := mastercat_indexer.New(src, readerResults["indexer"], writerInput, cfg.CatalogIndexer.Indexer)
		if err != nil {
			panic(err)
		}
		return ind
	})

	// Register mastercat indexer writer
	ctr.NamedSingleton("indexer_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
		ind *mastercat_indexer.IndexerActor,
	) writer.Writer[repository.Mastercat] {
		if cfg.CatalogIndexer.IndexerWriter == nil {
			panic("Indexer writer not configured")
		}
		switch cfg.CatalogIndexer.IndexerWriter.Type {
		case "parquet":
			cfg.CatalogIndexer.IndexerWriter.Schema = config.MastercatSchema
			w, err := parquet_writer.NewParquetWriter(ind.GetOutbox(), make(chan struct{}), cfg.CatalogIndexer.IndexerWriter, ctx)
			if err != nil {
				panic(err)
			}
			return w
		case "sqlite":
			w := sqlite_writer.NewSqliteWriter(repo, ind.GetOutbox(), make(chan struct{}), ctx, src)
			return w
		default:
			slog.Error("Writer type not allowed", "type", cfg.CatalogIndexer.IndexerWriter.Type)
			panic("Writer type not allowed")
		}
	})

	// Register metadata indexer
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *metadata_indexer.IndexerActor {
		outbox := make(chan writer.WriterInput[repository.Metadata], cfg.CatalogIndexer.ChannelSize)
		var parser metadata_indexer.MetadataParser[repository.Metadata]
		switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
		case "allwise":
			parser = metadata_indexer.AllwiseMetadataParser{}
		case "vlass":
			parser = metadata_indexer.VlassMetadataParser{}
		default:
			panic(fmt.Errorf("Unknown catalog name: %s", cfg.CatalogIndexer.Source.CatalogName))
		}
		return metadata_indexer.New(readerResults["metadata"], outbox, parser)
	})

	// Register metadata writer
	ctr.NamedSingleton("metadata_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
		ind *metadata_indexer.IndexerActor,
	) writer.Writer[repository.Metadata] {
		if cfg.CatalogIndexer.MetadataWriter == nil {
			slog.Info("Skipping registration for metadata writer. MetadataWriter not configured")
			return nil
		}
		switch cfg.CatalogIndexer.MetadataWriter.Type {
		case "parquet":
			switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
			case "allwise":
				cfg.CatalogIndexer.MetadataWriter.Schema = config.AllwiseSchema
			default:
				panic("Unknown catalog name")
			}
			w, err := parquet_writer.NewParquetWriter(ind.GetOutbox(), make(chan struct{}), cfg.CatalogIndexer.MetadataWriter, ctx)
			if err != nil {
				panic(err)
			}
			return w
		case "sqlite":
			w := sqlite_writer.NewSqliteWriter(repo, ind.GetOutbox(), make(chan struct{}), ctx, src)
			return w
		default:
			panic(fmt.Errorf("Unknown catalog name: %s", cfg.CatalogIndexer.Source.CatalogName))
		}
	})

	return ctr
}
