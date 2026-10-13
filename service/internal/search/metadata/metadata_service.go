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

// Package metadata provides a metadata service to query metadata from catalogs
package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type MetadataService struct {
	queries *repository.Queries
}

func NewMetadataService(queries *repository.Queries) (*MetadataService, error) {
	if queries == nil {
		return nil, fmt.Errorf("queries was nil while creating MetadataService")
	}
	return &MetadataService{queries: queries}, nil
}

func (m *MetadataService) FindByID(ctx context.Context, id string, catalogName string) (any, error) {
	if err := m.validateID(id); err != nil {
		return nil, err
	}

	return m.queryCatalog(ctx, id, catalogName)
}

func (m *MetadataService) BulkFindByID(ctx context.Context, ids []string, catalogName string) (any, error) {
	for i := range ids {
		if err := m.validateID(ids[i]); err != nil {
			return nil, err
		}
	}

	return m.bulkQueryCatalog(ctx, ids, catalogName)
}

// queryCatalog dispatches a by-id lookup to the generated query for the
// catalog. The switch default carries the unknown-catalog validation, so an
// unregistered name never reaches SQL.
func (m *MetadataService) queryCatalog(ctx context.Context, id string, catalogName string) (any, error) {
	switch strings.ToLower(catalogName) {
	case "allwise":
		return m.queries.GetAllwise(ctx, id)
	case "gaia":
		return m.queries.GetGaia(ctx, id)
	case "erosita":
		return m.queries.GetErosita(ctx, id)
	default:
		return nil, unknownCatalogError(catalogName)
	}
}

// bulkQueryCatalog dispatches a bulk lookup to the generated query for the
// catalog, with the same unknown-catalog default as queryCatalog.
func (m *MetadataService) bulkQueryCatalog(ctx context.Context, ids []string, catalogName string) (any, error) {
	switch strings.ToLower(catalogName) {
	case "allwise":
		return m.queries.BulkGetAllwise(ctx, ids)
	case "gaia":
		return m.queries.BulkGetGaia(ctx, ids)
	case "erosita":
		return m.queries.BulkGetErosita(ctx, ids)
	default:
		return nil, unknownCatalogError(catalogName)
	}
}

// GetFromPixelRanges returns the metadata envelopes of the objects whose
// mastercat pixel falls inside one of the ranges. It dispatches to the
// generated range primitive per catalog; a name without a registered catalog
// contributes nothing without error (the silent skip conesearch relies on).
func (m *MetadataService) GetFromPixelRanges(ctx context.Context, catalogName string, pixelRanges []healpix.PixelRange) ([]repository.Metadata, error) {
	switch strings.ToLower(catalogName) {
	case "allwise":
		return m.queries.QueryMetadataFromPixelRanges(ctx, "allwise", pixelRanges, func(rows *sql.Rows) (repository.Metadata, error) {
			return metadataFromRow[repository.Allwise](rows, "AllWISE")
		})
	case "gaia":
		return m.queries.QueryMetadataFromPixelRanges(ctx, "gaia", pixelRanges, func(rows *sql.Rows) (repository.Metadata, error) {
			return metadataFromRow[repository.Gaia](rows, "GAIA/DR3")
		})
	case "erosita":
		return m.queries.QueryMetadataFromPixelRanges(ctx, "erosita", pixelRanges, func(rows *sql.Rows) (repository.Metadata, error) {
			return metadataFromRow[repository.Erosita](rows, "eROSITA")
		})
	default:
		return nil, nil
	}
}

// scanInto scans the current row into dst followed by extra, in order. dst
// must point to a struct whose exported fields, in declaration order, match
// the query's leading columns: the generated models mirror the tables' CREATE
// TABLE order, which is the order `SELECT <table>.*` returns.
func scanInto(rows *sql.Rows, dst any, extra ...any) error {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("scanInto: dst must be a pointer to a struct, got %T", dst)
	}

	value = value.Elem()
	targets := make([]any, 0, value.NumField()+len(extra))
	for i := range value.NumField() {
		field := value.Field(i)
		if !field.CanAddr() || !field.CanInterface() {
			continue
		}
		targets = append(targets, field.Addr().Interface())
	}
	targets = append(targets, extra...)

	return rows.Scan(targets...)
}

// metadataFromRow scans one catalog row plus the trailing mastercat
// coordinates selected by QueryMetadataFromPixelRanges and builds the
// Metadata envelope. displayName is the catalog's display name literal; the
// envelope id comes from the model's ID field.
func metadataFromRow[T any](rows *sql.Rows, displayName string) (repository.Metadata, error) {
	var object T
	var ra, dec float64
	if err := scanInto(rows, &object, &ra, &dec); err != nil {
		return repository.Metadata{}, err
	}

	id := reflect.ValueOf(object).FieldByName("ID")
	if !id.IsValid() || id.Kind() != reflect.String {
		return repository.Metadata{}, fmt.Errorf("metadataFromRow: %T has no string ID field", object)
	}

	return repository.Metadata{
		ID:      id.String(),
		Catalog: displayName,
		Ra:      ra,
		Dec:     dec,
		Object:  object,
	}, nil
}

func unknownCatalogError(catalogName string) error {
	return ValidationError{
		Field:  "catalog",
		Reason: fmt.Sprintf("unknown catalog: %s", catalogName),
		Value:  catalogName,
	}
}

func (m *MetadataService) validateID(id string) error {
	if err := ensureNoSQLInjection(id); err != nil {
		return ValidationError{
			Field:  "id",
			Reason: err.Error(),
			Value:  id,
		}
	}

	return nil
}

func ensureNoSQLInjection(s string) error {
	dangerousPatterns := []string{
		"--",
		";",
		"'",
		"\"",
		"/*",
		"*/",
		"union",
		"select",
		"drop",
		"delete",
		"update",
		"insert",
		"exec",
		"execute",
		"alter",
		"truncate",
	}

	lowerID := strings.ToLower(s)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerID, pattern) {
			return fmt.Errorf("id contains potentially dangerous pattern: %s", pattern)
		}
	}

	return nil
}
