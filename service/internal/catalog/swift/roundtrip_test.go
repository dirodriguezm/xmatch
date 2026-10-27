package swift

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	queries := testutils.CatalogDB(t)
	ctx := context.Background()

	raw := InputSchema{ID: "LSXPS J000000.0+000000", RA: 12.3, Dec: -4.5}

	ra, dec := raw.GetCoordinates()
	mastercat := repository.Mastercat{
		ID:   raw.MastercatID(),
		Ipix: 42,
		Ra:   ra,
		Dec:  dec,
		Cat:  "swift",
	}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{mastercat}))

	metadata, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Swift{metadata}))

	byID, err := queries.GetSwift(ctx, "LSXPS J000000.0+000000")
	require.NoError(t, err)
	assert.Equal(t, mastercat.Ra, byID.Ra)
	assert.Equal(t, mastercat.Dec, byID.Dec)

	bulk, err := queries.BulkGetSwift(ctx, []string{"LSXPS J000000.0+000000"})
	require.NoError(t, err)
	require.Len(t, bulk, 1)
	assert.Equal(t, "LSXPS J000000.0+000000", bulk[0].ID)
}
