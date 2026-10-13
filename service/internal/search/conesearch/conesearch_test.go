package conesearch

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1, "all")
	require.NoError(t, err)
	repo.AssertExpectations(t)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].ID, "A")
}

func TestConesearch_WithRepositoryError(t *testing.T) {
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(nil, errors.New("Test error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	_, err = service.Conesearch(1, 1, 1, 1, "all")
	repo.AssertExpectations(t)
	if assert.Error(t, err) {
		require.Equal(t, errors.New("Test error"), err)
	}
}

func TestConesearch_WithMultipleMappers(t *testing.T) {
	vlassObjects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	ztfObjects := []repository.Mastercat{
		{ID: "ZTFA", Ra: 1, Dec: 1, Cat: "ztf"},
	}
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(vlassObjects, nil).Once()
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(ztfObjects, nil).Once()
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}, {Name: "ztf", Nside: 12}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 2, "all")
	repo.AssertExpectations(t)

	require.Len(t, result, 2)
	ids := make([]string, 2)
	cats := make([]string, 2)
	for i := range result {
		for j := range result[i].Data {
			ids[i] = result[i].Data[j].ID
			cats[i] = result[i].Data[j].Cat
		}
	}
	require.Subset(t, ids, []string{"A", "ZTFA"})
	require.Subset(t, cats, []string{"vlass", "ztf"})
}

func TestConesearch_WithUnknownCatalogReturnsNoResults(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
	}
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.Conesearch(1, 1, 1, 1, "unknown")
	require.NoError(t, err)
	repo.AssertExpectations(t)
	require.Empty(t, result)
}

func TestBulkConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 1.1, Dec: 1, Cat: "vlass"},
		{ID: "C", Ra: 1.2, Dec: 1, Cat: "ztf"},
		{ID: "D", Ra: 10, Dec: 10, Cat: "ztf"},
	}
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	type expectedRow struct {
		id             string
		index          int
		catalog        string
		total          int
		totalInCatalog int
	}

	type testCase struct {
		ra        []float64
		dec       []float64
		radius    float64
		nneighbor int
		expected  []expectedRow
	}

	testCases := []testCase{
		{
			ra: []float64{1}, dec: []float64{1}, radius: 3600, nneighbor: 100,
			expected: []expectedRow{
				{id: "A", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "B", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "C", index: 0, catalog: "ztf", total: 3, totalInCatalog: 1},
			},
		},
		{
			// The cut keeps only the nearest object, from vlass; the fully
			// cut ztf catalog still contributes to the position's total.
			ra: []float64{1}, dec: []float64{1}, radius: 3600, nneighbor: 1,
			expected: []expectedRow{
				{id: "A", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
			},
		},
		{
			ra: []float64{10}, dec: []float64{10}, radius: 3600, nneighbor: 100,
			expected: []expectedRow{
				{id: "D", index: 0, catalog: "ztf", total: 1, totalInCatalog: 1},
			},
		},
		{
			// Positions 1 and 2 have no in-radius objects, so they have no
			// rows: absence means zero matches for that position.
			ra: []float64{1, 2, 3}, dec: []float64{1, 2, 3}, radius: 3600, nneighbor: 100,
			expected: []expectedRow{
				{id: "A", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "B", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "C", index: 0, catalog: "ztf", total: 3, totalInCatalog: 1},
			},
		},
		{
			ra: []float64{1, 10}, dec: []float64{1, 10}, radius: 3600, nneighbor: 100,
			expected: []expectedRow{
				{id: "A", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "B", index: 0, catalog: "vlass", total: 3, totalInCatalog: 2},
				{id: "C", index: 0, catalog: "ztf", total: 3, totalInCatalog: 1},
				{id: "D", index: 1, catalog: "ztf", total: 1, totalInCatalog: 1},
			},
		},
	}

	for _, tc := range testCases {
		result, err := service.BulkConesearch(tc.ra, tc.dec, tc.radius, tc.nneighbor, "all", 2, 1)
		require.NoError(t, err)
		repo.AssertExpectations(t)

		require.Lenf(t, result, len(tc.expected), "test case: %v", tc)

		matched := make([]bool, len(tc.expected))
		for _, row := range result {
			for j := range row.Data {
				found := -1
				for k, expected := range tc.expected {
					if row.Index == expected.index && row.Catalog == expected.catalog && row.Data[j].ID == expected.id {
						found = k
						break
					}
				}
				require.Truef(t, found >= 0, "unexpected row %q for position %d in test case %v", row.Data[j].ID, row.Index, tc)
				require.Equalf(t, tc.expected[found].total, row.Total, "total for position %d in test case %v", row.Index, tc)
				require.Equalf(t, tc.expected[found].totalInCatalog, row.TotalInCatalog, "total_in_catalog for (position %d, catalog %q) in test case %v", row.Index, row.Catalog, tc)
				matched[found] = true
			}
		}
		for k, m := range matched {
			require.Truef(t, m, "expected row %+v was not returned, test case %v", tc.expected[k], tc)
		}

		// Positions with no expected rows must have no rows at all.
		expectedIndexes := make(map[int]bool)
		for _, expected := range tc.expected {
			expectedIndexes[expected.index] = true
		}
		for i := range tc.ra {
			if expectedIndexes[i] {
				continue
			}
			for _, row := range result {
				require.NotEqualf(t, i, row.Index, "unexpected row for position %d in test case %v", i, tc)
			}
		}
	}
}

func TestBulkConesearch_CollapsesDuplicateCandidates(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 1.1, Dec: 1, Cat: "vlass"},
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
	}
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.BulkConesearch([]float64{1}, []float64{1}, 3600, 100, "all", 1, 1)
	require.NoError(t, err)
	repo.AssertExpectations(t)

	// The same object appears twice in the candidate set for the input
	// position, so it must be returned once; distinct objects are kept.
	require.Len(t, result, 2)
	ids := make([]string, 0, len(result))
	for _, r := range result {
		require.Len(t, r.Data, 1)
		require.Equal(t, 0, r.Index)
		ids = append(ids, r.Data[0].ID)
	}
	require.ElementsMatch(t, []string{"A", "B"}, ids)
}

func TestBulkConesearch_WithRepositoryError(t *testing.T) {
	repo := repository.NewMockMastercatReader(t)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	_, err = service.BulkConesearch([]float64{1, 10}, []float64{1, 10}, 1, 100, "all", 2, 1)
	repo.AssertExpectations(t)
	require.Error(t, err)
	require.Equal(t, "repository error", err.Error())
}

func newAllwiseMetadataRepo(t *testing.T) *repository.Queries {
	t.Helper()

	dbFile := filepath.Join(t.TempDir(), "conesearch-metadata.db")
	mig, err := migrate.New(fmt.Sprintf("file://%s", conesearchMigrationsDir(t)), fmt.Sprintf("sqlite3://%s", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = mig.Close() })
	require.NoError(t, mig.Up())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.Ping())

	return repository.New(db)
}

func conesearchMigrationsDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Join(filepath.Dir(file), "..", "..", "db", "migrations")
}

func TestConesearch_WithMetadata(t *testing.T) {
	metadataRepo := newAllwiseMetadataRepo(t)
	ctx := context.Background()
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)
	require.NoError(t, metadataRepo.InsertObject(ctx, repository.InsertObjectParams{ID: "A", Ipix: mapper.PixelAt(healpix.RADec(1, 1)), Ra: 1, Dec: 1, Cat: "allwise"}))
	require.NoError(t, metadataRepo.InsertObject(ctx, repository.InsertObjectParams{ID: "B", Ipix: mapper.PixelAt(healpix.RADec(10, 10)), Ra: 10, Dec: 10, Cat: "allwise"}))
	require.NoError(t, metadataRepo.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: "A"}))
	require.NoError(t, metadataRepo.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: "B"}))

	metadataService, err := metadata.NewMetadataService(metadataRepo)
	require.NoError(t, err)

	repo := repository.NewMockMastercatReader(t)
	catalogs := []repository.Catalog{{Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithMetadataService(metadataService), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "allwise")
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].ID, "A")
}

// TestConesearch_WithMetadata_ErositaUsesMastercatCoordinates pins the
// coordinate source of erosita metadata envelopes: the radius cutoff, ranking,
// and reported distance use the mastercat pair the pixel filter is built from,
// not the table's own (nullable) ra/dec columns.
func TestConesearch_WithMetadata_ErositaUsesMastercatCoordinates(t *testing.T) {
	metadataRepo := newAllwiseMetadataRepo(t)
	ctx := context.Background()
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)
	require.NoError(t, metadataRepo.InsertObject(ctx, repository.InsertObjectParams{
		ID:   "E",
		Ipix: mapper.PixelAt(healpix.RADec(10, 10)),
		Ra:   10,
		Dec:  10,
		Cat:  "erosita",
	}))
	// The table's own coordinates are far outside the query radius: only an
	// envelope built from the mastercat pair can match at (10, 10).
	require.NoError(t, metadataRepo.InsertErosita(ctx, repository.InsertErositaParams{
		ID:  "E",
		Ra:  repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 100, Valid: true}},
		Dec: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 50, Valid: true}},
	}))

	metadataService, err := metadata.NewMetadataService(metadataRepo)
	require.NoError(t, err)

	repo := repository.NewMockMastercatReader(t)
	catalogs := []repository.Catalog{{Name: "erosita", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithMetadataService(metadataService), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(10, 10, 1, 1, "erosita")
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.Equal(t, "eROSITA", result[0].Catalog)
	require.Len(t, result[0].Data, 1)
	require.Equal(t, "E", result[0].Data[0].ID)
	// Distance ~0 proves the envelope coordinates are the mastercat pair.
	require.InDelta(t, 0, result[0].Data[0].Distance, 1e-9)

	// The serialized payload publishes that same mastercat pair under ra/dec;
	// the erosita table's own (100, 50) must not leak through.
	encoded, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded []struct {
		Data []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Len(t, decoded, 1)
	require.Len(t, decoded[0].Data, 1)
	ra, ok := decoded[0].Data[0]["ra"].(float64)
	require.True(t, ok)
	dec, ok := decoded[0].Data[0]["dec"].(float64)
	require.True(t, ok)
	require.InDelta(t, 10, ra, 1e-9)
	require.InDelta(t, 10, dec, 1e-9)
	require.NotContains(t, decoded[0].Data[0], "ra_2")
	require.NotContains(t, decoded[0].Data[0], "dec_2")
}

// TestConesearch_WithMetadata_SkipsUnregisteredCatalog pins the silent skip:
// expanding catalog=all over the catalogs-table names contributes nothing for
// a name the metadata service does not dispatch, without failing the query.
func TestConesearch_WithMetadata_SkipsUnregisteredCatalog(t *testing.T) {
	metadataRepo := newAllwiseMetadataRepo(t)
	ctx := context.Background()
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)
	require.NoError(t, metadataRepo.InsertObject(ctx, repository.InsertObjectParams{ID: "A", Ipix: mapper.PixelAt(healpix.RADec(1, 1)), Ra: 1, Dec: 1, Cat: "allwise"}))
	require.NoError(t, metadataRepo.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: "A"}))

	metadataService, err := metadata.NewMetadataService(metadataRepo)
	require.NoError(t, err)

	repo := repository.NewMockMastercatReader(t)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}, {Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithMetadataService(metadataService), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "all")
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.Equal(t, "AllWISE", result[0].Catalog)
	require.Len(t, result[0].Data, 1)
	require.Equal(t, "A", result[0].Data[0].ID)
}

func TestConesearch_WithMetadataUnknownCatalogReturnsNoResults(t *testing.T) {
	metadataRepo := newAllwiseMetadataRepo(t)

	metadataService, err := metadata.NewMetadataService(metadataRepo)
	require.NoError(t, err)

	repo := repository.NewMockMastercatReader(t)
	catalogs := []repository.Catalog{{Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithMetadataService(metadataService), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "unknown")
	require.NoError(t, err)
	require.Empty(t, result)
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := repository.NewMockMastercatReader(f)
	repo.On("FindObjectsInPixelRanges", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(f, err)

	f.Add(float64(1), float64(1), float64(1), int(1))
	f.Fuzz(func(t *testing.T, ra float64, dec float64, radius float64, nneighbor int) {
		_, err := service.Conesearch(ra, dec, radius, nneighbor, "all")
		if err == nil {
			repo.AssertExpectations(t)
		}
	})
}
