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
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func BuildServiceContainer() container.Container {
	ctr := container.New()
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
	ctr.Singleton(func() *sql.DB {
		conn := os.Getenv("DB_CONN")
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite3 connection", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})
	ctr.Singleton(func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})
	ctr.Singleton(func(r conesearch.Repository) (*conesearch.ConesearchService, error) {
		return conesearch.NewConesearchService(
			conesearch.WithNside(18),
			conesearch.WithScheme(healpix.Nest),
			conesearch.WithCatalog("vlass"),
			conesearch.WithRepository(r),
		)
	})
	ctr.Singleton(func(service *conesearch.ConesearchService) httpservice.HttpServer {
		return httpservice.NewHttpServer(service)
	})
	return ctr
}

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
	ctr.Singleton(func() *sql.DB {
		conn := os.Getenv("DB_CONN")
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite3 connection", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})
	ctr.Singleton(func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})

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
	ctr.Singleton(func(rdr *indexer.Reader, cfg *config.Config) *indexer.Indexer {
		idx, err := indexer.New(*rdr, readerResults, indexerResults, cfg.CatalogIndexer.Indexer)
		if err != nil {
			panic(err)
		}
		return idx
	})

	// Register writer
	ctr.Singleton(func(repo *conesearch.Repository) indexer.Writer {
		return writer.NewSqliteWriter(*repo, indexerResults, context.TODO())
	})
	return ctr
}
