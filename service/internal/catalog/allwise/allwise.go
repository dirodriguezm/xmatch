package allwise

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	fits_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/fits"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type AllwiseStore interface {
	InsertAllwiseWithoutParams(context.Context, repository.Allwise) error
	GetAllwise(context.Context, string) (repository.GetAllwiseRow, error)
	BulkInsertAllwise(context.Context, *sql.DB, []any) error
	BulkGetAllwise(context.Context, []string) ([]repository.BulkGetAllwiseRow, error)
	GetAllwiseFromPixels(context.Context, []int64) ([]repository.GetAllwiseFromPixelsRow, error)
}

type Adapter struct {
	store AllwiseStore
}

func init() {
	catalog.Register("allwise", func(store any) (catalog.CatalogIndexAdapter, error) {
		s, _ := store.(AllwiseStore)
		return &Adapter{store: s}, nil
	})
}

func (a Adapter) Name() string {
	return "allwise"
}

func (a Adapter) NewRawRecord() any {
	return repository.AllwiseInputSchema{}
}

func (a Adapter) NewParquetWriter(cfg config.WriterConfig, ctx context.Context) (writer.Writer, error) {
	return parquet_writer.New[repository.Allwise](cfg, ctx)
}

func (a Adapter) NewParquetReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return parquet_reader.NewParquetReader(
		src,
		parquet_reader.WithParquetBatchSize[repository.AllwiseInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) NewFitsReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return fits_reader.NewFitsReader(
		src,
		a,
		fits_reader.WithBatchSize[repository.AllwiseInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) BulkInsertFn() func(context.Context, *sql.DB, []any) error {
	if a.store == nil {
		return func(ctx context.Context, db *sql.DB, rows []any) error {
			return fmt.Errorf("allwise adapter has no store")
		}
	}
	return a.store.BulkInsertAllwise
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("allwise adapter has no store")
	}
	return a.store.GetAllwise(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("allwise adapter has no store")
	}
	return a.store.BulkGetAllwise(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error) {
	if a.store == nil {
		return nil, fmt.Errorf("allwise adapter has no store")
	}
	rows, err := a.store.GetAllwiseFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.Metadata, len(rows))
	for i, r := range rows {
		result[i] = convertAllwiseFromPixelsRowToMetadata(r)
	}
	return result, nil
}

func (a Adapter) ConvertToMetadata(obj any) repository.Metadata {
	row := obj.(repository.GetAllwiseFromPixelsRow)
	return convertAllwiseFromPixelsRowToMetadata(row)
}

func (a Adapter) ConvertToMastercat(raw any, mapper *healpix.HEALPixMapper) (repository.Mastercat, error) {
	schema := raw.(repository.AllwiseInputSchema)
	ra := 0.0
	dec := 0.0
	if schema.Ra != nil {
		ra = *schema.Ra
	}
	if schema.Dec != nil {
		dec = *schema.Dec
	}
	ipix := mapper.PixelAt(healpix.RADec(ra, dec))
	mc := repository.Mastercat{
		Ipix: ipix,
		Cat:  "allwise",
	}
	if schema.Source_id != nil {
		mc.ID = *schema.Source_id
	}
	if schema.Ra != nil {
		mc.Ra = *schema.Ra
	}
	if schema.Dec != nil {
		mc.Dec = *schema.Dec
	}
	return mc, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (any, error) {
	schema := raw.(repository.AllwiseInputSchema)
	return convertAllwiseInputToMetadata(schema), nil
}

func convertAllwiseFromPixelsRowToMetadata(row repository.GetAllwiseFromPixelsRow) repository.Metadata {
	return repository.Metadata{
		ID:      row.ID,
		Catalog: row.GetCatalog(),
		Ra:      row.Ra,
		Dec:     row.Dec,
		Object: repository.Allwise{
			ID:         row.ID,
			Cntr:       row.Cntr,
			W1mpro:     row.W1mpro,
			W1sigmpro:  row.W1sigmpro,
			W2mpro:     row.W2mpro,
			W2sigmpro:  row.W2sigmpro,
			W3mpro:     row.W3mpro,
			W3sigmpro:  row.W3sigmpro,
			W4mpro:     row.W4mpro,
			W4sigmpro:  row.W4sigmpro,
			JM2mass:    row.JM2mass,
			JMsig2mass: row.JMsig2mass,
			HM2mass:    row.HM2mass,
			HMsig2mass: row.HMsig2mass,
			KM2mass:    row.KM2mass,
			KMsig2mass: row.KMsig2mass,
		},
	}
}

func convertAllwiseInputToMetadata(schema repository.AllwiseInputSchema) repository.Allwise {
	allwise := repository.Allwise{}
	if schema.Source_id != nil {
		allwise.ID = *schema.Source_id
	}
	if schema.Cntr != nil {
		allwise.Cntr = *schema.Cntr
	}
	if schema.W1mpro != nil {
		allwise.W1mpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W1mpro, Valid: true}}
	} else {
		allwise.W1mpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W1sigmpro != nil {
		allwise.W1sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W1sigmpro, Valid: true}}
	} else {
		allwise.W1sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2mpro != nil {
		allwise.W2mpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W2mpro, Valid: true}}
	} else {
		allwise.W2mpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2sigmpro != nil {
		allwise.W2sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W2sigmpro, Valid: true}}
	} else {
		allwise.W2sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3mpro != nil {
		allwise.W3mpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W3mpro, Valid: true}}
	} else {
		allwise.W3mpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3sigmpro != nil {
		allwise.W3sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W3sigmpro, Valid: true}}
	} else {
		allwise.W3sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4mpro != nil {
		allwise.W4mpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W4mpro, Valid: true}}
	} else {
		allwise.W4mpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4sigmpro != nil {
		allwise.W4sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: *schema.W4sigmpro, Valid: true}}
	} else {
		allwise.W4sigmpro = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_m_2mass != nil {
		allwise.JM2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.J_m_2mass, Valid: true}}
	} else {
		allwise.JM2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_msig_2mass != nil {
		allwise.JMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.J_msig_2mass, Valid: true}}
	} else {
		allwise.JMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_m_2mass != nil {
		allwise.HM2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.H_m_2mass, Valid: true}}
	} else {
		allwise.HM2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_msig_2mass != nil {
		allwise.HMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.H_msig_2mass, Valid: true}}
	} else {
		allwise.HMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_m_2mass != nil {
		allwise.KM2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.K_m_2mass, Valid: true}}
	} else {
		allwise.KM2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_msig_2mass != nil {
		allwise.KMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: *schema.K_msig_2mass, Valid: true}}
	} else {
		allwise.KMsig2mass = repository.NullFloat64{sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	return allwise
}
