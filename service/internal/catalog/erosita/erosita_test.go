package erosita

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

	raw := repository.ErositaInputSchema{
		IAUNAME: "eROSITA J123456.7-765432",
		RA:      45.0,
		DEC:     30.0,
	}
	mc, err := adapter.ConvertToMastercat(raw, mapper)
	require.NoError(t, err)
	assert.Equal(t, "eROSITA J123456.7-765432", mc.ID)
	assert.Equal(t, 45.0, mc.Ra)
	assert.Equal(t, 30.0, mc.Dec)
	assert.Equal(t, "erosita", mc.Cat)
	assert.NotZero(t, mc.Ipix)
}

func TestConvertToMetadataFromRaw(t *testing.T) {
	adapter := Adapter{}

	raw := repository.ErositaInputSchema{
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
	md, err := adapter.ConvertToMetadataFromRaw(repository.ErositaInputSchema{})
	require.NoError(t, err)
	assert.Equal(t, "eROSITA", md.GetCatalog())
}

func TestConvertToMetadataFromRowType(t *testing.T) {
	adapter := Adapter{}
	row := repository.GetErositaFromPixelsRow{
		ID:     "from_db",
		Detuid: repository.NullString{NullString: sql.NullString{String: "det456", Valid: true}},
		Ra:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 50.0, Valid: true}},
	}
	md := adapter.ConvertToMetadata(row)
	result := md.(repository.Erosita)
	assert.Equal(t, "from_db", result.ID)
	assert.Equal(t, "det456", result.Detuid.String)
	assert.Equal(t, 50.0, result.Ra.Float64)
}
