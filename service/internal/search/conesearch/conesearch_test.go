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

package conesearch

import (
	"errors"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"

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
	repo := &MockRepository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1, "all")
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].ID, "A")
}

func TestConesearch_WithRepositoryError(t *testing.T) {
	repo := &MockRepository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("Test error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
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
	repo := &MockRepository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(vlassObjects, nil).Once()
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(ztfObjects, nil).Once()
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}, {Name: "ztf", Nside: 12}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 2, "all")
	repo.AssertExpectations(t)

	// both objects in the result should be in the same coordinates, but different catalog
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
	repo := &MockRepository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
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

func TestConesearch_WithMetadata(t *testing.T) {
	objects := []repository.GetAllwiseFromPixelsRow{
		{ID: "A", Ra: 1, Dec: 1},
		{ID: "B", Ra: 10, Dec: 10},
	}
	repo := &MockRepository{}
	repo.On("GetAllwiseFromPixels", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "allwise")
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].GetId(), "A")
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := &MockRepository{}
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithRepository(repo), WithCatalogs(catalogs))
	require.NoError(f, err)

	f.Add(float64(1), float64(1), float64(1), int(1))
	f.Fuzz(func(t *testing.T, ra float64, dec float64, radius float64, nneighbor int) {
		_, err := service.Conesearch(ra, dec, radius, nneighbor, "all")
		if err == nil {
			repo.AssertExpectations(t)
		}
	})
}
