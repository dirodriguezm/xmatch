package fitsreader

type FitsReaderOption func(r *FitsReader)

func WithBatchSize(size int) FitsReaderOption {
	return func(r *FitsReader) {
		if size <= 0 {
			size = 1
		}
		r.batchSize = size
	}
}
