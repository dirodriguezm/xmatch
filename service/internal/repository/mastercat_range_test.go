package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dirodriguezm/healpix"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestFindObjectsInPixelRanges(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS mastercat (
			id text not null,
			ipix bigint not null,
			ra double precision not null,
			dec double precision not null,
			cat text not null,
			PRIMARY KEY (id, cat)
		);
		CREATE INDEX IF NOT EXISTS mastercat_ipix_idx ON mastercat (ipix);
	`)
	require.NoError(t, err)

	ctx := context.Background()
	queries := New(db)
	objects := []Mastercat{
		{ID: "before", Ipix: 9, Ra: 0, Dec: 0, Cat: "test"},
		{ID: "start", Ipix: 10, Ra: 0, Dec: 0, Cat: "test"},
		{ID: "middle", Ipix: 15, Ra: 0, Dec: 0, Cat: "test"},
		{ID: "stop", Ipix: 20, Ra: 0, Dec: 0, Cat: "test"},
		{ID: "second", Ipix: 30, Ra: 0, Dec: 0, Cat: "test"},
	}
	for _, obj := range objects {
		require.NoError(t, queries.InsertObject(ctx, InsertObjectParams(obj)))
	}

	result, err := queries.FindObjectsInPixelRanges(ctx, []healpix.PixelRange{
		{Start: 10, Stop: 20},
		{Start: 30, Stop: 31},
	})
	require.NoError(t, err)

	ids := make([]string, 0, len(result))
	for _, obj := range result {
		ids = append(ids, obj.ID)
	}
	require.ElementsMatch(t, []string{"start", "middle", "second"}, ids)
}

func TestFindObjectsInPixelRangesAcceptsHealpixRangeList(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS mastercat (
			id text not null,
			ipix bigint not null,
			ra double precision not null,
			dec double precision not null,
			cat text not null,
			PRIMARY KEY (id, cat)
		);
		CREATE INDEX IF NOT EXISTS mastercat_ipix_idx ON mastercat (ipix);
	`)
	require.NoError(t, err)

	ctx := context.Background()
	queries := New(db)
	require.NoError(t, queries.InsertObject(ctx, InsertObjectParams{ID: "matched", Ipix: 999, Ra: 0, Dec: 0, Cat: "test"}))

	ranges := make([]healpix.PixelRange, 451)
	for i := range ranges {
		ranges[i] = healpix.PixelRange{Start: int64(i * 2), Stop: int64(i*2 + 1)}
	}
	ranges = append(ranges, healpix.PixelRange{Start: 999, Stop: 1000})

	result, err := queries.FindObjectsInPixelRanges(ctx, ranges)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "matched", result[0].ID)
}
