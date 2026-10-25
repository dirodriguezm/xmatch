package allwise

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCoordinates(t *testing.T) {
	t.Run("all fields populated", func(t *testing.T) {
		ra := 45.0
		dec := 30.0
		raw := InputSchema{
			Ra:  &ra,
			Dec: &dec,
		}
		raOut, decOut := raw.GetCoordinates()
		assert.Equal(t, ra, raOut)
		assert.Equal(t, dec, decOut)
	})

	t.Run("nil coordinates - zero values", func(t *testing.T) {
		raw := InputSchema{
			Ra:  nil,
			Dec: nil,
		}
		ra, dec := raw.GetCoordinates()
		assert.Equal(t, 0.0, ra)
		assert.Equal(t, 0.0, dec)
	})
}

func TestMastercatID(t *testing.T) {
	t.Run("source id populated", func(t *testing.T) {
		sourceID := "allwise_001"
		raw := InputSchema{Source_id: &sourceID}
		assert.Equal(t, "allwise_001", raw.MastercatID())
	})

	t.Run("nil source_id - empty id", func(t *testing.T) {
		raw := InputSchema{Source_id: nil}
		assert.Equal(t, "", raw.MastercatID())
	})
}

func TestConvertToMetadata(t *testing.T) {
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
		result, err := raw.ConvertToMetadata()
		require.NoError(t, err)
		assert.Equal(t, "allwise_001", result.ID)
		assert.Equal(t, int64(5), result.Cntr)
		require.NotNil(t, result.W1mpro)
		assert.Equal(t, 13.5, *result.W1mpro)
		assert.Equal(t, 0.05, *result.W1sigmpro)
		assert.Equal(t, 15.2, *result.JM2mass)
	})

	t.Run("nil fields - nil pointers", func(t *testing.T) {
		sourceID := "allwise_002"
		raw := InputSchema{
			Source_id: &sourceID,
			W1mpro:    nil,
			W1sigmpro: nil,
		}
		result, err := raw.ConvertToMetadata()
		require.NoError(t, err)
		assert.Equal(t, "allwise_002", result.ID)
		assert.Nil(t, result.W1mpro)
		assert.Nil(t, result.W1sigmpro)
	})

	t.Run("nil Source_id - empty id", func(t *testing.T) {
		raw := InputSchema{
			Source_id: nil,
		}
		result, err := raw.ConvertToMetadata()
		require.NoError(t, err)
		assert.Equal(t, "", result.ID)
	})
}

func TestConvertToMetadata_2massFields(t *testing.T) {
	jmag := 14.0
	jsig := 0.03
	hmag := 13.5
	hsig := 0.04
	kmag := 13.0
	ksig := 0.02
	raw := InputSchema{
		Source_id:    new("test"),
		J_m_2mass:    &jmag,
		J_msig_2mass: &jsig,
		H_m_2mass:    &hmag,
		H_msig_2mass: &hsig,
		K_m_2mass:    &kmag,
		K_msig_2mass: &ksig,
	}
	result, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	assert.Equal(t, 14.0, *result.JM2mass)
	assert.Equal(t, 0.03, *result.JMsig2mass)
	assert.Equal(t, 13.5, *result.HM2mass)
	assert.Equal(t, 0.04, *result.HMsig2mass)
	assert.Equal(t, 13.0, *result.KM2mass)
	assert.Equal(t, 0.02, *result.KMsig2mass)
}

func TestConvertToMetadata_WiseAllBands(t *testing.T) {
	w1 := 14.0
	w2 := 13.5
	w3 := 12.0
	w4 := 9.0
	w1s := 0.02
	w2s := 0.03
	w3s := 0.05
	w4s := 0.10
	raw := InputSchema{
		Source_id: new("test"),
		W1mpro:    &w1,
		W1sigmpro: &w1s,
		W2mpro:    &w2,
		W2sigmpro: &w2s,
		W3mpro:    &w3,
		W3sigmpro: &w3s,
		W4mpro:    &w4,
		W4sigmpro: &w4s,
	}
	result, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	assert.Equal(t, 14.0, *result.W1mpro)
	assert.Equal(t, 0.02, *result.W1sigmpro)
	assert.Equal(t, 13.5, *result.W2mpro)
	assert.Equal(t, 0.03, *result.W2sigmpro)
	assert.Equal(t, 12.0, *result.W3mpro)
	assert.Equal(t, 0.05, *result.W3sigmpro)
	assert.Equal(t, 9.0, *result.W4mpro)
	assert.Equal(t, 0.10, *result.W4sigmpro)
}
