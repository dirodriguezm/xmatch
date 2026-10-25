package fitsreader

type FitsReaderOption[T any] func(r *FitsReader[T])

func WithBatchSize[T any](size int) FitsReaderOption[T] {
	return func(r *FitsReader[T]) {
		if size <= 0 {
			size = 1
		}
		r.batchSize = size
	}
}
