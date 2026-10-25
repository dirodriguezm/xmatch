package gaia

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
		Designation:    "Gaia DR3 123",
		RA:             33.3,
		Dec:            44.4,
		PhotGMeanFlux:  1000.5,
		PhotBpMeanMag:  16.1,
		PhotRpMeanFlux: 600.7,
	}

	ra, dec := raw.GetCoordinates()
	mastercat := repository.Mastercat{
		ID:   raw.MastercatID(),
		Ipix: 7,
		Ra:   ra,
		Dec:  dec,
		Cat:  "gaia",
	}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{mastercat}))

	metadata, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Gaia{metadata}))

	byID, err := queries.GetGaia(ctx, "Gaia DR3 123")
	require.NoError(t, err)
	assert.Equal(t, mastercat.Ra, byID.Ra)
	assert.Equal(t, mastercat.Dec, byID.Dec)
	assert.True(t, byID.PhotGMeanFlux.Valid)
	assert.Equal(t, 1000.5, byID.PhotGMeanFlux.Float64)
	assert.InDelta(t, 16.1, byID.PhotBpMeanMag.Float64, 0.0001)

	bulk, err := queries.BulkGetGaia(ctx, []string{"Gaia DR3 123"})
	require.NoError(t, err)
	require.Len(t, bulk, 1)
	assert.Equal(t, "Gaia DR3 123", bulk[0].ID)

	objects, err := queries.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Equal(t, []repository.Mastercat{mastercat}, objects)
}
