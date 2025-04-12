package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToInsertParams(t *testing.T) {
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
	result := m.ToInsertParams()
	require.Equal(t, *m.ID, result.(InsertObjectParams).ID)
	require.Equal(t, *m.Ipix, result.(InsertObjectParams).Ipix)
	require.Equal(t, *m.Ra, result.(InsertObjectParams).Ra)
	require.Equal(t, *m.Dec, result.(InsertObjectParams).Dec)
	require.Equal(t, *m.Cat, result.(InsertObjectParams).Cat)
}
