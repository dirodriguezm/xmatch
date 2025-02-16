package metadata

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Repository interface {
	InsertAllwise(context.Context, repository.InsertAllwiseParams) error
	GetAllwise(context.Context, string) (repository.Allwise, error)
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

func (m *MetadataService) queryCatalog(ctx context.Context, id string, catalog string) (any, error) {
	switch strings.ToLower(catalog) {
	case "allwise":
		result, err := m.repository.GetAllwise(ctx, id)
		if err != nil {
			return nil, err
		}
		return result.ToAllwiseMetadata(), nil
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

func (m *MetadataService) validateID(_ string) error {
	return nil
}
