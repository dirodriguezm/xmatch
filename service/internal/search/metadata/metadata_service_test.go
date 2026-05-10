package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
)

type mockAllwiseStore struct {
	getAllwise    func(ctx context.Context, id string) (repository.GetAllwiseRow, error)
	bulkGetAllwise func(ctx context.Context, ids []string) ([]repository.BulkGetAllwiseRow, error)
}

func (m mockAllwiseStore) InsertAllwiseWithoutParams(context.Context, repository.Allwise) error {
	return nil
}
func (m mockAllwiseStore) GetAllwise(ctx context.Context, id string) (repository.GetAllwiseRow, error) {
	return m.getAllwise(ctx, id)
}
func (m mockAllwiseStore) BulkInsertAllwise(context.Context, *sql.DB, []any) error {
	return nil
}
func (m mockAllwiseStore) BulkGetAllwise(ctx context.Context, ids []string) ([]repository.BulkGetAllwiseRow, error) {
	return m.bulkGetAllwise(ctx, ids)
}
func (m mockAllwiseStore) GetAllwiseFromPixels(context.Context, []int64) ([]repository.GetAllwiseFromPixelsRow, error) {
	return nil, nil
}

func TestMetadata_ValidateCatalog(t *testing.T) {
	m := &MetadataService{resolver: catalog.NewResolver()}

	err := m.validateCatalog("allwise")
	require.Nil(t, err)
	err = m.validateCatalog("AllWise")
	require.Nil(t, err)

	err = m.validateCatalog("gaia")
	require.Nil(t, err)
	err = m.validateCatalog("GAIA")
	require.Nil(t, err)

	err = m.validateCatalog("erosita")
	require.Nil(t, err)
	err = m.validateCatalog("Erosita")
	require.Nil(t, err)

	err = m.validateCatalog("invalid")
	require.NotNil(t, err)
	require.Equal(t, "Could not parse field catalog with value invalid: unknown catalog: invalid", err.Error())
}

func TestMetadata_FindByID(t *testing.T) {
	store := mockAllwiseStore{
		getAllwise: func(ctx context.Context, id string) (repository.GetAllwiseRow, error) {
			return repository.GetAllwiseRow{ID: "allwise1", Ra: 12.34, Dec: 56.78}, nil
		},
	}
	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", store)

	m := &MetadataService{resolver: resolver}

	result, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.Nil(t, err)
	require.Equal(t, "allwise1", result.(repository.GetAllwiseRow).ID)
}

func TestMetadata_BulkFindByID(t *testing.T) {
	store := mockAllwiseStore{
		bulkGetAllwise: func(ctx context.Context, ids []string) ([]repository.BulkGetAllwiseRow, error) {
			return []repository.BulkGetAllwiseRow{
				{ID: "allwise1", Ra: 12.34, Dec: 56.78},
				{ID: "allwise2", Ra: 23.45, Dec: 67.89},
			}, nil
		},
	}
	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", store)

	m := &MetadataService{resolver: resolver}

	result, err := m.BulkFindByID(context.Background(), []string{"allwise1", "allwise2"}, "allwise")
	require.Nil(t, err)
	expectedIds := []string{"allwise1", "allwise2"}
	for i := 0; i < len(result.([]repository.BulkGetAllwiseRow)); i++ {
		require.Equal(t, expectedIds[i], result.([]repository.BulkGetAllwiseRow)[i].ID)
	}
}

func TestMetadata_Bulk_EmptyResult(t *testing.T) {
	store := mockAllwiseStore{
		getAllwise: func(ctx context.Context, id string) (repository.GetAllwiseRow, error) {
			return repository.GetAllwiseRow{}, sql.ErrNoRows
		},
	}
	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", store)

	m := &MetadataService{resolver: resolver}

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.NotNil(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestMetadata_SomeDBError(t *testing.T) {
	store := mockAllwiseStore{
		getAllwise: func(ctx context.Context, id string) (repository.GetAllwiseRow, error) {
			return repository.GetAllwiseRow{}, fmt.Errorf("db error")
		},
	}
	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", store)

	m := &MetadataService{resolver: resolver}

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.NotNil(t, err)
	require.EqualError(t, err, "db error")
}
