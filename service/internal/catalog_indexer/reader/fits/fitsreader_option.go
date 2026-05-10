package fitsreader

import "github.com/dirodriguezm/xmatch/service/internal/repository"

type FitsReaderOption[T repository.InputSchema] func(r *FitsReader[T])

func WithBatchSize[T repository.InputSchema](size int) FitsReaderOption[T] {
	return func(r *FitsReader[T]) {
		if size <= 0 {
			size = 1
		}
		r.batchSize = size
	}
}
