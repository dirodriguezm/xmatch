package gaia

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCoordinates(t *testing.T) {
	raw := InputSchema{
		RA:  45.0,
		Dec: 30.0,
	}
	ra, dec := raw.GetCoordinates()
	assert.Equal(t, 45.0, ra)
	assert.Equal(t, 30.0, dec)
}

func TestMastercatID(t *testing.T) {
	raw := InputSchema{Designation: "Gaia DR3 123456789"}
	assert.Equal(t, "Gaia DR3 123456789", raw.MastercatID())
}

func TestConvertToMetadata(t *testing.T) {
	raw := InputSchema{
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
	result, err := raw.ConvertToMetadata()
	require.NoError(t, err)
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
