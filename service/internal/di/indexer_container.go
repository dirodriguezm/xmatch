package di

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/golobby/container/v3"
)

func BuildIndexerContainer() container.Container {
	ctr := container.New()
	// read config
	ctr.Singleton(func() *config.Config {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
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
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
		slog.SetDefault(logger)
		programLevel.Set(levels[os.Getenv("LOG_LEVEL")])
		return programLevel
	})

	// Register DB
	ctr.Singleton(func(cfg *config.Config) *sql.DB {
		conn := cfg.CatalogIndexer.Database.Url
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite3 connection", "conn", conn)
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
		ctx := context.Background()
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
	readerResults := make(chan indexer.ReaderResult)
	ctr.Singleton(func(src *source.Source, cfg *config.Config) indexer.Reader {
		r, err := reader.ReaderFactory(src, readerResults, cfg.CatalogIndexer.Reader)
		if err != nil {
			slog.Error("Could not register reader")
			panic(err)
		}
		return r
	})

	// Register indexer
	indexerResults := make(chan indexer.IndexerResult)
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *indexer.Indexer {
		idx, err := indexer.New(src, readerResults, indexerResults, cfg.CatalogIndexer.Indexer)
		if err != nil {
			panic(err)
		}
		return idx
	})

	// Register writer
	ctr.Singleton(func(repo conesearch.Repository) indexer.Writer {
		return writer.NewSqliteWriter(repo, indexerResults, make(chan bool), context.TODO())
	})
	return ctr
}
