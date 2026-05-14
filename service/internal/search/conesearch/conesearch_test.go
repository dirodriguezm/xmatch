package conesearch

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
)

func TestConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
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
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("Test error"))
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
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(vlassObjects, nil).Once()
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(ztfObjects, nil).Once()
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

func TestBulkConesearch(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
	catalogs := []repository.Catalog{{Name: "vlass", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithCatalogs(catalogs))
	require.NoError(t, err)

	type testCase struct {
		ra        []float64
		dec       []float64
		radius    float64
		nneighbor int
		expected  []string
	}

	testCases := []testCase{
		{ra: []float64{1}, dec: []float64{1}, radius: 1, nneighbor: 100, expected: []string{"A"}},
		{ra: []float64{10}, dec: []float64{10}, radius: 1, nneighbor: 100, expected: []string{"B"}},
		{ra: []float64{1, 2, 3}, dec: []float64{1, 2, 3}, radius: 1, nneighbor: 100, expected: []string{"A"}},
		{ra: []float64{1, 10}, dec: []float64{1, 10}, radius: 1, nneighbor: 100, expected: []string{"A", "B"}},
	}

	for _, tc := range testCases {
		result, err := service.BulkConesearch(tc.ra, tc.dec, tc.radius, tc.nneighbor, "all", 2, 1)
		require.NoError(t, err)
		repo.AssertExpectations(t)

		require.Lenf(t, result, len(tc.expected), "test case: %v", tc)
		for i := range result {
			for j := range result[i].Data {
				require.Contains(t, tc.expected, result[i].Data[j].ID, "test case: %v", tc)
			}
		}
	}
}

func TestBulkConesearch_WithRepositoryError(t *testing.T) {
	repo := NewMockMastercatStore(t)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(nil, errors.New("repository error"))
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
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err = db.Exec(`
		CREATE TABLE mastercat (
			id text not null,
			ipix bigint not null,
			ra double precision not null,
			dec double precision not null,
			cat text not null,
			PRIMARY KEY (id, cat)
		);
		CREATE TABLE allwise (
			id text not null,
			cntr bigint not null,
			w1mpro double precision,
			w1sigmpro double precision,
			w2mpro double precision,
			w2sigmpro double precision,
			w3mpro double precision,
			w3sigmpro double precision,
			w4mpro double precision,
			w4sigmpro double precision,
			J_m_2mass double precision,
			J_msig_2mass double precision,
			H_m_2mass double precision,
			H_msig_2mass double precision,
			K_m_2mass double precision,
			K_msig_2mass double precision,
			PRIMARY KEY (id)
		);
	`)
	require.NoError(t, err)
	return repository.New(db)
}

func TestConesearch_WithMetadata(t *testing.T) {
	metadataRepo := newAllwiseMetadataRepo(t)
	ctx := context.Background()
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)
	require.NoError(t, metadataRepo.InsertMastercat(ctx, repository.Mastercat{ID: "A", Ipix: mapper.PixelAt(healpix.RADec(1, 1)), Ra: 1, Dec: 1, Cat: "allwise"}))
	require.NoError(t, metadataRepo.InsertMastercat(ctx, repository.Mastercat{ID: "B", Ipix: mapper.PixelAt(healpix.RADec(10, 10)), Ra: 10, Dec: 10, Cat: "allwise"}))
	require.NoError(t, metadataRepo.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: "A"}))
	require.NoError(t, metadataRepo.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: "B"}))
	resolver := catalog.NewResolver(metadataRepo)

	repo := NewMockMastercatStore(t)
	catalogs := []repository.Catalog{{Name: "allwise", Nside: 18}}
	service, err := NewConesearchService(WithScheme(healpix.Nest), WithMastercatStore(repo), WithResolver(resolver), WithCatalogs(catalogs))
	require.NoError(t, err)

	result, err := service.FindMetadataByConesearch(1, 1, 1, 1, "allwise")
	require.NoError(t, err)

	require.Len(t, result, 1)
	require.Equal(t, result[0].Data[0].ID, "A")
}

func FuzzConesearch(f *testing.F) {
	objects := []repository.Mastercat{
		{ID: "A", Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "B", Ra: 10, Dec: 10, Cat: "vlass"},
	}
	repo := NewMockMastercatStore(f)
	repo.On("FindObjects", mock.Anything, mock.Anything).Return(objects, nil)
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
