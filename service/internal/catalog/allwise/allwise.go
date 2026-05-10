package allwise

import (
	"context"
	"database/sql"
	"fmt"

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
	catalog.Register("allwise", func(store any) (catalog.CatalogAdapter, error) {
		s, _ := store.(AllwiseStore)
		return &Adapter{store: s}, nil
	})
}

func (a Adapter) Name() string {
	return "allwise"
}

func (a Adapter) NewInputSchema() repository.InputSchema {
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

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.MetadataWithCoordinates, error) {
	if a.store == nil {
		return nil, fmt.Errorf("allwise adapter has no store")
	}
	rows, err := a.store.GetAllwiseFromPixels(ctx, pixels)
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
	row := obj.(repository.GetAllwiseFromPixelsRow)
	return repository.Allwise{
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
	}
}
