package conesearch

import (
	"errors"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"testing"

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
	repo := &MockRepository{
		Objects: objects,
		Error:   nil,
	}
	repo.On("FindObjects", mock.Anything).Return(objects, nil)
	service, err := NewConesearchService(WithCatalog("vlass"), WithScheme(healpix.Nest), WithNside(18), WithRepository(repo))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1)
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].ID, "A")
}

func TestConesearchWithRepositoryError(t *testing.T) {
	repo := &MockRepository{
		Objects: nil,
		Error:   errors.New("Test error"),
	}
	repo.On("FindObjects", mock.Anything).Return(nil, errors.New("Test error"))
	service, err := NewConesearchService(WithCatalog("vlass"), WithScheme(healpix.Nest), WithNside(18), WithRepository(repo))
	require.NoError(t, err)

	_, err = service.Conesearch(1, 1, 1, 1)
	repo.AssertExpectations(t)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Test error"), err)
	}
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := &MockRepository{
		Objects: objects,
		Error:   nil,
	}
	repo.On("FindObjects", mock.Anything).Return(objects, nil)
	service, err := NewConesearchService(WithCatalog("vlass"), WithScheme(healpix.Nest), WithNside(18), WithRepository(repo))
	require.NoError(f, err)

	f.Add(float64(1), float64(1), float64(1), int(1))
	f.Fuzz(func(t *testing.T, ra float64, dec float64, radius float64, nneighbor int) {
		_, err := service.Conesearch(ra, dec, radius, nneighbor)
		if err == nil {
			repo.AssertExpectations(t)
		}
	})
}
