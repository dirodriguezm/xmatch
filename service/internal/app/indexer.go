package app

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

func Config(getenv func(string) string) (config.Config, error) {
	return config.Load(getenv)
}

func Repository(cfg config.Config) (*repository.Queries, error) {
	conn := cfg.CatalogIndexer.Database.Url
	if !strings.Contains(conn, "?") {
		conn += "?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
	}
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, fmt.Errorf("could not create sqlite connection: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	_, err = db.Exec("select 'test conn'")
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	return repository.New(db), nil
}

func CatalogRegister(ctx context.Context, registry conesearch.CatalogRegistry, srcConfig config.SourceConfig) *indexer.CatalogRegister {
	return indexer.NewCatalogRegister(ctx, registry, srcConfig)
}

func Source(cfg config.SourceConfig) (*source.Source, error) {
	return source.NewSource(cfg)
}
