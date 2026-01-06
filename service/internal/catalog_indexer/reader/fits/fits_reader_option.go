package fits_reader

type FitsReaderOption func(r FitsReader) FitsReader

func WithBatchSize(size int) FitsReaderOption {
	return func(r FitsReader) FitsReader {
		if size <= 0 {
			size = 1
		}
		r.batchSize = size
		return r
	}
}
