package parquet_reader

type ParquetReaderOption[T any] func(r *ParquetReader[T])

func WithParquetBatchSize[T any](size int) ParquetReaderOption[T] {
	return func(r *ParquetReader[T]) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}
