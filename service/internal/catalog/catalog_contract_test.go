package catalog_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

type catalogContractCase struct {
	name             string
	displayName      string
	raw              any
	expectedID       string
	expectedRA       float64
	expectedDec      float64
	assertRecordType func(*testing.T, any)
	assertMetadata   func(*testing.T, any)
	assertGetByID    func(*testing.T, any)
	assertBulkGet    func(*testing.T, any)
	assertObject     func(*testing.T, any)
}

func TestRealCatalogAdaptersSatisfyContract(t *testing.T) {
	allwiseID := "allwise-001"
	allwiseCntr := int64(42)
	allwiseRA := 11.1
	allwiseDec := -22.2
	allwiseW1 := 12.3

	contractCases := []catalogContractCase{
		{
			name:        "allwise",
			displayName: "AllWISE",
			raw: allwise.InputSchema{
				Source_id: &allwiseID,
				Cntr:      &allwiseCntr,
				Ra:        &allwiseRA,
				Dec:       &allwiseDec,
				W1mpro:    &allwiseW1,
			},
			expectedID:  allwiseID,
			expectedRA:  allwiseRA,
			expectedDec: allwiseDec,
			assertRecordType: func(t *testing.T, record any) {
				t.Helper()
				require.IsType(t, repository.Allwise{}, record)
			},
			assertMetadata: func(t *testing.T, metadata any) {
				t.Helper()
				row := metadata.(repository.Allwise)
				require.Equal(t, allwiseID, row.ID)
				require.Equal(t, allwiseCntr, row.Cntr)
				require.True(t, row.W1mpro.Valid)
				require.Equal(t, allwiseW1, row.W1mpro.Float64)
			},
			assertGetByID: func(t *testing.T, got any) {
				t.Helper()
				row := got.(repository.GetAllwiseRow)
				require.Equal(t, allwiseID, row.ID)
				require.Equal(t, allwiseCntr, row.Cntr)
				require.Equal(t, allwiseRA, row.Ra)
				require.Equal(t, allwiseDec, row.Dec)
			},
			assertBulkGet: func(t *testing.T, got any) {
				t.Helper()
				rows := got.([]repository.BulkGetAllwiseRow)
				require.Len(t, rows, 1)
				require.Equal(t, allwiseID, rows[0].ID)
			},
			assertObject: func(t *testing.T, object any) {
				t.Helper()
				row := object.(repository.Allwise)
				require.Equal(t, allwiseID, row.ID)
				require.Equal(t, allwiseCntr, row.Cntr)
			},
		},
		{
			name:        "gaia",
			displayName: "GAIA/DR3",
			raw: gaia.InputSchema{
				Designation:    "Gaia DR3 123",
				RA:             33.3,
				Dec:            44.4,
				PhotGMeanFlux:  1000.5,
				PhotBpMeanMag:  16.1,
				PhotRpMeanFlux: 600.7,
			},
			expectedID:  "Gaia DR3 123",
			expectedRA:  33.3,
			expectedDec: 44.4,
			assertRecordType: func(t *testing.T, record any) {
				t.Helper()
				require.IsType(t, repository.Gaia{}, record)
			},
			assertMetadata: func(t *testing.T, metadata any) {
				t.Helper()
				row := metadata.(repository.Gaia)
				require.Equal(t, "Gaia DR3 123", row.ID)
				require.True(t, row.PhotGMeanFlux.Valid)
				require.Equal(t, 1000.5, row.PhotGMeanFlux.Float64)
				require.InDelta(t, 16.1, row.PhotBpMeanMag.Float64, 0.0001)
			},
			assertGetByID: func(t *testing.T, got any) {
				t.Helper()
				row := got.(repository.GetGaiaRow)
				require.Equal(t, "Gaia DR3 123", row.ID)
				require.Equal(t, 33.3, row.Ra)
				require.Equal(t, 44.4, row.Dec)
			},
			assertBulkGet: func(t *testing.T, got any) {
				t.Helper()
				rows := got.([]repository.BulkGetGaiaRow)
				require.Len(t, rows, 1)
				require.Equal(t, "Gaia DR3 123", rows[0].ID)
			},
			assertObject: func(t *testing.T, object any) {
				t.Helper()
				row := object.(repository.Gaia)
				require.Equal(t, "Gaia DR3 123", row.ID)
				require.Equal(t, 1000.5, row.PhotGMeanFlux.Float64)
			},
		},
		{
			name:        "erosita",
			displayName: "eROSITA",
			raw: erosita.InputSchema{
				IAUNAME: "eROSITA J123",
				DETUID:  "det-123",
				SKYTILE: 17,
				ID_SRC:  99,
				UID:     1001,
				RA:      55.5,
				DEC:     -12.3,
			},
			expectedID:  "eROSITA J123",
			expectedRA:  55.5,
			expectedDec: -12.3,
			assertRecordType: func(t *testing.T, record any) {
				t.Helper()
				require.IsType(t, repository.Erosita{}, record)
			},
			assertMetadata: func(t *testing.T, metadata any) {
				t.Helper()
				row := metadata.(repository.Erosita)
				require.Equal(t, "eROSITA J123", row.ID)
				require.Equal(t, "det-123", row.Detuid.String)
				require.Equal(t, int64(17), row.Skytile.Int64)
				require.Equal(t, int64(99), row.IDSrc.Int64)
			},
			assertGetByID: func(t *testing.T, got any) {
				t.Helper()
				row := got.(repository.GetErositaRow)
				require.Equal(t, "eROSITA J123", row.ID)
				require.Equal(t, "det-123", row.Detuid.String)
				require.Equal(t, 55.5, row.Ra.Float64)
				require.Equal(t, -12.3, row.Dec.Float64)
			},
			assertBulkGet: func(t *testing.T, got any) {
				t.Helper()
				rows := got.([]repository.BulkGetErositaRow)
				require.Len(t, rows, 1)
				require.Equal(t, "eROSITA J123", rows[0].ID)
			},
			assertObject: func(t *testing.T, object any) {
				t.Helper()
				row := object.(repository.Erosita)
				require.Equal(t, "eROSITA J123", row.ID)
				require.Equal(t, "det-123", row.Detuid.String)
			},
		},
	}

	for _, tc := range contractCases {
		t.Run(tc.name, func(t *testing.T) {
			queries := migratedCatalogTestDB(t)
			adapter, err := catalog.NewResolver(queries).Get(tc.name)
			require.NoError(t, err)

			require.Equal(t, tc.name, adapter.Name())
			require.IsType(t, tc.raw, adapter.NewRawRecord())
			tc.assertRecordType(t, adapter.NewMetadataRecord())

			ra, dec, err := adapter.GetCoordinates(tc.raw)
			require.NoError(t, err)
			require.Equal(t, tc.expectedRA, ra)
			require.Equal(t, tc.expectedDec, dec)

			mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
			require.NoError(t, err)
			ipix := mapper.PixelAt(healpix.RADec(tc.expectedRA, tc.expectedDec))

			mastercat, err := adapter.ConvertToMastercat(tc.raw, ipix)
			require.NoError(t, err)
			require.Equal(t, repository.Mastercat{ID: tc.expectedID, Ipix: ipix, Ra: tc.expectedRA, Dec: tc.expectedDec, Cat: tc.name}, mastercat)

			metadata, err := adapter.ConvertToMetadataFromRaw(tc.raw)
			require.NoError(t, err)
			tc.assertMetadata(t, metadata)

			ctx := context.Background()
			require.NoError(t, queries.InsertObject(ctx, repository.InsertObjectParams(mastercat)))
			require.NoError(t, adapter.BulkInsertMetadata(ctx, []any{metadata}))

			byID, err := adapter.GetByID(ctx, tc.expectedID)
			require.NoError(t, err)
			tc.assertGetByID(t, byID)

			bulkByID, err := adapter.BulkGetByID(ctx, []string{tc.expectedID})
			require.NoError(t, err)
			tc.assertBulkGet(t, bulkByID)

			fromPixels, err := adapter.GetFromPixels(ctx, []int64{ipix})
			require.NoError(t, err)
			require.Len(t, fromPixels, 1)
			require.Equal(t, tc.expectedID, fromPixels[0].ID)
			require.Equal(t, tc.displayName, fromPixels[0].Catalog)
			require.Equal(t, tc.expectedRA, fromPixels[0].Ra)
			require.Equal(t, tc.expectedDec, fromPixels[0].Dec)
			tc.assertObject(t, fromPixels[0].Object)
		})
	}
}

func migratedCatalogTestDB(t *testing.T) *repository.Queries {
	t.Helper()

	dbFile := filepath.Join(t.TempDir(), "catalog-contract.db")
	mig, err := migrate.New(fmt.Sprintf("file://%s", catalogMigrationsDir(t)), fmt.Sprintf("sqlite3://%s", dbFile))
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

	return repository.New(db)
}

func catalogMigrationsDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(file), "..", "db", "migrations")
}
