package core

import (
	"testing"

	"github.com/dirodriguezm/healpix"
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

func TestWithNside(t *testing.T) {
	service := &ConesearchService{}

	err := WithNside(4)(service)
	require.NoError(t, err)

	require.Equal(t, service.Nside, 4)
}

func TestWithNsideInvalid(t *testing.T) {
	service := &ConesearchService{}
	nsides := []int{0, -1, 30}

	for _, nside := range nsides {
		err := WithNside(nside)(service)
		require.Error(t, err)
		require.Zero(t, service.Nside)
	}
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

func TestWithCatalog(t *testing.T) {
	service := &ConesearchService{}

	err := WithCatalog("vlass")(service)
	require.NoError(t, err)
	require.Equal(t, "vlass", service.Catalog)
	err = WithCatalog("VLASS")(service)
	require.NoError(t, err)
	require.Equal(t, "vlass", service.Catalog)
}

func TestWithCatalogInvalid(t *testing.T) {
	service := &ConesearchService{}

	err := WithCatalog("invalid")(service)
	require.Error(t, err)
	require.Zero(t, service.Catalog)
}
