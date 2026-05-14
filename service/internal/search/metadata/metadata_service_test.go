package metadata

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
)

func newAllwiseMetadataService(t *testing.T) (*MetadataService, *repository.Queries) {
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

	queries := repository.New(db)
	return &MetadataService{resolver: catalog.NewResolver(queries)}, queries
}

func insertAllwiseMetadata(t *testing.T, queries *repository.Queries, id string, ra, dec float64) {
	t.Helper()
	ctx := context.Background()
	require.NoError(t, queries.InsertMastercat(ctx, repository.Mastercat{ID: id, Ipix: 1, Ra: ra, Dec: dec, Cat: "allwise"}))
	require.NoError(t, queries.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: id}))
}

func TestMetadata_ValidateCatalog(t *testing.T) {
	m := &MetadataService{resolver: catalog.NewResolver(nil)}

	err := m.validateCatalog("allwise")
	require.Nil(t, err)
	err = m.validateCatalog("AllWise")
	require.Nil(t, err)

	err = m.validateCatalog("gaia")
	require.Nil(t, err)
	err = m.validateCatalog("GAIA")
	require.Nil(t, err)

	err = m.validateCatalog("erosita")
	require.Nil(t, err)
	err = m.validateCatalog("Erosita")
	require.Nil(t, err)

	err = m.validateCatalog("invalid")
	require.NotNil(t, err)
	require.Equal(t, "Could not parse field catalog with value invalid: unknown catalog: invalid", err.Error())
}

func TestMetadata_FindByID(t *testing.T) {
	m, queries := newAllwiseMetadataService(t)
	insertAllwiseMetadata(t, queries, "allwise1", 12.34, 56.78)

	result, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.Nil(t, err)
	require.Equal(t, "allwise1", result.(repository.GetAllwiseRow).ID)
}

func TestMetadata_BulkFindByID(t *testing.T) {
	m, queries := newAllwiseMetadataService(t)
	insertAllwiseMetadata(t, queries, "allwise1", 12.34, 56.78)
	insertAllwiseMetadata(t, queries, "allwise2", 23.45, 67.89)

	result, err := m.BulkFindByID(context.Background(), []string{"allwise1", "allwise2"}, "allwise")
	require.Nil(t, err)
	expectedIds := []string{"allwise1", "allwise2"}
	for i := 0; i < len(result.([]repository.BulkGetAllwiseRow)); i++ {
		require.Equal(t, expectedIds[i], result.([]repository.BulkGetAllwiseRow)[i].ID)
	}
}

func TestMetadata_Bulk_EmptyResult(t *testing.T) {
	m, _ := newAllwiseMetadataService(t)

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.NotNil(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestMetadata_SomeDBError(t *testing.T) {
	m, queries := newAllwiseMetadataService(t)
	require.NoError(t, queries.GetDbInstance().Close())

	_, err := m.FindByID(context.Background(), "allwise1", "allwise")
	require.NotNil(t, err)
}
