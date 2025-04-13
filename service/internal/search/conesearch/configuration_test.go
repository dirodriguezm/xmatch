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
	"slices"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestIsPowerOfTwo(t *testing.T) {
	require.True(t, isPowerOfTwo(32))
	require.False(t, isPowerOfTwo(23))
}

func TestWithScheme(t *testing.T) {
	service := &ConesearchService{}

	err := WithScheme(healpix.Nest)(service)
	require.NoError(t, err)

	require.Equal(t, service.Scheme, healpix.Nest)
}

func TestWithResolution(t *testing.T) {
	nestService := &ConesearchService{Scheme: healpix.Nest}
	ringService := &ConesearchService{Scheme: healpix.Ring}

	err := WithResolution(4)(nestService)
	require.NoError(t, err)
	require.Equal(t, nestService.Resolution, 4)
	err = WithResolution(4)(ringService)
	require.NoError(t, err)
	require.Equal(t, ringService.Resolution, 4)
	err = WithResolution(5)(ringService)
	require.NoError(t, err)
	require.Equal(t, ringService.Resolution, 5)
}

func TestWithResolutionInvalid(t *testing.T) {
	nestService := &ConesearchService{Scheme: healpix.Nest}
	ringService := &ConesearchService{Scheme: healpix.Ring}

	err := WithResolution(5)(nestService)
	require.Error(t, err)
	require.Zero(t, nestService.Resolution)
	err = WithResolution(0)(ringService)
	require.Error(t, err)
	require.Zero(t, ringService.Resolution)
	err = WithResolution(-1)(ringService)
	require.Error(t, err)
	require.Zero(t, ringService.Resolution)
}

func TestWithCatalogs(t *testing.T) {
	service := &ConesearchService{}

	err := WithCatalogs([]repository.Catalog{{Name: "vlass", Nside: 18}})(service)
	require.NoError(t, err)
	require.True(t, slices.Contains(service.Catalogs, repository.Catalog{Name: "vlass", Nside: 18}))
	err = WithCatalogs([]repository.Catalog{{Name: "VLASS", Nside: 18}})(service)
	require.NoError(t, err)
	require.True(t, slices.Contains(service.Catalogs, repository.Catalog{Name: "VLASS", Nside: 18}))
}

func TestWithCatalogsInvalid(t *testing.T) {
	service := &ConesearchService{}

	err := WithCatalogs([]repository.Catalog{{Name: "invalid"}})(service)
	require.Error(t, err)
	require.Zero(t, service.Catalogs)
}
