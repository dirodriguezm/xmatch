# Adding a New Catalog to the ALeRCE xmatch Indexer

This guide provides step-by-step instructions for adding a new astronomical catalog to the ALeRCE xmatch indexer. The process involves modifying database schemas, implementing Go interfaces, and updating configuration files.

## Overview

The indexer supports multiple catalogs (AllWISE, Gaia, eRosita) through a modular architecture. Each catalog requires:
1. Database table definition
2. SQLC configuration for code generation
3. Input schema struct with interface implementations
4. Reader factory registration
5. App initialization updates

## Step 1: Create Database Migration

### 1.1 Create Migration File
Create a new migration file in `service/internal/db/migrations/`:

```bash
# Create migration for new catalog (e.g., "newcatalog")
cd service/internal/db/migrations
touch 003_newcatalog.up.sql
touch 003_newcatalog.down.sql
```

### 1.2 Define Table Schema
Edit `003_newcatalog.up.sql`:

```sql
CREATE TABLE newcatalog (
    id TEXT PRIMARY KEY,
    ra REAL NOT NULL,
    dec REAL NOT NULL,
    -- Add your catalog-specific columns here
    column1 REAL,
    column2 REAL,
    -- Include nullable columns with appropriate defaults
    optional_column REAL DEFAULT NULL
);

-- Create indexes for performance
CREATE INDEX idx_newcatalog_ra_dec ON newcatalog(ra, dec);
CREATE INDEX idx_newcatalog_ipix ON newcatalog(ipix);
```

Edit `003_newcatalog.down.sql`:

```sql
DROP TABLE IF EXISTS newcatalog;
```

### 1.3 Apply Migration
```bash
just migrate dev  # or your database name
```

## Step 2: Update SQLC Configuration

### 2.1 Edit `service/internal/db/sqlc.yaml`
Add column overrides for your new catalog in the `overrides` section:

```yaml
overrides:
  go:
    rename:
      newcatalogum: Newcatalog  # SQLite table name + "um" suffix -> Go struct name
  - column: "newcatalog.id"
    go_struct_tag: 'parquet:"name=id, type=BYTE_ARRAY" json:"id"'
  - column: "newcatalog.ra"
    go_struct_tag: 'parquet:"name=ra, type=DOUBLE" json:"ra"'
  - column: "newcatalog.dec"
    go_struct_tag: 'parquet:"name=dec, type=DOUBLE" json:"dec"'
  # Add overrides for all your catalog columns
  - column: "newcatalog.column1"
    go_struct_tag: 'parquet:"name=column1, type=DOUBLE" json:"column1"'
```

### 2.2 Regenerate SQLC Code
```bash
cd service
sqlc generate
```

## Step 3: Implement Repository Interfaces

### 3.1 Create Input Schema File
Create `service/internal/repository/newcatalog.go`:

```go
package repository

import (
    "context"
    "database/sql"
)

type NewcatalogInputSchema struct {
    ID      string  `parquet:"name=id, type=BYTE_ARRAY"`
    Ra      float64 `parquet:"name=ra, type=DOUBLE"`
    Dec     float64 `parquet:"name=dec, type=DOUBLE"`
    Column1 float64 `parquet:"name=column1, type=DOUBLE"`
    Column2 float64 `parquet:"name=column2, type=DOUBLE"`
}

func (schema NewcatalogInputSchema) GetCoordinates() (float64, float64) {
    return schema.Ra, schema.Dec
}

func (schema NewcatalogInputSchema) GetId() string {
    return schema.ID
}

func (schema NewcatalogInputSchema) FillMetadata() Metadata {
    return &Newcatalog{
        ID:      schema.ID,
        Column1: sql.NullFloat64{Float64: schema.Column1, Valid: true},
        Column2: sql.NullFloat64{Float64: schema.Column2, Valid: true},
    }
}

func (schema NewcatalogInputSchema) FillMastercat(ipix int64) Mastercat {
    return Mastercat{
        ID:   schema.ID,
        Ipix: ipix,
        Ra:   schema.Ra,
        Dec:  schema.Dec,
        Cat:  "newcatalog",
    }
}

// Implement required interface methods for the model
func (n Newcatalog) GetId() string {
    return n.ID
}

func (n Newcatalog) GetCatalog() string {
    return "NewCatalog"
}

// Add bulk insert helper if needed
func (q *Queries) InsertNewcatalogWithoutParams(ctx context.Context, arg Newcatalog) error {
    _, err := q.db.ExecContext(ctx, insertNewcatalog,
        arg.ID,
        arg.Column1,
        arg.Column2,
    )
    return err
}
```

### 3.2 Update `service/internal/repository/input_schema.go`
Add a variable for your input schema:

```go
var NewcatalogInputSchema NewcatalogInputSchema
```

## Step 4: Update Reader Factory

### 4.1 Edit `service/internal/catalog_indexer/reader/factory/reader_factory.go`
Update the `parquetFactory` function to support your catalog:

```go
func parquetFactory(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
    switch strings.ToLower(src.CatalogName) {
    case "allwise":
        return parquet_reader.NewParquetReader(
            src,
            parquet_reader.WithParquetBatchSize[repository.AllwiseInputSchema](cfg.BatchSize),
        )
    case "gaia":
        return parquet_reader.NewParquetReader(
            src,
            parquet_reader.WithParquetBatchSize[repository.GaiaInputSchema](cfg.BatchSize),
        )
    case "newcatalog":  // Add your catalog
        return parquet_reader.NewParquetReader(
            src,
            parquet_reader.WithParquetBatchSize[repository.NewcatalogInputSchema](cfg.BatchSize),
        )
    default:
        return nil, fmt.Errorf("Schema not found for catalog %s", src.CatalogName)
    }
}
```

## Step 5: Update App Initialization

### 5.1 Edit `service/internal/app/indexer.go`
Add constants and update switch statements:

```go
const ALLWISE = "allwise"
const GAIA = "gaia"
const EROSITA = "erosita"
const NEWCATALOG = "newcatalog"  // Add your catalog

// Update MetadataWriter function
func MetadataWriter(ctx context.Context, cfg config.Config, repo conesearch.Repository, src *source.Source) (*actor.Actor, error) {
    switch cfg.CatalogIndexer.MetadataWriter.Type {
    case "parquet":
        var w writer.Writer
        var err error
        switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
        case ALLWISE:
            w, err = parquet_writer.New[repository.Allwise](cfg.CatalogIndexer.MetadataWriter, ctx)
        case GAIA:
            w, err = parquet_writer.New[repository.Gaia](cfg.CatalogIndexer.MetadataWriter, ctx)
        case EROSITA:
            w, err = parquet_writer.New[repository.Erosita](cfg.CatalogIndexer.MetadataWriter, ctx)
        case NEWCATALOG:  // Add your catalog
            w, err = parquet_writer.New[repository.Newcatalog](cfg.CatalogIndexer.MetadataWriter, ctx)
        default:
            err = fmt.Errorf("Unknown catalog %s", cfg.CatalogIndexer.Source.CatalogName)
        }
        // ... rest of function
    case "sqlite":
        switch strings.ToLower(cfg.CatalogIndexer.Source.CatalogName) {
        case ALLWISE:
            w := sqlite_writer.New(repo, ctx, repo.BulkInsertAllwise)
            return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
        case GAIA:
            w := sqlite_writer.New(repo, ctx, repo.BulkInsertGaia)
            return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
        case EROSITA:
            w := sqlite_writer.New(repo, ctx, repo.BulkInsertErosita)
            return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
        case NEWCATALOG:  // Add your catalog
            w := sqlite_writer.New(repo, ctx, repo.BulkInsertNewcatalog)
            return actor.New("metadata writer", cfg.CatalogIndexer.ChannelSize, w.Write, w.Stop, nil, ctx), nil
        default:
            return nil, fmt.Errorf("Unknown catalog %s", cfg.CatalogIndexer.Source.CatalogName)
        }
    }
}

// Update MastercatIndexer function
func MastercatIndexer(cfg config.CatalogIndexerConfig, writer *actor.Actor, ctx context.Context) (*actor.Actor, error) {
    fillMastercat := func(schema repository.InputSchema, ipix int64) repository.Mastercat {
        switch cfg.Source.CatalogName {
        case ALLWISE:
            return repository.AllwiseInputSchema.FillMastercat(schema.(repository.AllwiseInputSchema), ipix)
        case GAIA:
            return repository.GaiaInputSchema.FillMastercat(schema.(repository.GaiaInputSchema), ipix)
        case EROSITA:
            return repository.ErositaInputSchema.FillMastercat(schema.(repository.ErositaInputSchema), ipix)
        case NEWCATALOG:  // Add your catalog
            return repository.NewcatalogInputSchema.FillMastercat(schema.(repository.NewcatalogInputSchema), ipix)
        default:
            panic("Catalog not supported")
        }
    }
    // ... rest of function
}

// Update MetadataIndexer function
func MetadataIndexer(cfg config.CatalogIndexerConfig, writer *actor.Actor, ctx context.Context) *actor.Actor {
    fillMetadata := func(schema repository.InputSchema) repository.Metadata {
        switch cfg.Source.CatalogName {
        case ALLWISE:
            return repository.AllwiseInputSchema.FillMetadata(schema.(repository.AllwiseInputSchema))
        case GAIA:
            return repository.GaiaInputSchema.FillMetadata(schema.(repository.GaiaInputSchema))
        case EROSITA:
            return repository.ErositaInputSchema.FillMetadata(schema.(repository.ErositaInputSchema))
        case NEWCATALOG:  // Add your catalog
            return repository.NewcatalogInputSchema.FillMetadata(schema.(repository.NewcatalogInputSchema))
        default:
            panic("Catalog not supported")
        }
    }
    // ... rest of function
}
```

## Step 6: Create Configuration

### 6.1 Create Catalog Configuration
Create a configuration file in `service/configs/` (e.g., `newcatalog.yaml`):

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
  indexer_writer:
    type: "sqlite"
  metadata_writer:
    type: "sqlite"
  channel_size: 50000
```

### 6.2 Update Main Config (Optional)
Add your catalog to the main `service/config.yaml` if needed:

```yaml
catalog_indexer:
  source:
    catalog_name: "vlass|ztf|allwise|gaia|erosita|newcatalog"
```

## Step 7: Add Bulk Insert Support (Optional)

If you need bulk insert functionality, update `service/internal/repository/bulk_insert.go`:

```go
func (q *Queries) BulkInsertNewcatalog(ctx context.Context, arg []Newcatalog) error {
    // Implement bulk insert logic similar to other catalogs
    // See BulkInsertAllwise for reference
}
```

## Common Patterns and Examples

### Handling Nullable Fields
```go
type NewcatalogInputSchema struct {
    ID         string   `parquet:"name=id, type=BYTE_ARRAY"`
    Ra         float64  `parquet:"name=ra, type=DOUBLE"`
    Dec        float64  `parquet:"name=dec, type=DOUBLE"`
    Optional   *float64 `parquet:"name=optional, type=DOUBLE"`  // Pointer for nullable
}

func (schema NewcatalogInputSchema) FillMetadata() Metadata {
    newcatalog := &Newcatalog{
        ID: schema.ID,
    }
    
    if schema.Optional != nil {
        newcatalog.Optional = sql.NullFloat64{Float64: *schema.Optional, Valid: true}
    } else {
        newcatalog.Optional = sql.NullFloat64{Float64: -9999.0, Valid: false}
    }
    
    return newcatalog
}
```

### CSV Reader Support
If your catalog uses CSV format, ensure column mappings are correct in your configuration.

## References

- Existing implementations: `service/internal/repository/allwise.go`, `gaia.go`, `erosita.go`
- Reader factory: `service/internal/catalog_indexer/reader/factory/reader_factory.go`
- App initialization: `service/internal/app/indexer.go`
- SQLC configuration: `service/internal/db/sqlc.yaml`

