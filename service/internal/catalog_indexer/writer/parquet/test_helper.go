package parquet_writer

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type TestStruct struct {
	Oid string  `parquet:"name=oid, type=BYTE_ARRAY"`
	Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
}

type ParquetWriterBuilder[T any] struct {
	t *testing.T

	cfg   *config.WriterConfig
	input chan writer.WriterInput[T]
	done  chan bool
}

func AWriter[T any](t *testing.T) *ParquetWriterBuilder[T] {
	t.Helper()

	return &ParquetWriterBuilder[T]{
		t:     t,
		cfg:   &config.WriterConfig{OutputFile: "test.parquet", Schema: config.TestSchema},
		input: make(chan writer.WriterInput[T]),
		done:  make(chan bool),
	}
}

func (b *ParquetWriterBuilder[T]) WithOutputFile(file string) *ParquetWriterBuilder[T] {
	b.t.Helper()

	b.cfg.OutputFile = file
	return b
}

func (b *ParquetWriterBuilder[T]) WithMessages(messages []writer.WriterInput[T]) *ParquetWriterBuilder[T] {
	b.t.Helper()

	for i := range messages {
		b.input <- messages[i]
	}
	return b
}

func (b *ParquetWriterBuilder[T]) Build() *ParquetWriter[T] {
	b.t.Helper()

	w, err := NewParquetWriter(b.input, b.done, b.cfg)
	if err != nil {
		b.t.Fatal(err)
	}
	return w
}
