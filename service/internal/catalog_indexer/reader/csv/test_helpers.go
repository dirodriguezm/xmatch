package csv_reader

import (
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type ReaderBuilder struct {
	ReaderConfig  *config.ReaderConfig
	t             *testing.T
	Source        *source.Source
	OutputChannel []chan indexer.ReaderResult
}

func AReader(t *testing.T) *ReaderBuilder {
	outputs := make([]chan indexer.ReaderResult, 1)
	outputs[0] = make(chan indexer.ReaderResult)
	return &ReaderBuilder{
		t: t,
		ReaderConfig: &config.ReaderConfig{
			Type:            "csv",
			FirstLineHeader: true,
			BatchSize:       1,
		},
		OutputChannel: outputs,
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

func (builder *ReaderBuilder) WithOutputChannels(ch []chan indexer.ReaderResult) *ReaderBuilder {
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

type TestSchema struct {
	Ra  float64
	Dec float64
	Oid string
}

func (t *TestSchema) ToMastercat() repository.ParquetMastercat {
	cat := "vlass"
	return repository.ParquetMastercat{
		ID:  &t.Oid,
		Ra:  &t.Ra,
		Dec: &t.Dec,
		Cat: &cat,
	}
}

func (t *TestSchema) SetField(name string, val interface{}) {
	switch name {
	case "ra":
		if v, ok := val.(float64); ok {
			t.Ra = v
		}
		if v, ok := val.(string); ok {
			ra, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}
			t.Ra = ra
		}
	case "dec":
		if v, ok := val.(float64); ok {
			t.Dec = v
		}
		if v, ok := val.(string); ok {
			dec, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}
			t.Dec = dec
		}
	case "oid":
		t.Oid = val.(string)
	}
}
