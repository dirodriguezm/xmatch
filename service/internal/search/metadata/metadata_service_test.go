package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestMetadata_FindByID(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	ctx := context.Background()
	seedAllwise(t, queries, "allwise-1", 0, 12.34, 56.78, repository.Allwise{Cntr: 7})
	seedGaia(t, queries, "gaia-1", 0, 1.5, 2.5, repository.Gaia{SourceID: nullInt64(424242)})
	seedErosita(t, queries, "erosita-1", 0, 3.5, 4.5, repository.Erosita{
		Detuid: nullString("det-1"),
		Ra:     nullFloat64(30.5),
		Dec:    nullFloat64(40.5),
	})

	t.Run("allwise", func(t *testing.T) {
		result, err := service.FindByID(ctx, "allwise-1", "allwise")
		require.NoError(t, err)
		row, ok := result.(repository.GetAllwiseRow)
		require.True(t, ok)
		require.Equal(t, "allwise-1", row.ID)
		require.Equal(t, int64(7), row.Cntr)
		require.InDelta(t, 12.34, row.Ra, 1e-9)
		require.InDelta(t, 56.78, row.Dec, 1e-9)
	})

	t.Run("gaia", func(t *testing.T) {
		result, err := service.FindByID(ctx, "gaia-1", "gaia")
		require.NoError(t, err)
		row, ok := result.(repository.GetGaiaRow)
		require.True(t, ok)
		require.Equal(t, "gaia-1", row.ID)
		require.Equal(t, int64(424242), row.SourceID.Int64)
		require.InDelta(t, 1.5, row.Ra, 1e-9)
		require.InDelta(t, 2.5, row.Dec, 1e-9)
	})

	t.Run("erosita", func(t *testing.T) {
		result, err := service.FindByID(ctx, "erosita-1", "erosita")
		require.NoError(t, err)
		row, ok := result.(repository.GetErositaRow)
		require.True(t, ok)
		require.Equal(t, "erosita-1", row.ID)
		require.Equal(t, "det-1", row.Detuid.String)
		// The response carries the mastercat pair under ra/dec, not the
		// catalog table's own (nullable) coordinates.
		require.InDelta(t, 3.5, row.Ra, 1e-9)
		require.InDelta(t, 4.5, row.Dec, 1e-9)
	})
}

func TestMetadata_BulkFindByID(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	ctx := context.Background()
	seedAllwise(t, queries, "allwise-1", 0, 12.34, 56.78, repository.Allwise{})
	seedAllwise(t, queries, "allwise-2", 0, 23.45, 67.89, repository.Allwise{})
	seedGaia(t, queries, "gaia-1", 0, 1.5, 2.5, repository.Gaia{})
	seedGaia(t, queries, "gaia-2", 0, 3.5, 4.5, repository.Gaia{})
	seedErosita(t, queries, "erosita-1", 0, 5.5, 6.5, repository.Erosita{})
	seedErosita(t, queries, "erosita-2", 0, 7.5, 8.5, repository.Erosita{})

	t.Run("allwise", func(t *testing.T) {
		result, err := service.BulkFindByID(ctx, []string{"allwise-1", "allwise-2"}, "allwise")
		require.NoError(t, err)
		rows, ok := result.([]repository.BulkGetAllwiseRow)
		require.True(t, ok)
		require.Len(t, rows, 2)
		ids := make([]string, 0, len(rows))
		for i := range rows {
			ids = append(ids, rows[i].ID)
		}
		require.ElementsMatch(t, []string{"allwise-1", "allwise-2"}, ids)
	})

	t.Run("gaia", func(t *testing.T) {
		result, err := service.BulkFindByID(ctx, []string{"gaia-1", "gaia-2"}, "gaia")
		require.NoError(t, err)
		rows, ok := result.([]repository.BulkGetGaiaRow)
		require.True(t, ok)
		require.Len(t, rows, 2)
		ids := make([]string, 0, len(rows))
		for i := range rows {
			ids = append(ids, rows[i].ID)
		}
		require.ElementsMatch(t, []string{"gaia-1", "gaia-2"}, ids)
	})

	t.Run("erosita", func(t *testing.T) {
		result, err := service.BulkFindByID(ctx, []string{"erosita-1", "erosita-2"}, "erosita")
		require.NoError(t, err)
		rows, ok := result.([]repository.BulkGetErositaRow)
		require.True(t, ok)
		require.Len(t, rows, 2)
		byID := make(map[string]repository.BulkGetErositaRow, len(rows))
		for i := range rows {
			byID[rows[i].ID] = rows[i]
		}
		require.ElementsMatch(t, []string{"erosita-1", "erosita-2"}, []string{rows[0].ID, rows[1].ID})
		// Bulk responses carry the mastercat pair under ra/dec too.
		require.InDelta(t, 5.5, byID["erosita-1"].Ra, 1e-9)
		require.InDelta(t, 6.5, byID["erosita-1"].Dec, 1e-9)
		require.InDelta(t, 7.5, byID["erosita-2"].Ra, 1e-9)
		require.InDelta(t, 8.5, byID["erosita-2"].Dec, 1e-9)
	})
}

func TestMetadata_FindByID_CaseInsensitiveCatalog(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	ctx := context.Background()
	seedGaia(t, queries, "gaia-1", 0, 1.5, 2.5, repository.Gaia{})

	result, err := service.FindByID(ctx, "gaia-1", "GAIA")
	require.NoError(t, err)
	require.Equal(t, "gaia-1", result.(repository.GetGaiaRow).ID)

	bulk, err := service.BulkFindByID(ctx, []string{"gaia-1"}, "Gaia")
	require.NoError(t, err)
	require.Len(t, bulk.([]repository.BulkGetGaiaRow), 1)
}

func TestMetadata_FindByID_UnknownCatalog(t *testing.T) {
	service, _, _ := newMetadataTestService(t)

	result, err := service.FindByID(context.Background(), "allwise-1", "invalid")
	require.Nil(t, result)

	var validationErr ValidationError
	require.ErrorAs(t, err, &validationErr)
	require.Equal(t, "catalog", validationErr.Field)
	require.Equal(t, "Could not parse field catalog with value invalid: unknown catalog: invalid", err.Error())
}

func TestMetadata_BulkFindByID_UnknownCatalog(t *testing.T) {
	service, _, _ := newMetadataTestService(t)

	result, err := service.BulkFindByID(context.Background(), []string{"allwise-1"}, "invalid")
	require.Nil(t, result)

	var validationErr ValidationError
	require.ErrorAs(t, err, &validationErr)
	require.Equal(t, "catalog", validationErr.Field)
}

func TestMetadata_FindByID_UnknownID(t *testing.T) {
	service, _, _ := newMetadataTestService(t)

	_, err := service.FindByID(context.Background(), "missing", "allwise")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestMetadata_FindByID_SQLInjectionPattern(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	seedAllwise(t, queries, "allwise-1", 0, 12.34, 56.78, repository.Allwise{})

	result, err := service.FindByID(context.Background(), "allwise-1; DROP TABLE allwise", "allwise")
	require.Nil(t, result)

	var validationErr ValidationError
	require.ErrorAs(t, err, &validationErr)
	require.Equal(t, "id", validationErr.Field)
}

func TestMetadata_BulkFindByID_SQLInjectionPattern(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	seedAllwise(t, queries, "allwise-1", 0, 12.34, 56.78, repository.Allwise{})

	result, err := service.BulkFindByID(context.Background(), []string{"allwise-1", "1' OR '1'='1"}, "allwise")
	require.Nil(t, result)

	var validationErr ValidationError
	require.ErrorAs(t, err, &validationErr)
	require.Equal(t, "id", validationErr.Field)
}

func TestMetadata_FindByID_DatabaseError(t *testing.T) {
	service, _, db := newMetadataTestService(t)
	require.NoError(t, db.Close())

	_, err := service.FindByID(context.Background(), "allwise-1", "allwise")
	require.Error(t, err)
}

func TestNewMetadataService_NilQueries(t *testing.T) {
	service, err := NewMetadataService(nil)
	require.Nil(t, service)
	require.Error(t, err)
}

func TestMetadata_GetFromPixelRanges_PerCatalog(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	ctx := context.Background()
	ranges := []healpix.PixelRange{{Start: 10, Stop: 20}}

	// The seeded non-zero fields sit at the beginning, middle, and end of the
	// tables so a column-order mismatch between SELECT <table>.* and the
	// generated model cannot pass unnoticed.
	seedAllwise(t, queries, "allwise-1", 15, 12.5, -3.25, repository.Allwise{
		Cntr:       7,
		W1mpro:     nullFloat64(1.5),
		JM2mass:    nullFloat64(5.0),
		KMsig2mass: nullFloat64(0.7),
	})
	seedGaia(t, queries, "gaia-1", 15, 60.0, -12.5, repository.Gaia{
		SourceID:        nullInt64(424242),
		PhotGMeanFlux:   nullFloat64(1234.5),
		DistanceGspphot: nullFloat64(987.6),
	})
	// The erosita table carries its own coordinates; they are deliberately
	// different from the mastercat ones to pin the envelope behavior below.
	seedErosita(t, queries, "erosita-1", 15, 10.0, 20.0, repository.Erosita{
		Detuid:  nullString("det-1"),
		Skytile: nullInt64(12345),
		Ra:      nullFloat64(100.0),
		Dec:     nullFloat64(50.0),
		MlFlux1: nullFloat64(1.5e-13),
		FlagOpt: nullInt64(9),
	})

	t.Run("allwise", func(t *testing.T) {
		result, err := service.GetFromPixelRanges(ctx, "allwise", ranges)
		require.NoError(t, err)
		require.Len(t, result, 1)

		item := result[0]
		require.Equal(t, "allwise-1", item.ID)
		require.Equal(t, "AllWISE", item.Catalog)
		require.InDelta(t, 12.5, item.Ra, 1e-9)
		require.InDelta(t, -3.25, item.Dec, 1e-9)

		object, ok := item.Object.(repository.Allwise)
		require.True(t, ok)
		require.Equal(t, "allwise-1", object.ID)
		require.Equal(t, int64(7), object.Cntr)
		require.True(t, object.W1mpro.Valid)
		require.InDelta(t, 1.5, object.W1mpro.Float64, 1e-9)
		require.InDelta(t, 5.0, object.JM2mass.Float64, 1e-9)
		require.InDelta(t, 0.7, object.KMsig2mass.Float64, 1e-9)
	})

	t.Run("gaia", func(t *testing.T) {
		result, err := service.GetFromPixelRanges(ctx, "gaia", ranges)
		require.NoError(t, err)
		require.Len(t, result, 1)

		item := result[0]
		require.Equal(t, "gaia-1", item.ID)
		require.Equal(t, "GAIA/DR3", item.Catalog)
		require.InDelta(t, 60.0, item.Ra, 1e-9)
		require.InDelta(t, -12.5, item.Dec, 1e-9)

		object, ok := item.Object.(repository.Gaia)
		require.True(t, ok)
		require.Equal(t, int64(424242), object.SourceID.Int64)
		require.InDelta(t, 1234.5, object.PhotGMeanFlux.Float64, 1e-9)
		require.InDelta(t, 987.6, object.DistanceGspphot.Float64, 1e-9)
	})

	t.Run("erosita envelope uses the mastercat coordinates", func(t *testing.T) {
		result, err := service.GetFromPixelRanges(ctx, "erosita", ranges)
		require.NoError(t, err)
		require.Len(t, result, 1)

		item := result[0]
		require.Equal(t, "erosita-1", item.ID)
		require.Equal(t, "eROSITA", item.Catalog)
		// The envelope carries the mastercat pair (10.0, 20.0), like every
		// other catalog, while the row object keeps the table's own
		// (nullable) coordinates (100.0, 50.0).
		require.InDelta(t, 10.0, item.Ra, 1e-9)
		require.InDelta(t, 20.0, item.Dec, 1e-9)

		object, ok := item.Object.(repository.Erosita)
		require.True(t, ok)
		require.Equal(t, "det-1", object.Detuid.String)
		require.Equal(t, int64(12345), object.Skytile.Int64)
		require.InDelta(t, 100.0, object.Ra.Float64, 1e-9)
		require.InDelta(t, 50.0, object.Dec.Float64, 1e-9)
		require.InDelta(t, 1.5e-13, object.MlFlux1.Float64, 1e-20)
		require.Equal(t, int64(9), object.FlagOpt.Int64)
	})
}

func TestMetadata_GetFromPixelRanges_CaseInsensitiveCatalog(t *testing.T) {
	service, queries, _ := newMetadataTestService(t)
	seedAllwise(t, queries, "allwise-1", 15, 12.5, -3.25, repository.Allwise{})

	result, err := service.GetFromPixelRanges(context.Background(), "ALLWISE", []healpix.PixelRange{{Start: 10, Stop: 20}})
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "allwise-1", result[0].ID)
}

// TestMetadata_GetFromPixelRanges_NoUsableRanges closes the database before
// querying: returning a result without error proves no query was issued for
// empty, nil, or fully degenerate range inputs.
func TestMetadata_GetFromPixelRanges_NoUsableRanges(t *testing.T) {
	service, _, db := newMetadataTestService(t)
	require.NoError(t, db.Close())

	cases := []struct {
		name   string
		ranges []healpix.PixelRange
	}{
		{name: "nil input"},
		{name: "empty input", ranges: []healpix.PixelRange{}},
		{name: "all ranges degenerate", ranges: []healpix.PixelRange{
			{Start: 5, Stop: 5},
			{Start: 9, Stop: 3},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.GetFromPixelRanges(context.Background(), "allwise", tc.ranges)
			require.NoError(t, err)
			require.Empty(t, result)
		})
	}
}

// TestMetadata_GetFromPixelRanges_UnknownCatalog closes the database for the
// same reason as above: an unregistered name must contribute nothing without
// touching SQL.
func TestMetadata_GetFromPixelRanges_UnknownCatalog(t *testing.T) {
	service, _, db := newMetadataTestService(t)
	require.NoError(t, db.Close())

	result, err := service.GetFromPixelRanges(
		context.Background(),
		"notacatalog",
		[]healpix.PixelRange{{Start: 10, Stop: 20}},
	)
	require.NoError(t, err)
	require.Empty(t, result)
}

func newMetadataTestService(t *testing.T) (*MetadataService, *repository.Queries, *sql.DB) {
	t.Helper()

	dbFile := filepath.Join(t.TempDir(), "metadata.db")
	mig, err := migrate.New(fmt.Sprintf("file://%s", metadataMigrationsDir(t)), fmt.Sprintf("sqlite3://%s", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = mig.Close() })
	require.NoError(t, mig.Up())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.Ping())

	queries := repository.New(db)
	service, err := NewMetadataService(queries)
	require.NoError(t, err)

	return service, queries, db
}

func metadataMigrationsDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(file), "..", "..", "db", "migrations")
}

func seedObject(t *testing.T, q *repository.Queries, id string, ipix int64, ra, dec float64, cat string) {
	t.Helper()

	err := q.InsertObject(context.Background(), repository.InsertObjectParams{
		ID:   id,
		Ipix: ipix,
		Ra:   ra,
		Dec:  dec,
		Cat:  cat,
	})
	require.NoError(t, err)
}

func seedAllwise(t *testing.T, q *repository.Queries, id string, ipix int64, ra, dec float64, metadata repository.Allwise) {
	t.Helper()

	seedObject(t, q, id, ipix, ra, dec, "allwise")
	metadata.ID = id
	require.NoError(t, q.InsertAllwise(context.Background(), repository.InsertAllwiseParams(metadata)))
}

func seedGaia(t *testing.T, q *repository.Queries, id string, ipix int64, ra, dec float64, metadata repository.Gaia) {
	t.Helper()

	seedObject(t, q, id, ipix, ra, dec, "gaia")
	metadata.ID = id
	require.NoError(t, q.InsertGaia(context.Background(), repository.InsertGaiaParams(metadata)))
}

func seedErosita(t *testing.T, q *repository.Queries, id string, ipix int64, ra, dec float64, metadata repository.Erosita) {
	t.Helper()

	seedObject(t, q, id, ipix, ra, dec, "erosita")
	metadata.ID = id
	require.NoError(t, q.InsertErosita(context.Background(), repository.InsertErositaParams(metadata)))
}

func nullString(s string) repository.NullString {
	return repository.NullString{NullString: sql.NullString{String: s, Valid: true}}
}

func nullInt64(i int64) repository.NullInt64 {
	return repository.NullInt64{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}

func nullFloat64(f float64) repository.NullFloat64 {
	return repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: f, Valid: true}}
}
