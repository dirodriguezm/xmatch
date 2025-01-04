package csv_reader

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
)

type ReaderBuilder struct {
	ReaderConfig  *config.ReaderConfig
	t             *testing.T
	Source        *source.Source
	OutputChannel chan indexer.ReaderResult
}

func AReader(t *testing.T) *ReaderBuilder {
	return &ReaderBuilder{
		t: t,
		ReaderConfig: &config.ReaderConfig{
			Type:            "csv",
			FirstLineHeader: true,
			BatchSize:       1,
		},
		OutputChannel: make(chan indexer.ReaderResult),
	}
}

func (builder *ReaderBuilder) WithType(t string) *ReaderBuilder {
	builder.t.Helper()

	builder.ReaderConfig = &config.ReaderConfig{
		Type:      t,
		BatchSize: 1,
	}
	return builder
}

func (builder *ReaderBuilder) WithBatchSize(size int) *ReaderBuilder {
	builder.t.Helper()

	builder.ReaderConfig.BatchSize = size
	return builder
}

func (builder *ReaderBuilder) WithOutputChannel(ch chan indexer.ReaderResult) *ReaderBuilder {
	builder.t.Helper()

	builder.OutputChannel = ch
	return builder
}

func (builder *ReaderBuilder) WithSource(src *source.Source) *ReaderBuilder {
	builder.t.Helper()

	builder.Source = src
	return builder
}

func (builder *ReaderBuilder) Build() indexer.Reader {
	builder.t.Helper()

	r, err := NewCsvReader(
		builder.Source,
		builder.OutputChannel,
		WithCsvBatchSize(builder.ReaderConfig.BatchSize),
		WithHeader(builder.ReaderConfig.Header),
		WithFirstLineHeader(builder.ReaderConfig.FirstLineHeader),
	)
	require.NoError(builder.t, err)
	return r
}
