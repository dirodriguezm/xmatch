package gaia

import (
	"database/sql"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToMastercat(t *testing.T) {
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Ring)
	require.NoError(t, err)
	adapter := Adapter{}

	raw := repository.GaiaInputSchema{
		Designation: "Gaia DR3 123456789",
		RA:          45.0,
		Dec:         30.0,
	}
	mc, err := adapter.ConvertToMastercat(raw, mapper)
	require.NoError(t, err)
	assert.Equal(t, "Gaia DR3 123456789", mc.ID)
	assert.Equal(t, 45.0, mc.Ra)
	assert.Equal(t, 30.0, mc.Dec)
	assert.Equal(t, "gaia", mc.Cat)
	assert.NotZero(t, mc.Ipix)
}

func TestConvertToMetadataFromRaw(t *testing.T) {
	adapter := Adapter{}

	raw := repository.GaiaInputSchema{
		Designation:         "Gaia DR3 123456789",
		PhotGMeanFlux:       1000.5,
		PhotGMeanFluxError:  10.2,
		PhotGMeanMag:        15.3,
		PhotBpMeanFlux:      800.3,
		PhotBpMeanFluxError: 8.1,
		PhotBpMeanMag:       16.0,
		PhotRpMeanFlux:      600.7,
		PhotRpMeanFluxError: 6.5,
		PhotRpMeanMag:       14.8,
	}
	md, err := adapter.ConvertToMetadataFromRaw(raw)
	require.NoError(t, err)
	result := md.(repository.Gaia)
	assert.Equal(t, "Gaia DR3 123456789", result.ID)
	assert.Equal(t, 1000.5, result.PhotGMeanFlux.Float64)
	assert.InDelta(t, 10.2, result.PhotGMeanFluxError.Float64, 0.001)
	assert.InDelta(t, 15.3, result.PhotGMeanMag.Float64, 0.001)
	assert.Equal(t, 800.3, result.PhotBpMeanFlux.Float64)
	assert.InDelta(t, 8.1, result.PhotBpMeanFluxError.Float64, 0.001)
	assert.InDelta(t, 16.0, result.PhotBpMeanMag.Float64, 0.001)
	assert.Equal(t, 600.7, result.PhotRpMeanFlux.Float64)
	assert.InDelta(t, 6.5, result.PhotRpMeanFluxError.Float64, 0.001)
	assert.InDelta(t, 14.8, result.PhotRpMeanMag.Float64, 0.001)
}

func TestConvertToMetadataFromRaw_ImplementsMetadataInterface(t *testing.T) {
	adapter := Adapter{}
	md, err := adapter.ConvertToMetadataFromRaw(repository.GaiaInputSchema{})
	require.NoError(t, err)
	assert.Equal(t, "GAIA/DR3", md.GetCatalog())
}

func TestConvertToMetadataFromRowType(t *testing.T) {
	adapter := Adapter{}
	row := repository.GetGaiaFromPixelsRow{
		ID:                  "from_db",
		PhotGMeanFlux:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 1000.5, Valid: true}},
		PhotGMeanFluxError:  repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 10.2, Valid: true}},
		PhotGMeanMag:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 15.3, Valid: true}},
		PhotBpMeanFlux:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 800.3, Valid: true}},
		PhotBpMeanFluxError: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 8.1, Valid: true}},
		PhotBpMeanMag:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 16.0, Valid: true}},
	}
	md := adapter.ConvertToMetadata(row)
	result := md.(repository.Gaia)
	assert.Equal(t, "from_db", result.ID)
	assert.Equal(t, 1000.5, result.PhotGMeanFlux.Float64)
	assert.Equal(t, 16.0, result.PhotBpMeanMag.Float64)
}
