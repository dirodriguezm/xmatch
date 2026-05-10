package gaia

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

type GaiaStore interface {
	InsertGaiaWithoutParams(context.Context, repository.Gaia) error
	GetGaia(context.Context, string) (repository.GetGaiaRow, error)
	BulkInsertGaia(context.Context, *sql.DB, []any) error
	BulkGetGaia(context.Context, []string) ([]repository.BulkGetGaiaRow, error)
	GetGaiaFromPixels(context.Context, []int64) ([]repository.GetGaiaFromPixelsRow, error)
}

type Adapter struct {
	store GaiaStore
}

func init() {
	catalog.Register("gaia", func(store any) (catalog.CatalogAdapter, error) {
		s, _ := store.(GaiaStore)
		return &Adapter{store: s}, nil
	})
}

func (a Adapter) Name() string {
	return "gaia"
}

func (a Adapter) NewRawRecord() any {
	return repository.GaiaInputSchema{}
}

func (a Adapter) NewParquetWriter(cfg config.WriterConfig, ctx context.Context) (writer.Writer, error) {
	return parquet_writer.New[repository.Gaia](cfg, ctx)
}

func (a Adapter) NewParquetReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return parquet_reader.NewParquetReader(
		src,
		parquet_reader.WithParquetBatchSize[repository.GaiaInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) NewFitsReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return fits_reader.NewFitsReader(
		src,
		fits_reader.WithBatchSize[repository.GaiaInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) BulkInsertFn() func(context.Context, *sql.DB, []any) error {
	if a.store == nil {
		return func(ctx context.Context, db *sql.DB, rows []any) error {
			return fmt.Errorf("gaia adapter has no store")
		}
	}
	return a.store.BulkInsertGaia
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("gaia adapter has no store")
	}
	return a.store.GetGaia(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("gaia adapter has no store")
	}
	return a.store.BulkGetGaia(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.MetadataWithCoordinates, error) {
	if a.store == nil {
		return nil, fmt.Errorf("gaia adapter has no store")
	}
	rows, err := a.store.GetGaiaFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.MetadataWithCoordinates, len(rows))
	for i, r := range rows {
		result[i] = r
	}
	return result, nil
}

func (a Adapter) ConvertToMetadata(obj repository.MetadataWithCoordinates) repository.Metadata {
	row := obj.(repository.GetGaiaFromPixelsRow)
	return repository.Gaia{
		ID:                  row.ID,
		PhotGMeanFlux:       row.PhotGMeanFlux,
		PhotGMeanFluxError:  row.PhotGMeanFluxError,
		PhotGMeanMag:        row.PhotGMeanMag,
		PhotBpMeanFlux:      row.PhotBpMeanFlux,
		PhotBpMeanFluxError: row.PhotBpMeanFluxError,
		PhotBpMeanMag:       row.PhotBpMeanMag,
		PhotRpMeanFlux:      row.PhotRpMeanFlux,
		PhotRpMeanFluxError: row.PhotRpMeanFluxError,
		PhotRpMeanMag:       row.PhotRpMeanMag,
	}
}

func (a Adapter) ConvertToMastercat(raw any, mapper *healpix.HEALPixMapper) (repository.Mastercat, error) {
	schema := raw.(repository.GaiaInputSchema)
	ipix := mapper.PixelAt(healpix.RADec(schema.RA, schema.Dec))
	return repository.Mastercat{
		ID:   schema.Designation,
		Ipix: ipix,
		Ra:   schema.RA,
		Dec:  schema.Dec,
		Cat:  "gaia",
	}, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (repository.Metadata, error) {
	schema := raw.(repository.GaiaInputSchema)
	return repository.Gaia{
		ID:                  schema.Designation,
		PhotGMeanFlux:       repository.NullFloat64{sql.NullFloat64{Float64: schema.PhotGMeanFlux, Valid: true}},
		PhotGMeanFluxError:  repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotGMeanFluxError), Valid: true}},
		PhotGMeanMag:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotGMeanMag), Valid: true}},
		PhotBpMeanFlux:      repository.NullFloat64{sql.NullFloat64{Float64: schema.PhotBpMeanFlux, Valid: true}},
		PhotBpMeanFluxError: repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotBpMeanFluxError), Valid: true}},
		PhotBpMeanMag:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotBpMeanMag), Valid: true}},
		PhotRpMeanFlux:      repository.NullFloat64{sql.NullFloat64{Float64: schema.PhotRpMeanFlux, Valid: true}},
		PhotRpMeanFluxError: repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotRpMeanFluxError), Valid: true}},
		PhotRpMeanMag:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.PhotRpMeanMag), Valid: true}},
	}, nil
}
