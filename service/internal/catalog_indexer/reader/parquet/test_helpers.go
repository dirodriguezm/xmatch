package parquet_reader

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
)

type ReaderBuilder[T any] struct {
	ReaderConfig  *config.ReaderConfig
	t             *testing.T
	Source        *source.Source
	OutputChannel chan indexer.ReaderResult
}

func AReader[T any](t *testing.T) *ReaderBuilder[T] {
	return &ReaderBuilder[T]{
		t: t,
		ReaderConfig: &config.ReaderConfig{
			Type:            "csv",
			FirstLineHeader: true,
			BatchSize:       1,
		},
		OutputChannel: make(chan indexer.ReaderResult),
	}
}

func (builder *ReaderBuilder[T]) WithType(t string) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.ReaderConfig = &config.ReaderConfig{
		Type:      t,
		BatchSize: 1,
	}
	return builder
}

func (builder *ReaderBuilder[T]) WithBatchSize(size int) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.ReaderConfig.BatchSize = size
	return builder
}

func (builder *ReaderBuilder[T]) WithOutputChannel(ch chan indexer.ReaderResult) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.OutputChannel = ch
	return builder
}

func (builder *ReaderBuilder[T]) WithSource(src *source.Source) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.Source = src
	return builder
}

func (builder *ReaderBuilder[T]) Build() indexer.Reader {
	builder.t.Helper()

	r, err := NewParquetReader(
		builder.Source,
		builder.OutputChannel,
		WithParquetBatchSize[T](builder.ReaderConfig.BatchSize),
	)
	require.NoError(builder.t, err)
	return r
}
