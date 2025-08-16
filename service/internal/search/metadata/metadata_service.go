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

package metadata

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Repository interface {
	InsertAllwise(context.Context, repository.InsertAllwiseParams) error
	GetAllwise(context.Context, string) (repository.Allwise, error)
	BulkGetAllwise(context.Context, []string) ([]repository.Allwise, error)
}

type MetadataService struct {
	repository Repository
}

func NewMetadataService(repo Repository) (*MetadataService, error) {
	if repo == nil {
		return nil, fmt.Errorf("Repository was nil while creating MetadataService")
	}
	return &MetadataService{repository: repo}, nil
}

func (m *MetadataService) FindByID(ctx context.Context, id string, catalog string) (any, error) {
	if err := m.validateCatalog(catalog); err != nil {
		return nil, err
	}
	if err := m.validateID(id); err != nil {
		return nil, err
	}

	return m.queryCatalog(ctx, id, catalog)
}

func (m *MetadataService) BulkFindByID(ctx context.Context, ids []string, catalog string) (any, error) {
	if err := m.validateCatalog(catalog); err != nil {
		return nil, err
	}
	for i := range ids {
		if err := m.validateID(ids[i]); err != nil {
			return nil, err
		}
	}

	return m.bulkQueryCatalog(ctx, ids, catalog)
}

func (m *MetadataService) queryCatalog(ctx context.Context, id string, catalog string) (any, error) {
	switch strings.ToLower(catalog) {
	case "allwise":
		result, err := m.repository.GetAllwise(ctx, id)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "vlass":
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Search not yet implemented for catalog"}
	case "ztf":
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Search not yet implemented for catalog"}
	default:
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Unknown catalog"}
	}
}

func (m *MetadataService) bulkQueryCatalog(ctx context.Context, ids []string, catalog string) (any, error) {
	switch strings.ToLower(catalog) {
	case "allwise":
		result, err := m.repository.BulkGetAllwise(ctx, ids)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "vlass":
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Search not yet implemented for catalog"}
	case "ztf":
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Search not yet implemented for catalog"}
	default:
		return nil, ArgumentError{Name: "catalog", Value: catalog, Reason: "Unknown catalog"}
	}
}

func (m *MetadataService) validateCatalog(catalog string) error {
	allowedCatalogs := []string{"allwise", "vlass", "ztf"}
	if !slices.Contains(allowedCatalogs, strings.ToLower(catalog)) {
		return ValidationError{
			Field:  "catalog",
			Reason: fmt.Sprintf("Allowed catalogs are %v", allowedCatalogs),
			Value:  catalog,
		}
	}
	return nil
}

func (m *MetadataService) validateID(id string) error {
	if err := ensureNoSQLInjection(id); err != nil {
		return err
	}

	if err := ensureAlphanumeric(id); err != nil {
		return err
	}

	return nil
}

func ensureAlphanumeric(s string) error {
	validIDPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validIDPattern.MatchString(s) {
		return fmt.Errorf("id must contain only alphanumeric characters, underscores, or hyphens")
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
