// Package catalog provides all catalog related operations
package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type CatalogIndexAdapter interface {
	Name() string

	NewRawRecord() any

	NewParquetWriter(cfg config.WriterConfig, ctx context.Context) (writer.Writer, error)

	NewParquetReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error)

	NewFitsReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error)

	BulkInsertFn() func(context.Context, *sql.DB, []any) error

	ConvertToMastercat(raw any, mapper *healpix.HEALPixMapper) (repository.Mastercat, error)

	ConvertToMetadataFromRaw(raw any) (any, error)
}

type CatalogQueryAdapter interface {
	Name() string

	GetByID(ctx context.Context, id string) (any, error)

	BulkGetByID(ctx context.Context, ids []string) (any, error)

	GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error)

	ConvertToMetadata(obj any) repository.Metadata
}

var factories = map[string]func(any) (CatalogIndexAdapter, error){}

func Register(name string, factory func(any) (CatalogIndexAdapter, error)) {
	factories[strings.ToLower(name)] = factory
}

type Resolver struct {
	stores map[string]any
}

func NewResolver() *Resolver {
	return &Resolver{stores: map[string]any{}}
}

func (r *Resolver) RegisterStore(name string, store any) {
	r.stores[strings.ToLower(name)] = store
}

func (r *Resolver) Has(name string) bool {
	_, ok := factories[strings.ToLower(name)]
	return ok
}

func (r *Resolver) Get(name string) (CatalogIndexAdapter, error) {
	factory, ok := factories[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown catalog: %s", name)
	}
	store, ok := r.stores[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("no store registered for catalog: %s", name)
	}
	return factory(store)
}

func (r *Resolver) GetQuery(name string) (CatalogQueryAdapter, error) {
	adapter, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	queryAdapter, ok := adapter.(CatalogQueryAdapter)
	if !ok {
		return nil, fmt.Errorf("catalog %s does not support query operations", name)
	}
	return queryAdapter, nil
}
