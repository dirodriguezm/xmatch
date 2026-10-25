package allwise

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

	sourceID := "allwise-001"
	cntr := int64(42)
	ra := 11.1
	dec := -22.2
	w1 := 12.3
	raw := InputSchema{
		Source_id: &sourceID,
		Cntr:      &cntr,
		Ra:        &ra,
		Dec:       &dec,
		W1mpro:    &w1,
	}

	rawRA, rawDec := raw.GetCoordinates()
	mastercat := repository.Mastercat{
		ID:   raw.MastercatID(),
		Ipix: 3,
		Ra:   rawRA,
		Dec:  rawDec,
		Cat:  "allwise",
	}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{mastercat}))

	metadata, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Allwise{metadata}))

	byID, err := queries.GetAllwise(ctx, "allwise-001")
	require.NoError(t, err)
	assert.Equal(t, int64(42), byID.Cntr)
	assert.Equal(t, mastercat.Ra, byID.Ra)
	assert.Equal(t, mastercat.Dec, byID.Dec)
	assert.NotNil(t, byID.W1mpro)
	assert.Equal(t, 12.3, *byID.W1mpro)

	bulk, err := queries.BulkGetAllwise(ctx, []string{"allwise-001"})
	require.NoError(t, err)
	require.Len(t, bulk, 1)
	assert.Equal(t, "allwise-001", bulk[0].ID)

	objects, err := queries.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Equal(t, []repository.Mastercat{mastercat}, objects)

	// A nil photometry column maps to a NULL stored value.
	missingID := "allwise-002"
	missing := InputSchema{Source_id: &missingID, Cntr: &cntr, Ra: &ra, Dec: &dec}
	missingMastercat := repository.Mastercat{
		ID:   missing.MastercatID(),
		Ipix: 4,
		Ra:   rawRA,
		Dec:  rawDec,
		Cat:  "allwise",
	}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{missingMastercat}))

	missingMetadata, err := missing.ConvertToMetadata()
	require.NoError(t, err)
	require.Nil(t, missingMetadata.W1mpro)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Allwise{missingMetadata}))

	missingByID, err := queries.GetAllwise(ctx, missingID)
	require.NoError(t, err)
	assert.Nil(t, missingByID.W1mpro)
}
