package di

import (
	"context"
	"database/sql"
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
	readerResults := make(map[string]chan reader.ReaderResult)
	readerResults["indexer"] = make(chan reader.ReaderResult)
	readerResults["metadata"] = make(chan reader.ReaderResult)
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

	// Register indexer
	writerInput := make(chan writer.WriterInput[any])
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *mastercat_indexer.IndexerActor {
		actor, err := mastercat_indexer.New(src, readerResults["indexer"], writerInput, cfg.CatalogIndexer.Indexer)
		if err != nil {
			panic(err)
		}
		return actor
	})

	// Register metadata indexer
	writerInputMetadata := make(chan writer.WriterInput[any])
	ctr.Singleton(func(src *source.Source, cfg *config.Config) *metadata_indexer.IndexerActor {
		if !cfg.CatalogIndexer.Source.Metadata {
			return nil
		}
		actor := metadata_indexer.New(readerResults["metadata"], writerInputMetadata)
		return actor
	})

	// Register mastercat indexer writer
	ctr.NamedSingleton("indexer_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
	) writer.Writer[any] {
		if cfg.CatalogIndexer.IndexerWriter == nil {
			panic("Indexer writer not configured")
		}
		switch cfg.CatalogIndexer.IndexerWriter.Type {
		case "parquet":
			cfg.CatalogIndexer.IndexerWriter.Schema = config.MastercatSchema
			w, err := parquet_writer.NewParquetWriter(writerInput, make(chan bool), cfg.CatalogIndexer.IndexerWriter)
			if err != nil {
				panic(err)
			}
			return w
		case "sqlite":
			w := sqlite_writer.NewSqliteWriter(repo, writerInput, make(chan bool), context.TODO(), src)
			return w
		default:
			slog.Error("Writer type not allowed", "type", cfg.CatalogIndexer.IndexerWriter.Type)
			panic("Writer type not allowed")
		}
	})

	// Register metadata writer
	ctr.NamedSingleton("metadata_writer", func(
		cfg *config.Config,
		repo conesearch.Repository,
		src *source.Source,
	) writer.Writer[any] {
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
			w, err := parquet_writer.NewParquetWriter(writerInputMetadata, make(chan bool), cfg.CatalogIndexer.MetadataWriter)
			if err != nil {
				panic(err)
			}
			return w
		case "sqlite":
			w := sqlite_writer.NewSqliteWriter(repo, writerInputMetadata, make(chan bool), context.TODO(), src)
			return w
		default:
			slog.Error("Writer type not allowed", "type", cfg.CatalogIndexer.MetadataWriter.Type)
			panic("Writer type not allowed")
		}
	})

	// Register partition writer if is configured
	// ctr.Singleton(func(cfg *config.Config) indexer.Writer {
	// 	if cfg.CatalogIndexer.PartitionWriter == nil {
	// 		return nil
	// 	}
	// 	type PartitionSchema struct {
	// 		ID  string  `parquet:"name=id, type=BYTE_ARRAY"`
	// 		Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	// 		Dec float64 `parquet:"name=dec, type=DOUBLE"`
	// 	}
	// 	// TODO: Configure partition reader and writer
	// 	return nil
	// })

	// Register reducer writer if is configured
	// ctr.Singleton(func(cfg *config.Config) indexer.Writer {
	// 	if cfg.CatalogIndexer.ReducerWriter == nil {
	// 		return nil
	// 	}
	// 	type ReducerSchema struct {
	// 		ID  string  `parquet:"name=id, type=BYTE_ARRAY"`
	// 		Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	// 		Dec float64 `parquet:"name=dec, type=DOUBLE"`
	// 	}
	// 	// TODO: Configure pre partition reader and writer
	// 	return nil
	// })

	return ctr
}
