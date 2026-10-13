package erosita

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoundTrip pins the catalog conversion and insert path end to end:
// InputSchema -> conversion methods -> sqlite insert -> generated queries.
func TestRoundTrip(t *testing.T) {
	queries := testutils.CatalogDB(t)
	ctx := context.Background()

	raw := InputSchema{
		IAUNAME: "eROSITA J123",
		DETUID:  "det-123",
		SKYTILE: 17,
		ID_SRC:  99,
		UID:     1001,
		RA:      55.5,
		DEC:     -12.3,
	}

	ra, dec := raw.GetCoordinates()
	mastercat := repository.Mastercat{
		ID:   raw.MastercatID(),
		Ipix: 11,
		Ra:   ra,
		Dec:  dec,
		Cat:  "erosita",
	}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{mastercat}))

	metadata, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Erosita{metadata}))

	byID, err := queries.GetErosita(ctx, "eROSITA J123")
	require.NoError(t, err)
	assert.Equal(t, "det-123", byID.Detuid.String)
	assert.Equal(t, int64(17), byID.Skytile.Int64)
	assert.Equal(t, int64(99), byID.IDSrc.Int64)
	assert.Equal(t, int64(1001), byID.Uid.Int64)
	// The by-id response exposes the mastercat pair under ra/dec; the
	// catalog table's own coordinates are not part of the response.
	assert.Equal(t, 55.5, byID.Ra)
	assert.Equal(t, -12.3, byID.Dec)

	bulk, err := queries.BulkGetErosita(ctx, []string{"eROSITA J123"})
	require.NoError(t, err)
	require.Len(t, bulk, 1)
	assert.Equal(t, "eROSITA J123", bulk[0].ID)
	assert.Equal(t, 55.5, bulk[0].Ra)
	assert.Equal(t, -12.3, bulk[0].Dec)

	objects, err := queries.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Equal(t, []repository.Mastercat{mastercat}, objects)
}
