package conesearch

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
)

func TestConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1, "all")
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].ID, "A")
}

func TestConesearch_WithRepositoryError(t *testing.T) {
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("Test error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	_, err = service.Conesearch(1, 1, 1, 1, "all")
	repo.AssertExpectations(t)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Test error"), err)
	}
}

func TestConesearch_WithMultipleMappers(t *testing.T) {
	vlassObjects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	ztfObjects := []repository.Mastercat{
		{ID: "ZTFA", Ra: 1, Dec: 1, Cat: "ztf"},
	}
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(vlassObjects, nil).Once()
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(ztfObjects, nil).Once()
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}, {Name: "ztf", Nside: 12}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 2, "all")
	repo.AssertExpectations(t)

	require.Len(t, result, 2)
	ids := make([]string, 2)
	cats := make([]string, 2)
	for i := range result {
		for j := range result[i].Data {
			ids[i] = result[i].Data[j].ID
			cats[i] = result[i].Data[j].Cat
		}
	}
	require.Subset(t, ids, []string{"A", "ZTFA"})
	require.Subset(t, cats, []string{"vlass", "ztf"})
}

func TestBulkConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	type testCase struct {
		ra        []float64
		dec       []float64
		radius    float64
		nneighbor int
		expected  []string
	}

	testCases := []testCase{
		{ra: []float64{1}, dec: []float64{1}, radius: 1, nneighbor: 100, expected: []string{"A"}},
		{ra: []float64{10}, dec: []float64{10}, radius: 1, nneighbor: 100, expected: []string{"B"}},
		{ra: []float64{1, 2, 3}, dec: []float64{1, 2, 3}, radius: 1, nneighbor: 100, expected: []string{"A"}},
		{ra: []float64{1, 10}, dec: []float64{1, 10}, radius: 1, nneighbor: 100, expected: []string{"A", "B"}},
	}

	for _, tc := range testCases {
		result, err := service.BulkConesearch(tc.ra, tc.dec, tc.radius, tc.nneighbor, "all", 2, 1)
		require.NoError(t, err)
		repo.AssertExpectations(t)

		require.Lenf(t, result, len(tc.expected), "test case: %v", tc)
		for i := range result {
			for j := range result[i].Data {
				require.Contains(t, tc.expected, result[i].Data[j].ID, "test case: %v", tc)
			}
		}
	}
}

func TestBulkConesearch_WithRepositoryError(t *testing.T) {
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	_, err = service.BulkConesearch([]float64{1, 10}, []float64{1, 10}, 1, 100, "all", 2, 1)
	repo.AssertExpectations(t)
	require.Error(t, err)
	require.Equal(t, "repository error", err.Error())
}

type mockAllwiseStore struct {
	t       *testing.T
	objects []repository.GetAllwiseFromPixelsRow
}

func (m mockAllwiseStore) InsertAllwiseWithoutParams(context.Context, repository.Allwise) error   { return nil }
func (m mockAllwiseStore) GetAllwise(context.Context, string) (repository.GetAllwiseRow, error)    { return repository.GetAllwiseRow{}, nil }
func (m mockAllwiseStore) BulkInsertAllwise(context.Context, *sql.DB, []any) error                 { return nil }
func (m mockAllwiseStore) BulkGetAllwise(context.Context, []string) ([]repository.BulkGetAllwiseRow, error) {
	return nil, nil
}
func (m mockAllwiseStore) GetAllwiseFromPixels(ctx context.Context, pixels []int64) ([]repository.GetAllwiseFromPixelsRow, error) {
	return m.objects, nil
}

func TestConesearch_WithMetadata(t *testing.T) {
	objects := []repository.GetAllwiseFromPixelsRow{
		{ID: "A", Ra: 1, Dec: 1},
		{ID: "B", Ra: 10, Dec: 10},
	}

	store := mockAllwiseStore{t: t, objects: objects}
	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", store)

	repo := NewMockMastercatStore(t)
	catalogs := []repository.Catalog{{Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithResolver(resolver), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "allwise")
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].GetId(), "A")
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(f)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(f, err)

	f.Add(float64(1), float64(1), float64(1), int(1))
	f.Fuzz(func(t *testing.T, ra float64, dec float64, radius float64, nneighbor int) {
		_, err := service.Conesearch(ra, dec, radius, nneighbor, "all")
		if err == nil {
			repo.AssertExpectations(t)
		}
	})
}
