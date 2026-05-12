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
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
)

type MetadataService struct {
	resolver *catalog.Resolver
}

func NewMetadataService(resolver *catalog.Resolver) (*MetadataService, error) {
	if resolver == nil {
		return nil, fmt.Errorf("resolver was nil while creating MetadataService")
	}
	return &MetadataService{resolver: resolver}, nil
}

func (m *MetadataService) FindByID(ctx context.Context, id string, catalogName string) (any, error) {
	if err := m.validateCatalog(catalogName); err != nil {
		return nil, err
	}
	if err := m.validateID(id); err != nil {
		return nil, err
	}

	return m.queryCatalog(ctx, id, catalogName)
}

func (m *MetadataService) BulkFindByID(ctx context.Context, ids []string, catalogName string) (any, error) {
	if err := m.validateCatalog(catalogName); err != nil {
		return nil, err
	}
	for i := range ids {
		if err := m.validateID(ids[i]); err != nil {
			return nil, err
		}
	}

	return m.bulkQueryCatalog(ctx, ids, catalogName)
}

func (m *MetadataService) queryCatalog(ctx context.Context, id string, catalogName string) (any, error) {
	adapter, err := m.resolver.GetQuery(catalogName)
	if err != nil {
		return nil, ArgumentError{Name: "catalog", Value: catalogName, Reason: err.Error()}
	}
	return adapter.GetByID(ctx, id)
}

func (m *MetadataService) bulkQueryCatalog(ctx context.Context, ids []string, catalogName string) (any, error) {
	adapter, err := m.resolver.GetQuery(catalogName)
	if err != nil {
		return nil, ArgumentError{Name: "catalog", Value: catalogName, Reason: err.Error()}
	}
	return adapter.BulkGetByID(ctx, ids)
}

func (m *MetadataService) validateCatalog(catalogName string) error {
	if !m.resolver.Has(catalogName) {
		return ValidationError{
			Field:  "catalog",
			Reason: fmt.Sprintf("unknown catalog: %s", catalogName),
			Value:  catalogName,
		}
	}
	return nil
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
