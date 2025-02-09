package indexer

import "github.com/dirodriguezm/xmatch/service/internal/repository"

type ReaderResult struct {
	Rows  []repository.InputSchema
	Error error
}

type IndexerResult struct {
	Objects []repository.Mastercat
	Error   error
}

type WriterInput[T any] struct {
	Error error
	Rows  []T
}

type Reader interface {
	Start()
	Read() ([]repository.InputSchema, error)
	ReadBatch() ([]repository.InputSchema, error)
}

type Writer[T any] interface {
	Start()
	Done()
	Stop()
	Receive(WriterInput[T])
}
