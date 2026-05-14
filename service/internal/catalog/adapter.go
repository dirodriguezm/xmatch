// Package catalog provides all catalog related operations
package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type CatalogAdapter interface {
	Name() string

	NewRawRecord() any

	NewMetadataRecord() any

	BulkInsertMetadata(ctx context.Context, rows []any) error

	GetByID(ctx context.Context, id string) (any, error)

	BulkGetByID(ctx context.Context, ids []string) (any, error)

	GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error)

	GetCoordinates(raw any) (float64, float64, error)

	ConvertToMastercat(raw any, ipix int64) (repository.Mastercat, error)

	ConvertToMetadataFromRaw(raw any) (any, error)
}

var factories = map[string]func(*repository.Queries) (CatalogAdapter, error){}

func Register(name string, factory func(*repository.Queries) (CatalogAdapter, error)) {
	factories[strings.ToLower(name)] = factory
}

type Resolver struct {
	repo *repository.Queries
}

func NewResolver(repo *repository.Queries) *Resolver {
	return &Resolver{repo: repo}
}

func (r *Resolver) Has(name string) bool {
	_, ok := factories[strings.ToLower(name)]
	return ok
}

func (r *Resolver) Get(name string) (CatalogAdapter, error) {
	factory, ok := factories[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown catalog: %s", name)
	}
	if r.repo == nil {
		return nil, fmt.Errorf("no repository registered for catalog: %s", name)
	}
	return factory(r.repo)
}
