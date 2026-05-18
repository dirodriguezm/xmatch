package allwise

import (
	"database/sql"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToMastercat(t *testing.T) {
	adapter := Adapter{}
	ipix := int64(123)

	t.Run("all fields populated", func(t *testing.T) {
		sourceID := "allwise_001"
		ra := 45.0
		dec := 30.0
		raw := InputSchema{
			Source_id: &sourceID,
			Ra:        &ra,
			Dec:       &dec,
		}
		raOut, decOut, err := adapter.GetCoordinates(raw)
		require.NoError(t, err)
		assert.Equal(t, ra, raOut)
		assert.Equal(t, dec, decOut)

		mc, err := adapter.ConvertToMastercat(raw, ipix)
		require.NoError(t, err)
		assert.Equal(t, "allwise_001", mc.ID)
		assert.Equal(t, 45.0, mc.Ra)
		assert.Equal(t, 30.0, mc.Dec)
		assert.Equal(t, "allwise", mc.Cat)
		assert.Equal(t, ipix, mc.Ipix)
	})

	t.Run("nil coordinates - zero values", func(t *testing.T) {
		sourceID := "allwise_002"
		raw := InputSchema{
			Source_id: &sourceID,
			Ra:        nil,
			Dec:       nil,
		}
		ra, dec, err := adapter.GetCoordinates(raw)
		require.NoError(t, err)
		assert.Equal(t, 0.0, ra)
		assert.Equal(t, 0.0, dec)

		mc, err := adapter.ConvertToMastercat(raw, ipix)
		require.NoError(t, err)
		assert.Equal(t, "allwise_002", mc.ID)
		assert.Equal(t, 0.0, mc.Ra)
		assert.Equal(t, 0.0, mc.Dec)
		assert.Equal(t, "allwise", mc.Cat)
	})

	t.Run("nil source_id - empty id", func(t *testing.T) {
		ra := 10.0
		dec := 20.0
		raw := InputSchema{
			Source_id: nil,
			Ra:        &ra,
			Dec:       &dec,
		}
		mc, err := adapter.ConvertToMastercat(raw, ipix)
		require.NoError(t, err)
		assert.Equal(t, "", mc.ID)
		assert.Equal(t, 10.0, mc.Ra)
		assert.Equal(t, 20.0, mc.Dec)
	})
}

func TestConvertToMetadataFromRaw(t *testing.T) {
	adapter := Adapter{}

	t.Run("all fields populated", func(t *testing.T) {
		sourceID := "allwise_001"
		cntr := int64(5)
		w1 := 13.5
		w1sig := 0.05
		jmag := 15.2
		raw := InputSchema{
			Source_id: &sourceID,
			Cntr:      &cntr,
			W1mpro:    &w1,
			W1sigmpro: &w1sig,
			J_m_2mass: &jmag,
		}
		md, err := adapter.ConvertToMetadataFromRaw(raw)
		require.NoError(t, err)
		result := md.(repository.Allwise)
		assert.Equal(t, "allwise_001", result.ID)
		assert.Equal(t, int64(5), result.Cntr)
		assert.Equal(t, 13.5, result.W1mpro.Float64)
		assert.True(t, result.W1mpro.Valid)
		assert.Equal(t, 0.05, result.W1sigmpro.Float64)
		assert.Equal(t, 15.2, result.JM2mass.Float64)
	})

	t.Run("nil fields - sentinel values", func(t *testing.T) {
		sourceID := "allwise_002"
		raw := InputSchema{
			Source_id: &sourceID,
			W1mpro:    nil,
			W1sigmpro: nil,
		}
		md, err := adapter.ConvertToMetadataFromRaw(raw)
		require.NoError(t, err)
		result := md.(repository.Allwise)
		assert.Equal(t, "allwise_002", result.ID)
		assert.Equal(t, -9999.0, result.W1mpro.Float64)
		assert.False(t, result.W1mpro.Valid)
		assert.Equal(t, -9999.0, result.W1sigmpro.Float64)
		assert.False(t, result.W1sigmpro.Valid)
	})

	t.Run("nil Source_id - empty id", func(t *testing.T) {
		raw := InputSchema{
			Source_id: nil,
		}
		md, err := adapter.ConvertToMetadataFromRaw(raw)
		require.NoError(t, err)
		result := md.(repository.Allwise)
		assert.Equal(t, "", result.ID)
	})

	t.Run("metadata type assertions", func(t *testing.T) {
		md, _ := adapter.ConvertToMetadataFromRaw(InputSchema{})
		assert.IsType(t, repository.Allwise{}, md)
	})
}

func TestConvertToMetadataFromRaw_2massFields(t *testing.T) {
	adapter := Adapter{}

	jmag := 14.0
	jsig := 0.03
	hmag := 13.5
	hsig := 0.04
	kmag := 13.0
	ksig := 0.02
	raw := InputSchema{
		Source_id:    strPtr("test"),
		J_m_2mass:    &jmag,
		J_msig_2mass: &jsig,
		H_m_2mass:    &hmag,
		H_msig_2mass: &hsig,
		K_m_2mass:    &kmag,
		K_msig_2mass: &ksig,
	}
	md, err := adapter.ConvertToMetadataFromRaw(raw)
	require.NoError(t, err)
	result := md.(repository.Allwise)
	assert.Equal(t, 14.0, result.JM2mass.Float64)
	assert.Equal(t, 0.03, result.JMsig2mass.Float64)
	assert.Equal(t, 13.5, result.HM2mass.Float64)
	assert.Equal(t, 0.04, result.HMsig2mass.Float64)
	assert.Equal(t, 13.0, result.KM2mass.Float64)
	assert.Equal(t, 0.02, result.KMsig2mass.Float64)
}

func TestConvertToMetadataFromRaw_WiseAllBands(t *testing.T) {
	adapter := Adapter{}

	w1 := 14.0
	w2 := 13.5
	w3 := 12.0
	w4 := 9.0
	w1s := 0.02
	w2s := 0.03
	w3s := 0.05
	w4s := 0.10
	raw := InputSchema{
		Source_id: strPtr("test"),
		W1mpro:    &w1,
		W1sigmpro: &w1s,
		W2mpro:    &w2,
		W2sigmpro: &w2s,
		W3mpro:    &w3,
		W3sigmpro: &w3s,
		W4mpro:    &w4,
		W4sigmpro: &w4s,
	}
	md, err := adapter.ConvertToMetadataFromRaw(raw)
	require.NoError(t, err)
	result := md.(repository.Allwise)
	assert.Equal(t, 14.0, result.W1mpro.Float64)
	assert.Equal(t, 0.02, result.W1sigmpro.Float64)
	assert.Equal(t, 13.5, result.W2mpro.Float64)
	assert.Equal(t, 0.03, result.W2sigmpro.Float64)
	assert.Equal(t, 12.0, result.W3mpro.Float64)
	assert.Equal(t, 0.05, result.W3sigmpro.Float64)
	assert.Equal(t, 9.0, result.W4mpro.Float64)
	assert.Equal(t, 0.10, result.W4sigmpro.Float64)
}

func TestConvertToMetadata_ImplementsMetadataInterface(t *testing.T) {
	adapter := Adapter{}
	md, err := adapter.ConvertToMetadataFromRaw(InputSchema{})
	require.NoError(t, err)
	assert.IsType(t, repository.Allwise{}, md)
}

func TestConvertFromPixelsRowToMetadata(t *testing.T) {
	row := repository.GetAllwiseFromPixelsRow{
		ID:     "from_db",
		Cntr:   42,
		W1mpro: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 14.0, Valid: true}},
	}

	md := convertAllwiseFromPixelsRowToMetadata(row)
	result := md.Object.(repository.Allwise)
	assert.Equal(t, "from_db", result.ID)
	assert.Equal(t, int64(42), result.Cntr)
	assert.Equal(t, 14.0, result.W1mpro.Float64)
}

func strPtr(s string) *string {
	return &s
}
