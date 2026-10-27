package swift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCoordinates(t *testing.T) {
	raw := InputSchema{RA: 45.0, Dec: -30.0}
	ra, dec := raw.GetCoordinates()
	assert.Equal(t, 45.0, ra)
	assert.Equal(t, -30.0, dec)
}

func TestMastercatID(t *testing.T) {
	raw := InputSchema{ID: "LSXPS J000000.0+000000"}
	assert.Equal(t, "LSXPS J000000.0+000000", raw.MastercatID())
}

func TestConvertToMetadata(t *testing.T) {
	raw := InputSchema{ID: "LSXPS J000000.0+000000", RA: 1.0, Dec: 2.0}
	result, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	assert.Equal(t, "LSXPS J000000.0+000000", result.ID)
}
