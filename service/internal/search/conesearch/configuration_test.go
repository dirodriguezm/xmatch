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
