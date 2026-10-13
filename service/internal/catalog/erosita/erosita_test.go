package erosita

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCoordinates(t *testing.T) {
	raw := InputSchema{
		RA:  45.0,
		DEC: 30.0,
	}
	ra, dec := raw.GetCoordinates()
	assert.Equal(t, 45.0, ra)
	assert.Equal(t, 30.0, dec)
}

func TestMastercatID(t *testing.T) {
	raw := InputSchema{IAUNAME: "eROSITA J123456.7-765432"}
	assert.Equal(t, "eROSITA J123456.7-765432", raw.MastercatID())
}

func TestConvertToMetadata(t *testing.T) {
	raw := InputSchema{
		IAUNAME: "eROSITA J123456.7-765432",
		DETUID:  "det123",
		SKYTILE: 42,
		ID_SRC:  100,
		UID:     9999,
		RA:      45.0,
		DEC:     30.0,
	}
	result, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	assert.Equal(t, "eROSITA J123456.7-765432", result.ID)
	assert.Equal(t, "det123", result.Detuid.String)
	assert.Equal(t, int64(42), result.Skytile.Int64)
	assert.Equal(t, int64(100), result.IDSrc.Int64)
	assert.Equal(t, int64(9999), result.Uid.Int64)
	assert.Equal(t, 45.0, result.Ra.Float64)
	assert.Equal(t, 30.0, result.Dec.Float64)
}
