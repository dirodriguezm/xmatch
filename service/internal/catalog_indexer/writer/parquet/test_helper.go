package parquet_writer

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type ParquetWriterBuilder[T any] struct {
	t *testing.T

	cfg   *config.WriterConfig
	input chan indexer.WriterInput
	done  chan bool
}

func AWriter[T any](t *testing.T) *ParquetWriterBuilder[T] {
	t.Helper()

	return &ParquetWriterBuilder[T]{
		t:     t,
		cfg:   &config.WriterConfig{OutputFile: "test.parquet"},
		input: make(chan indexer.WriterInput),
		done:  make(chan bool),
	}
}

func (b *ParquetWriterBuilder[T]) WithOutputFile(file string) *ParquetWriterBuilder[T] {
	b.t.Helper()

	b.cfg.OutputFile = file
	return b
}

func (b *ParquetWriterBuilder[T]) WithMessages(messages []indexer.WriterInput) *ParquetWriterBuilder[T] {
	b.t.Helper()

	for i := range messages {
		b.input <- messages[i]
	}
	return b
}

func (b *ParquetWriterBuilder[T]) Build() *ParquetWriter[T] {
	b.t.Helper()

	w, err := NewParquetWriter[T](b.input, b.done, b.cfg)
	if err != nil {
		b.t.Fatal(err)
	}
	return w
}
