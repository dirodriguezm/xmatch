# Adding a New Catalog to the ALeRCE xmatch Indexer

This guide describes the current path for adding a new astronomical catalog to the xmatch service. Catalog-specific behavior lives behind `catalog.CatalogAdapter`; the pipeline builds readers, writers, HEALPix mappers, and SQLite bulk writers around that adapter.

Use `newcatalog` below as the placeholder catalog name.

## Overview

A catalog needs these pieces:

1. A metadata table migration under `service/internal/db/migrations/`.
2. SQL queries in `service/internal/db/query.sql`.
3. SQLC overrides in `service/internal/db/sqlc.yaml` when generated types or parquet tags need customization.
4. Regenerated SQLC repository code.
5. A catalog adapter under `service/internal/catalog/newcatalog/`.
6. Adapter registration through a blank import in `service/cmd/main.go`.
7. A config file passed with `CONFIG_PATH` when running the indexer.

The adapter is used by both the indexer and the HTTP API. There is no per-catalog reader factory, writer factory, bulk-insert method, or HTTP resolver registration to edit.

## Step 1: Add Database Migrations

Create the next numbered migration files in `service/internal/db/migrations/`. The current last catalog migration is `005_erosita`, so the next one should be `006_newcatalog` unless new migrations have been added since this guide was written.

```bash
touch service/internal/db/migrations/006_newcatalog.up.sql
touch service/internal/db/migrations/006_newcatalog.down.sql
```

Metadata tables store catalog-specific fields. Coordinates and HEALPix pixels are stored in `mastercat`, so do not add `ipix` to the catalog table unless the catalog itself has an unrelated pixel column.

Example `006_newcatalog.up.sql`:

```sql
CREATE TABLE newcatalog (
    id text NOT NULL,
    mag_g double precision,
    mag_r double precision,
    flag smallint,
    PRIMARY KEY (id)
);

CREATE INDEX idx_newcatalog_mag_g ON newcatalog(mag_g);
```

Example `006_newcatalog.down.sql`:

```sql
DROP TABLE IF EXISTS newcatalog;
```

Apply migrations to a local database when you need to test against SQLite. If you use the example `file:dev.db` config with `just run`, the process runs from `service/`, so migrate `service/dev.db`:

```bash
just migrate service/dev
```

## Step 2: Add SQLC Queries

Edit `service/internal/db/query.sql` and add insert, lookup, bulk lookup, cleanup, and pixel lookup queries.

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

-- name: RemoveAllNewcatalog :exec
DELETE FROM newcatalog;

-- name: GetNewcatalogFromPixels :many
SELECT newcatalog.*, mastercat.ra, mastercat.dec
FROM newcatalog
JOIN mastercat ON mastercat.id = newcatalog.id
WHERE mastercat.ipix IN (sqlc.slice(ipix));
```

The joins are important because metadata responses include RA and DEC from `mastercat`.

## Step 3: Update SQLC Configuration

Edit `service/internal/db/sqlc.yaml`.

Add column overrides under `sql[0].gen.go.overrides` if the generated struct needs parquet tags, JSON tags, or custom nullable wrapper types:

```yaml
          - column: "newcatalog.id"
            go_struct_tag: 'json:"id" parquet:"name=id, type=BYTE_ARRAY"'
          - column: "newcatalog.mag_g"
            go_struct_tag: 'json:"mag_g" parquet:"name=mag_g, type=DOUBLE"'
            go_type:
              type: "NullFloat64"
          - column: "newcatalog.mag_r"
            go_struct_tag: 'json:"mag_r" parquet:"name=mag_r, type=DOUBLE"'
            go_type:
              type: "NullFloat64"
          - column: "newcatalog.flag"
            go_struct_tag: 'json:"flag" parquet:"name=flag, type=INT32"'
            go_type:
              type: "NullInt64"
```

The existing `db_type` overrides near the end of `sqlc.yaml` already map nullable SQLite numeric/text columns to repository null wrappers. Use column overrides when you need specific parquet or JSON tags, or when SQLC's inferred type is not what the adapter should expose.

If SQLC generates an awkward Go type name, add a rename under the root `overrides.go.rename` block at the end of the file:

```yaml
overrides:
  go:
    rename:
      gaium: Gaia
      erositum: Erosita
      newcatalogum: Newcatalog
```

Only add the rename if SQLC actually needs it.

Regenerate repository code from the directory containing `sqlc.yaml`:

```bash
cd service/internal/db
sqlc generate
```

This updates generated files such as `service/internal/repository/models.go` and `service/internal/repository/query.sql.go`.

## Step 4: Add the Catalog Adapter

Create `service/internal/catalog/newcatalog/newcatalog.go`. The adapter connects the generic reader/indexer/writer pipeline to generated repository methods.

```go
package newcatalog

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

const displayName = "NewCatalog"

type InputSchema struct {
	ID   string  `json:"id" parquet:"name=id, type=BYTE_ARRAY"`
	RA   float64 `json:"ra" parquet:"name=ra, type=DOUBLE"`
	Dec  float64 `json:"dec" parquet:"name=dec, type=DOUBLE"`
	MagG float64 `json:"mag_g" parquet:"name=mag_g, type=DOUBLE"`
	MagR float64 `json:"mag_r" parquet:"name=mag_r, type=DOUBLE"`
	Flag int64   `json:"flag" parquet:"name=flag, type=INT64"`
}

type Adapter struct {
	repo *repository.Queries
}

func init() {
	catalog.Register("newcatalog", func(repo *repository.Queries) (catalog.CatalogAdapter, error) {
		return &Adapter{repo: repo}, nil
	})
}

func (a Adapter) Name() string {
	return "newcatalog"
}

func (a Adapter) NewRawRecord() any {
	return InputSchema{}
}

func (a Adapter) NewMetadataRecord() any {
	return repository.Newcatalog{}
}

func (a Adapter) BulkInsertMetadata(ctx context.Context, rows []any) error {
	if a.repo == nil {
		return fmt.Errorf("newcatalog adapter has no repository")
	}
	return repository.BulkInsert(ctx, a.repo, rows, func(ctx context.Context, qtx *repository.Queries, row repository.Newcatalog) error {
		return qtx.InsertNewcatalog(ctx, repository.InsertNewcatalogParams(row))
	})
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("newcatalog adapter has no repository")
	}
	return a.repo.GetNewcatalog(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("newcatalog adapter has no repository")
	}
	return a.repo.BulkGetNewcatalog(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("newcatalog adapter has no repository")
	}
	rows, err := a.repo.GetNewcatalogFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.Metadata, len(rows))
	for i, r := range rows {
		result[i] = convertNewcatalogFromPixelsRowToMetadata(r)
	}
	return result, nil
}

func convertNewcatalogFromPixelsRowToMetadata(row repository.GetNewcatalogFromPixelsRow) repository.Metadata {
	return repository.Metadata{
		ID:      row.ID,
		Catalog: displayName,
		Ra:      row.Ra,
		Dec:     row.Dec,
		Object: repository.Newcatalog{
			ID:   row.ID,
			MagG: row.MagG,
			MagR: row.MagR,
			Flag: row.Flag,
		},
	}
}

func (a Adapter) GetCoordinates(raw any) (float64, float64, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return 0, 0, fmt.Errorf("expected newcatalog.InputSchema, got %T", raw)
	}
	return schema.RA, schema.Dec, nil
}

func (a Adapter) ConvertToMastercat(raw any, ipix int64) (repository.Mastercat, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return repository.Mastercat{}, fmt.Errorf("expected newcatalog.InputSchema, got %T", raw)
	}
	return repository.Mastercat{
		ID:   schema.ID,
		Ipix: ipix,
		Ra:   schema.RA,
		Dec:  schema.Dec,
		Cat:  "newcatalog",
	}, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (any, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return nil, fmt.Errorf("expected newcatalog.InputSchema, got %T", raw)
	}
	return repository.Newcatalog{
		ID:   schema.ID,
		MagG: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: schema.MagG, Valid: true}},
		MagR: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: schema.MagR, Valid: true}},
		Flag: repository.NullInt64{NullInt64: sql.NullInt64{Int64: schema.Flag, Valid: true}},
	}, nil
}
```

Run `gofmt` after editing Go files. Adjust generated field names in the example to match SQLC output. If the catalog table also has `ra` or `dec` columns, SQLC may rename the joined `mastercat.ra`/`mastercat.dec` fields in row structs; copy the pattern from `service/internal/catalog/erosita/erosita.go` in that case.

Notes:

1. `InputSchema` belongs in the catalog adapter package, not `repository`; generated metadata types belong in `repository`.
2. `NewRawRecord` tells CSV, Parquet, and FITS readers which input struct to decode into.
3. `NewMetadataRecord` tells the Parquet metadata writer which generated repository struct to write.
4. `GetCoordinates` returns source coordinates; the mastercat indexer computes `ipix` and passes it into `ConvertToMastercat`.
5. `BulkInsertMetadata` uses the generic `repository.BulkInsert`; do not add a catalog-specific method to `repository/bulk_insert.go` unless the generic path is insufficient.
6. The CSV reader maps records by struct field order and currently supports scalar string, integer, float, and bool fields. For CSV, keep `InputSchema` in source column order or preprocess the file.
7. For Parquet/FITS sources that need to distinguish null from zero, use pointer fields in `InputSchema` and handle nil values in `GetCoordinates`, `ConvertToMastercat`, and `ConvertToMetadataFromRaw`.

## Step 5: Register the Adapter

Edit `service/cmd/main.go` and add a blank import so the adapter `init()` function runs:

```go
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/newcatalog"
```

This registration is shared by the indexer and the HTTP server because both commands are built from `service/cmd/*.go`. The HTTP metadata API uses `catalog.NewResolver(queries)` and does not need per-catalog calls in `start_http_server.go`.

Also update API/setup tests if they assert the exact set of registered catalogs.

## Step 6: Configure and Run the Indexer

Create a local config file and run with `CONFIG_PATH`. There is no required `service/configs/` directory; any path is fine. When using `just run`, `CONFIG_PATH` is resolved from the `service/` working directory unless you pass an absolute path.

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
  indexer_writer:
    type: "sqlite"
  metadata_writer:
    type: "sqlite"
  channel_size: 50000
```

Keep `source.nside` and `indexer.nside` aligned. `source.nside` is stored in the `catalogs` table for search, while `indexer.nside` is used to compute `mastercat.ipix` during indexing.

Run it:

```bash
CONFIG_PATH=newcatalog.yaml just run indexer
```

For CSV sources, set `source.type` and `reader.type` to `csv`. The current CSV reader maps records into `InputSchema` by struct field order, not by header name. If a CSV file has no header row, set `reader.header` in config to any field list so the first data row is not consumed as a header. If every file in a `files:` source has its own header row, set `reader.first_line_header: true` so the header is skipped for each file.

## Step 7: Verify

Run these checks from the repository root unless the command changes directories itself:

```bash
cd service/internal/db
sqlc generate
```

```bash
cd service
gofmt -w internal/catalog/newcatalog cmd
go test ./... -race
```

If you changed `catalog.CatalogAdapter` itself, regenerate mocks:

```bash
just mock
```

## Reference Files

- `service/internal/catalog/adapter.go`: catalog adapter interface and resolver.
- `service/internal/catalog/allwise/allwise.go`: compact adapter example with nullable source values.
- `service/internal/catalog/gaia/gaia.go`: simple adapter example.
- `service/internal/catalog/erosita/erosita.go`: large metadata adapter example with catalog RA/DEC fields.
- `service/internal/catalog_indexer/pipeline/pipeline.go`: generic pipeline using adapters.
- `service/internal/catalog_indexer/indexer/mastercat/indexer.go`: HEALPix pixel computation and mastercat conversion.
- `service/internal/catalog_indexer/indexer/metadata/indexer.go`: raw-to-metadata conversion.
- `service/internal/repository/bulk_insert.go`: generic SQLite bulk insert helper.
- `service/internal/db/query.sql`: SQLC query definitions.
- `service/internal/db/sqlc.yaml`: SQLC generation and type/tag overrides.
