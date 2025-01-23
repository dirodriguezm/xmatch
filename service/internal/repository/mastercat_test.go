package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToInsertObjectParams(t *testing.T) {
	m := Mastercat{
		ID:   "id",
		Ipix: 1,
		Ra:   1,
		Dec:  1,
		Cat:  "cat",
	}
	result := m.ToInsertObjectParams()
	require.Equal(t, m.ID, result.ID)
	require.Equal(t, m.Ipix, result.Ipix)
	require.Equal(t, m.Ra, result.Ra)
	require.Equal(t, m.Dec, result.Dec)
	require.Equal(t, m.Cat, result.Cat)
}

func TestToParquetMastercat(t *testing.T) {
	m := Mastercat{
		ID:   "id",
		Ipix: 1,
		Ra:   1,
		Dec:  1,
		Cat:  "cat",
	}
	result := m.ToParquetMastercat()
	require.Equal(t, m.ID, result.ID)
	require.Equal(t, m.Ipix, result.Ipix)
	require.Equal(t, m.Ra, result.Ra)
	require.Equal(t, m.Dec, result.Dec)
	require.Equal(t, m.Cat, result.Cat)
}
