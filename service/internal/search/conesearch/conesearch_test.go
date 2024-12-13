package conesearch

import (
	"errors"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/mocks"

	"github.com/dirodriguezm/healpix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := &mocks.Repository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1)
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].ID, "A")
}

func TestConesearch_WithRepositoryError(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("Test error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	_, err = service.Conesearch(1, 1, 1, 1)
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
	repo := &mocks.Repository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(vlassObjects, nil).Once()
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(ztfObjects, nil).Once()
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}, {Name: "ztf", Nside: 12}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 2)
	repo.AssertExpectations(t)

	// both objects in the result should be in the same coordinates, but different catalog
	require.Len(t, result, 2)
	ids := make([]string, 2)
	cats := make([]string, 2)
	for i := range result {
		ids[i] = result[i].ID
		cats[i] = result[i].Cat
	}
	require.Subset(t, ids, []string{"A", "ZTFA"})
	require.Subset(t, cats, []string{"vlass", "ztf"})
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := &mocks.Repository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(f, err)

	f.Add(float64(1), float64(1), float64(1), int(1))
	f.Fuzz(func(t *testing.T, ra float64, dec float64, radius float64, nneighbor int) {
		_, err := service.Conesearch(ra, dec, radius, nneighbor)
		if err == nil {
			repo.AssertExpectations(t)
		}
	})
}
