package core

import (
	"errors"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepository struct {
	mock.Mock
	Objects []MastercatObject
	Error   error
}

func (m *MockRepository) FindObjects(pixelList []int64) ([]MastercatObject, error) {
	m.Called(pixelList)
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Objects, nil
}

func TestConesearch(t *testing.T) {
	objects := []MastercatObject{
		{id: "A", ra: 1, dec: 1, catalog: "vlass"},
		{id: "B", ra: 10, dec: 10, catalog: "vlass"},
	}
	repo := &MockRepository{
		Objects: objects,
		Error:   nil,
	}
	repo.On("FindObjects", mock.Anything).Return(objects)
	service, err := NewConesearchService(WithCatalog("vlass"), WithScheme(healpix.Nest), WithNside(18), WithRepository(repo))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1)
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].id, "A")
}

func TestConesearchWithRepositoryError(t *testing.T) {
	repo := &MockRepository{
		Objects: nil,
		Error:   errors.New("Test error"),
	}
	repo.On("FindObjects", mock.Anything).Return(errors.New("Test error"))
	service, err := NewConesearchService(WithCatalog("vlass"), WithScheme(healpix.Nest), WithNside(18), WithRepository(repo))
	require.NoError(t, err)

	_, err = service.Conesearch(1, 1, 1, 1)
	repo.AssertExpectations(t)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Test error"), err)
	}
}

func FuzzConesearch(f *testing.F) {
	objects := []MastercatObject{
		{id: "A", ra: 1, dec: 1, catalog: "vlass"},
		{id: "B", ra: 10, dec: 10, catalog: "vlass"},
	}
	repo := &MockRepository{
		Objects: objects,
		Error:   nil,
	}
	repo.On("FindObjects", mock.Anything).Return(objects)
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