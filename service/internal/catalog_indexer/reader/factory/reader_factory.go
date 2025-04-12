package reader_factory

import (
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func ReaderFactory(
	src *source.Source,
	outbox []chan reader.ReaderResult,
	cfg *config.ReaderConfig,
) (reader.Reader, error) {
	readerType := strings.ToLower(cfg.Type)
	switch readerType {
	case "csv":
		return csv_reader.NewCsvReader(
			src,
			outbox,
			csv_reader.WithCsvBatchSize(cfg.BatchSize),
			csv_reader.WithHeader(cfg.Header),
			csv_reader.WithFirstLineHeader(cfg.FirstLineHeader),
		)
	case "parquet":
		return parquetFactory(src, outbox, cfg)
	default:
		return nil, fmt.Errorf("Reader type not allowed")
	}
}

func parquetFactory(src *source.Source, outbox []chan reader.ReaderResult, cfg *config.ReaderConfig) (reader.Reader, error) {
	catalog := strings.ToLower(src.CatalogName)
	if catalog == "allwise" {
		return parquet_reader.NewParquetReader(
			src,
			outbox,
			parquet_reader.WithParquetBatchSize[repository.AllwiseInputSchema](cfg.BatchSize),
		)
	}
	return nil, fmt.Errorf("Schema not found for catalog")
}
