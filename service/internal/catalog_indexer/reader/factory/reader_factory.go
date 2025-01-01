package reader_factory

import (
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

func ReaderFactory(
	src *source.Source,
	outbox chan indexer.ReaderResult,
	cfg *config.ReaderConfig,
) (indexer.Reader, error) {
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
		return parquet_reader.NewParquetReader(
			src,
			outbox,
			parquet_reader.WithParquetBatchSize(cfg.BatchSize),
			parquet_reader.WithParquetMetadata(cfg.Metadata),
		)
	default:
		return nil, fmt.Errorf("Reader type not allowed")
	}
}
