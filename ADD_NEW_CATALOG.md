# Adding a New Catalog to the ALeRCE xmatch Service

This guide describes the end-state path for adding a new astronomical catalog to the xmatch service. There is no adapter, resolver, factory registration, or blank import: the indexer dispatches on the configured catalog name with one `switch` case in `service/cmd/start_catalog_indexer.go`, and a catalog's behavior is a small package whose `InputSchema` implements a few methods the pipeline calls. The HTTP metadata API dispatches through switch cases in `service/internal/search/metadata/metadata_service.go`.

Use `newcatalog` below as the placeholder catalog name.

## Overview

A catalog needs these pieces:

1. A metadata table migration under `service/internal/db/migrations/`.
2. SQL queries in `service/internal/db/query.sql`: the metadata insert plus the by-id and bulk lookups the HTTP API uses.
3. SQLC config updates in `service/internal/db/sqlc.yaml` only for type exceptions the existing `db_type` rules do not cover (Step 3).
4. Regenerated SQLC repository code.
5. A catalog package under `service/internal/catalog/newcatalog/` whose `InputSchema` implements `GetCoordinates`, `MastercatID`, and `ConvertToMetadata`.
6. An `Insert` method for the generated metadata model in `service/internal/repository/insert.go`.
7. One case in the catalog switch in `service/cmd/start_catalog_indexer.go`.
8. Three cases in the HTTP metadata service switches (by-id, bulk, pixel ranges).
9. Unit and round-trip tests in the catalog package.
10. A config file passed with `CONFIG_PATH` when running the indexer.

The compiler carries the catalog through the pipeline: the switch instantiates the generic readers, indexers, writers, and pipeline with the catalog's `InputSchema` and metadata model, and each stage constrains its type parameter to the methods it calls, so a missing or mistyped method is a build error rather than a runtime assertion.

## Step 1: Add Database Migrations

Create the next numbered migration files in `service/internal/db/migrations/`. The current last catalog migration is `005_erosita`, so the next one should be `006_newcatalog` unless new migrations have been added since this guide was written.

```bash
touch service/internal/db/migrations/006_newcatalog.up.sql
touch service/internal/db/migrations/006_newcatalog.down.sql
```

Metadata tables store catalog-specific fields. Coordinates and HEALPix pixels are stored in `mastercat`, so do not add `ipix` to the catalog table unless the catalog itself has an unrelated pixel column.

Example `006_newcatalog.up.sql`:

```sql
CREATE TABLE IF NOT EXISTS newcatalog (
    id text NOT NULL,
    mag_g double precision,
    mag_r double precision,
    flag smallint,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_newcatalog_mag_g ON newcatalog(mag_g);
```

Example `006_newcatalog.down.sql`:

```sql
DROP TABLE IF EXISTS newcatalog;
```

## Step 2: Add SQLC Queries

Edit `service/internal/db/query.sql` and add the metadata insert and the two lookups the HTTP metadata switches call. Both lookups join `mastercat` so responses carry the object coordinates; the join key is the catalog's object identifier (the same value `MastercatID` returns).

```sql
-- name: InsertNewcatalog :exec
INSERT INTO newcatalog (
    id, mag_g, mag_r, flag
) VALUES (
    ?, ?, ?, ?
);

-- name: GetNewcatalog :one
SELECT newcatalog.*, mastercat.ra, mastercat.dec
FROM newcatalog
JOIN mastercat ON mastercat.id = newcatalog.id
WHERE newcatalog.id = ?;

-- name: BulkGetNewcatalog :many
SELECT newcatalog.*, mastercat.ra, mastercat.dec
FROM newcatalog
JOIN mastercat ON mastercat.id = newcatalog.id
WHERE newcatalog.id IN (sqlc.slice(id));
```

Metadata cone-search candidates are selected by pixel range through the generic `repository.QueryMetadataFromPixelRanges` primitive (Step 6), so no per-catalog pixel query is needed.

## Step 3: Update SQLC Configuration

Edit `service/internal/db/sqlc.yaml` only when the generated types need adjusting:

- Usually nothing. The `db_type` overrides already map nullable `double precision`, `real`, `integer`, `bigint`, `smallint`, and `text` columns to `NullFloat64`, `NullInt64`, and `NullString`, so the Step 1 migration needs no entries. `db_type` matching is an exact comparison against the engine's normalized type name, so `double precision` in a migration is spelled `doubleprecision` (no space) in `sqlc.yaml`, and `nullable: true` matches nullable columns only.
- Add a `column`-scoped `go_type` only for a one-off exception. Types outside the rules above fall back to the stdlib types: `varchar(...)`/`clob` → `sql.NullString`, `numeric`/`decimal` → `sql.NullFloat64`, `boolean` → `sql.NullBool`, and `date`/`datetime`/`timestamp` → `sql.NullTime`. Use a column override to get a project wrapper instead where one exists (`varchar` → `NullString`, `numeric` → `NullFloat64`); `repository` has no `NullBool` or `NullTime`, so those stay stdlib.
- Add a `go_struct_tag` with `json:"..."` only when the generated JSON key must differ from the column name; SQLC already lowercases column names into tags.
- Add a root `overrides.go.rename` entry if SQLC generates an awkward type name.

Do **not** add parquet tags to generated models: nothing exports parquet anymore. Source-reading parquet tags live on the catalog package's `InputSchema` (Step 4).

Example entries for columns the `db_type` rules do not cover (not needed for the Step 1 schema):

```yaml
          - column: "newcatalog.kind"
            go_type:
              type: "NullString"
          - column: "newcatalog.score"
            go_type:
              type: "NullFloat64"
```

Regenerate repository code from the directory containing `sqlc.yaml` with the `sqlc` binary from the dev shell. The pinned nixpkgs ships v1.31.1, which matches the `// sqlc v1.31.1` header of the generated files; do not run `go run github.com/sqlc-dev/sqlc/cmd/sqlc@<version>`, since a different version rewrites that header and churns the generated files:

```bash
cd service/internal/db
sqlc generate
```

This updates `service/internal/repository/models.go` and `service/internal/repository/query.sql.go`.

## Step 4: Add the Catalog Package

Create `service/internal/catalog/newcatalog/newcatalog.go`. A catalog package is just the source schema plus the methods the pipeline calls through its constraints; it imports only `repository` (plus the standard library).

```go
package newcatalog

import (
	"database/sql"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type InputSchema struct {
	ID   string  `parquet:"name=id, type=BYTE_ARRAY"`
	RA   float64 `parquet:"name=ra, type=DOUBLE"`
	Dec  float64 `parquet:"name=dec, type=DOUBLE"`
	MagG float64 `parquet:"name=mag_g, type=DOUBLE"`
	MagR float64 `parquet:"name=mag_r, type=DOUBLE"`
	Flag int64   `parquet:"name=flag, type=INT64"`
}

func (r InputSchema) GetCoordinates() (float64, float64) {
	return r.RA, r.Dec
}

func (r InputSchema) MastercatID() string {
	return r.ID
}

func (r InputSchema) ConvertToMetadata() (repository.Newcatalog, error) {
	return repository.Newcatalog{
		ID:   r.ID,
		MagG: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: r.MagG, Valid: true}},
		MagR: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: r.MagR, Valid: true}},
		Flag: repository.NullInt64{NullInt64: sql.NullInt64{Int64: r.Flag, Valid: true}},
	}, nil
}
```

Notes:

1. `InputSchema` belongs in the catalog package, not `repository`; generated metadata types belong in `repository`. It carries only `parquet:` tags — they describe the source file layout read by the parquet reader. Do not add `json:` tags.
2. `GetCoordinates` returns the source coordinates. The mastercat indexer computes `ipix` from them with the configured nside/ordering and builds the `mastercat` row itself: id from `MastercatID`, coordinates, and `cat` = `strings.ToLower(source.catalog_name)`. There is no per-catalog mastercat conversion method.
3. `MastercatID` is the object identifier and the `mastercat`↔metadata join key. Examples: Gaia's `designation`, eROSITA's `IAUNAME`, AllWISE's `source_id`.
4. `ConvertToMetadata` keeps the `(model, error)` shape so conversion failures still travel the pipeline's error path. Use pointer fields in `InputSchema` when a source format must distinguish null from zero, and encode the catalog's null policy here (for example AllWISE maps missing photometry to `Valid: false`).
5. The metadata model's insert lives on the model itself (Step 5): the writer only requires `repository.Insertable`, so the switch passes no insert function.
6. The CSV reader fills `InputSchema` by struct field order, not by header name. For CSV sources keep fields in source column order. For Parquet/FITS, use pointer fields for nullable columns.
7. If the catalog table also has `ra`/`dec` columns, keep them as ordinary model fields. The HTTP envelope coordinates always come from the trailing `mastercat.ra`/`mastercat.dec` selected by the lookups.

## Step 5: Add the Insert Method and the Indexer Switch Case

Add the metadata model's `Insert` method to `service/internal/repository/insert.go`, next to the generated models and the existing catalog methods:

```go
func (r Newcatalog) Insert(ctx context.Context, q *Queries) error {
	return q.InsertNewcatalog(ctx, InsertNewcatalogParams(r))
}
```

Then edit `service/cmd/start_catalog_indexer.go`: add a direct import for the catalog package and one case to the switch on the lowercased catalog name.

```go
	"github.com/dirodriguezm/xmatch/service/internal/catalog/newcatalog"
```

```go
	case "newcatalog":
		return runCatalogPipeline[repository.Newcatalog, newcatalog.InputSchema](
			ctx, cfg, srcCfg, src, repo,
		)
```

The type arguments are the metadata model first (it must be `repository.Insertable`) and the row type second (it must implement `GetCoordinates`, `MastercatID`, and `ConvertToMetadata`). The case composes the five typed constructors from `pipeline` (mastercat writer, metadata writer, mastercat indexer, metadata indexer, source reader) and runs them to completion. The metadata branch is built only when `catalog_indexer.source.metadata` is true. Any name without a case reaches the default and fails startup with `unknown catalog: <name>` before a source is read.

There is no registration step and no blank import in `service/cmd/main.go`. `CatalogRegister` (the `catalogs` table row used by cone search) is unaffected and keeps working from the source config.

## Step 6: Add the HTTP Metadata Dispatch Cases

The HTTP metadata API does not use catalog packages. `MetadataService` in `service/internal/search/metadata/metadata_service.go` dispatches by catalog name through three switches; add a case for the new catalog to each, with the table name and display-name literals:

```go
// queryCatalog: by-id lookup
case "newcatalog":
	return m.queries.GetNewcatalog(ctx, id)

// bulkQueryCatalog: bulk lookup
case "newcatalog":
	return m.queries.BulkGetNewcatalog(ctx, ids)

// GetFromPixelRanges: pixel-range lookup
case "newcatalog":
	return m.queries.QueryMetadataFromPixelRanges(ctx, "newcatalog", pixelRanges, func(rows *sql.Rows) (repository.Metadata, error) {
		return metadataFromRow[repository.Newcatalog](rows, "NewCatalog")
	})
```

`metadataFromRow` scans the generated model plus the trailing `mastercat.ra`/`mastercat.dec` columns and builds the response envelope (id, display name, coordinates, object). Unknown names keep their existing behavior: by-id and bulk return a `ValidationError` for the `catalog` field, and pixel-range lookups contribute nothing without error.

## Step 7: Add Tests

Put tests in the catalog package itself. The round-trip test replaces the old adapter contract test: it converts a fixture `InputSchema`, inserts through the catalog methods, and reads back with the generated queries against a migrated SQLite database.

`service/internal/catalog/newcatalog/newcatalog_test.go` covers the conversion semantics directly:

```go
func TestConvertToMetadata(t *testing.T) {
	raw := InputSchema{ID: "newcatalog-1", MagG: 18.5}
	got, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	assert.Equal(t, "newcatalog-1", got.ID)
	assert.Equal(t, 18.5, got.MagG.Float64)
}
```

`service/internal/catalog/newcatalog/roundtrip_test.go` pins the stored values:

```go
func TestRoundTrip(t *testing.T) {
	queries := testutils.CatalogDB(t)
	ctx := context.Background()

	raw := InputSchema{ID: "newcatalog-1", RA: 12.3, Dec: -4.5, MagG: 18.5}
	ra, dec := raw.GetCoordinates()
	mastercat := repository.Mastercat{ID: raw.MastercatID(), Ipix: 42, Ra: ra, Dec: dec, Cat: "newcatalog"}
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Mastercat{mastercat}))

	metadata, err := raw.ConvertToMetadata()
	require.NoError(t, err)
	require.NoError(t, repository.BulkInsert(ctx, queries, []repository.Newcatalog{metadata}))

	byID, err := queries.GetNewcatalog(ctx, "newcatalog-1")
	require.NoError(t, err)
	assert.Equal(t, 18.5, byID.MagG.Float64)
}
```

Use a migrated database from `testutils.CatalogDB(t)`; it resolves the migrations relative to its own source file, so tests in any package of the module work without working-directory setup.

## Step 8: Configure and Run the Indexer

Create a config file and run with `CONFIG_PATH`. There is no required `service/configs/` directory; any path is fine. `xwave-run` changes into `service/` before starting the binary, so `CONFIG_PATH` is resolved from there unless you pass an absolute path.

Example `service/newcatalog.yaml`:

```yaml
catalog_indexer:
  database:
    url: "file:dev.db"
  source:
    url: "file:/path/to/newcatalog.parquet"
    type: "parquet"
    catalog_name: "newcatalog"
    nside: 18
    metadata: true
  reader:
    batch_size: 1000
    type: "parquet"
  indexer:
    ordering_scheme: "nested"
    nside: 18
  channel_size: 50000
```

Sqlite output is unconditional: there are no `indexer_writer`/`metadata_writer` keys anymore. Keep `source.nside` and `indexer.nside` aligned. `source.nside` is stored in the `catalogs` table for search, while `indexer.nside` computes `mastercat.ipix`.

For CSV sources, set `source.type` and `reader.type` to `csv`. The CSV reader maps records into `InputSchema` by struct field order, not by header name. If a CSV file has no header row, set `reader.header` in config so the first data row is not consumed as a header. If every file in a `files:` source has its own header row, set `reader.first_line_header: true`.

```bash
CONFIG_PATH=newcatalog.yaml xwave-run indexer
```

## Step 9: Verify

Run these checks from the repository root unless the command changes directories itself:

```bash
cd service/internal/db
sqlc generate
```

```bash
cd service
gofmt -w internal/catalog/newcatalog internal/repository internal/search/metadata cmd
go test ./... -race
```

No mock regeneration is needed for a new catalog: catalog packages have no interfaces and are not mocked. Regenerate mocks (`xwave-mock`) only if you changed a shared interface such as `repository.MastercatReader`.

## Reference Files

- `service/cmd/start_catalog_indexer.go`: catalog switch and typed component composition.
- `service/internal/catalog_indexer/pipeline/pipeline.go`: exported typed constructors (`DefaultMastercatWriter`, `DefaultMetadataWriter`, `DefaultMastercatIndexer`, `DefaultMetadataIndexer`, `DefaultSourceReader`) and `Pipeline`.
- `service/internal/catalog_indexer/indexer/mastercat/indexer.go`: HEALPix pixel computation and mastercat row assembly.
- `service/internal/catalog_indexer/indexer/metadata/indexer.go`: raw-to-metadata conversion stage.
- `service/internal/search/metadata/metadata_service.go`: HTTP metadata dispatch switches and the range scan helper.
- `service/internal/catalog/allwise/allwise.go`: compact catalog package with nullable source values and sentinel handling.
- `service/internal/catalog/gaia/gaia.go`: simple catalog package.
- `service/internal/catalog/erosita/erosita.go`: large catalog package with catalog RA/DEC fields.
- `service/internal/repository/bulk_insert.go`: generic bulk insert helper over `repository.Insertable`.
- `service/internal/repository/insert.go`: hand-written `Insert` methods on the generated metadata models.
- `service/internal/db/query.sql`: SQLC query definitions.
- `service/internal/db/sqlc.yaml`: SQLC generation and type override configuration.
- `service/internal/testutils/catalog_db.go`: migrated SQLite database helper for round-trip tests.
