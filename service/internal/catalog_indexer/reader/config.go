package reader

type CsvReaderOption func(r *CsvReader)

func WithHeader(header []string) CsvReaderOption {
	return func(r *CsvReader) {
		r.Header = header
	}
}

func WithBatchSize(size int) CsvReaderOption {
	return func(r *CsvReader) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}
