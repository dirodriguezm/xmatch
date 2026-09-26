package erosita

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToMastercat(t *testing.T) {
	adapter := Adapter{}
	ipix := int64(123)

	raw := InputSchema{
		IAUNAME: "eROSITA J123456.7-765432",
		RA:      45.0,
		DEC:     30.0,
	}
	ra, dec, err := adapter.GetCoordinates(raw)
	require.NoError(t, err)
	assert.Equal(t, 45.0, ra)
	assert.Equal(t, 30.0, dec)

	mc, err := adapter.ConvertToMastercat(raw, ipix)
	require.NoError(t, err)
	assert.Equal(t, "eROSITA J123456.7-765432", mc.ID)
	assert.Equal(t, 45.0, mc.Ra)
	assert.Equal(t, 30.0, mc.Dec)
	assert.Equal(t, "erosita", mc.Cat)
	assert.Equal(t, ipix, mc.Ipix)
}

func TestConvertToMetadataFromRaw(t *testing.T) {
	adapter := Adapter{}

	raw := InputSchema{
		IAUNAME: "eROSITA J123456.7-765432",
		DETUID:  "det123",
		SKYTILE: 42,
		ID_SRC:  100,
		UID:     9999,
		RA:      45.0,
		DEC:     30.0,
	}
	md, err := adapter.ConvertToMetadataFromRaw(raw)
	require.NoError(t, err)
	result := md.(repository.Erosita)
	assert.Equal(t, "eROSITA J123456.7-765432", result.ID)
	assert.Equal(t, "det123", result.Detuid.String)
	assert.Equal(t, int64(42), result.Skytile.Int64)
	assert.Equal(t, int64(100), result.IDSrc.Int64)
	assert.Equal(t, int64(9999), result.Uid.Int64)
	assert.Equal(t, 45.0, result.Ra.Float64)
	assert.Equal(t, 30.0, result.Dec.Float64)
}

func TestConvertToMetadataFromRaw_ImplementsMetadataInterface(t *testing.T) {
	adapter := Adapter{}
	md, err := adapter.ConvertToMetadataFromRaw(InputSchema{})
	require.NoError(t, err)
	assert.IsType(t, repository.Erosita{}, md)
}

func TestGetFromPixelsRequiresRepository(t *testing.T) {
	adapter := Adapter{}

	_, err := adapter.GetFromPixels(context.Background(), []int64{1})
	require.EqualError(t, err, "erosita adapter has no repository")
}
