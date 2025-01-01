package parquet_reader

type ParquetReaderOption func(r *ParquetReader)

func WithParquetBatchSize(size int) ParquetReaderOption {
	return func(r *ParquetReader) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}

func WithParquetMetadata(md []string) ParquetReaderOption {
	return func(r *ParquetReader) {
		if md != nil && len(md) > 0 {
			r.metadata = md
		}
	}
}
