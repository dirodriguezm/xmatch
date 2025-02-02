package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToInsertObjectParams(t *testing.T) {
	id := "id"
	cat := "cat"
	ra := 1.0
	dec := 1.0
	ipix := int64(1)
	m := ParquetMastercat{
		ID:   &id,
		Ipix: &ipix,
		Ra:   &ra,
		Dec:  &dec,
		Cat:  &cat,
	}
	result := m.ToInsertObjectParams()
	require.Equal(t, *m.ID, result.ID)
	require.Equal(t, *m.Ipix, result.Ipix)
	require.Equal(t, *m.Ra, result.Ra)
	require.Equal(t, *m.Dec, result.Dec)
	require.Equal(t, *m.Cat, result.Cat)
}
