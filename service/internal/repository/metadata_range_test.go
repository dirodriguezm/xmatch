package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestQueryMetadataFromPixelRanges(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		queries := newMetadataRangeTestQueries(t)
		seedMetadataRangeTestObjects(t, queries, []Mastercat{
			{ID: "before", Ipix: 9, Cat: "test"},
			{ID: "start", Ipix: 10, Cat: "test"},
			{ID: "middle", Ipix: 15, Cat: "test"},
			{ID: "stop", Ipix: 20, Cat: "test"},
			{ID: "second", Ipix: 30, Cat: "test"},
		})

		result, err := queries.QueryMetadataFromPixelRanges(context.Background(), "allwise", []healpix.PixelRange{
			{Start: 10, Stop: 20},
			{Start: 30, Stop: 31},
		}, metadataRangeIDMapper)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{"start", "middle", "second"}, metadataIDs(result))
	})

	t.Run("stress 451 ranges", func(t *testing.T) {
		queries := newMetadataRangeTestQueries(t)
		seedMetadataRangeTestObjects(t, queries, []Mastercat{{ID: "matched", Ipix: 999, Cat: "test"}})

		ranges := make([]healpix.PixelRange, 451)
		for i := range ranges {
			ranges[i] = healpix.PixelRange{Start: int64(i * 2), Stop: int64(i*2 + 1)}
		}
		ranges = append(ranges, healpix.PixelRange{Start: 999, Stop: 1000})

		result, err := queries.QueryMetadataFromPixelRanges(context.Background(), "allwise", ranges, metadataRangeIDMapper)
		require.NoError(t, err)
		require.Equal(t, []string{"matched"}, metadataIDs(result))
	})

	t.Run("empty and degenerate ranges return nothing", func(t *testing.T) {
		queries := newMetadataRangeTestQueries(t)
		seedMetadataRangeTestObjects(t, queries, []Mastercat{{ID: "matched", Ipix: 999, Cat: "test"}})

		cases := []struct {
			name   string
			ranges []healpix.PixelRange
		}{
			{name: "nil input"},
			{name: "empty input", ranges: []healpix.PixelRange{}},
			{name: "all ranges degenerate", ranges: []healpix.PixelRange{
				{Start: 999, Stop: 999},
				{Start: 1000, Stop: 999},
			}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := queries.QueryMetadataFromPixelRanges(context.Background(), "allwise", tc.ranges, metadataRangeIDMapper)
				require.NoError(t, err)
				require.Empty(t, result)
			})
		}
	})

	t.Run("degenerate ranges are skipped", func(t *testing.T) {
		queries := newMetadataRangeTestQueries(t)
		seedMetadataRangeTestObjects(t, queries, []Mastercat{{ID: "matched", Ipix: 999, Cat: "test"}})

		result, err := queries.QueryMetadataFromPixelRanges(context.Background(), "allwise", []healpix.PixelRange{
			{Start: 999, Stop: 999},
			{Start: 999, Stop: 1000},
		}, metadataRangeIDMapper)
		require.NoError(t, err)
		require.Equal(t, []string{"matched"}, metadataIDs(result))
	})

	t.Run("map errors propagate", func(t *testing.T) {
		queries := newMetadataRangeTestQueries(t)
		seedMetadataRangeTestObjects(t, queries, []Mastercat{{ID: "matched", Ipix: 999, Cat: "test"}})

		mapErr := errors.New("map failed")
		result, err := queries.QueryMetadataFromPixelRanges(
			context.Background(),
			"allwise",
			[]healpix.PixelRange{{Start: 999, Stop: 1000}},
			func(*sql.Rows) (Metadata, error) { return Metadata{}, mapErr },
		)
		require.ErrorIs(t, err, mapErr)
		require.Nil(t, result)
	})
}

// metadataRangeIDMapper scans any metadata table row and keeps its id column
// (always the first column of the SELECT built by QueryMetadataFromPixelRanges).
func metadataRangeIDMapper(rows *sql.Rows) (Metadata, error) {
	columns, err := rows.Columns()
	if err != nil {
		return Metadata{}, err
	}
	values := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	if err := rows.Scan(dest...); err != nil {
		return Metadata{}, err
	}
	id, _ := values[0].(string)
	return Metadata{ID: id}, nil
}

func metadataIDs(result []Metadata) []string {
	ids := make([]string, 0, len(result))
	for _, item := range result {
		ids = append(ids, item.ID)
	}
	return ids
}

func newMetadataRangeTestQueries(t *testing.T) *Queries {
	t.Helper()

	dbFile := filepath.Join(t.TempDir(), "metadata-range.db")
	mig, err := migrate.New(fmt.Sprintf("file://%s", metadataRangeMigrationsDir(t)), fmt.Sprintf("sqlite3://%s", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = mig.Close()
	})
	require.NoError(t, mig.Up())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, db.Ping())

	return New(db)
}

func metadataRangeMigrationsDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(file), "..", "db", "migrations")
}

// seedMetadataRangeTestObjects inserts one mastercat row plus one allwise
// metadata row per object, sharing the object id.
func seedMetadataRangeTestObjects(t *testing.T, q *Queries, objects []Mastercat) {
	t.Helper()

	ctx := context.Background()
	for _, obj := range objects {
		require.NoError(t, q.InsertObject(ctx, InsertObjectParams(obj)))
		require.NoError(t, q.InsertAllwise(ctx, InsertAllwiseParams{ID: obj.ID}))
	}
}
