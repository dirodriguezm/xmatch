package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllwiseInputSchemaToMastercat(t *testing.T) {
	a := AllwiseInputSchema{
		Designation:  "designation",
		Ra:           0.0,
		Dec:          0.0,
		W1mpro:       0.0,
		W1sigmpro:    0.0,
		W2mpro:       0.0,
		W2sigmpro:    0.0,
		W3mpro:       0.0,
		W3sigmpro:    0.0,
		W4mpro:       0.0,
		W4sigmpro:    0.0,
		J_m_2mass:    0.0,
		H_m_2mass:    0.0,
		K_m_2mass:    0.0,
		J_msig_2mass: 0.0,
		H_msig_2mass: 0.0,
		K_msig_2mass: 0.0,
	}
	expected := Mastercat{
		ID:   "designation",
		Ipix: 0,
		Ra:   0.0,
		Dec:  0.0,
		Cat:  "allwise",
	}
	actual := a.ToMastercat()
	require.Equal(t, expected, actual)
}
