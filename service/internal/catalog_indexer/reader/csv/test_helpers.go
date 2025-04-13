// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csv_reader

import (
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type ReaderBuilder struct {
	ReaderConfig  *config.ReaderConfig
	t             *testing.T
	Source        *source.Source
	OutputChannel []chan reader.ReaderResult
}

func AReader(t *testing.T) *ReaderBuilder {
	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
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

func (builder *ReaderBuilder) WithOutputChannels(ch []chan reader.ReaderResult) *ReaderBuilder {
	builder.t.Helper()

	builder.OutputChannel = ch
	return builder
}

func (builder *ReaderBuilder) WithSource(src *source.Source) *ReaderBuilder {
	builder.t.Helper()

	builder.Source = src
	return builder
}

func (builder *ReaderBuilder) Build() reader.Reader {
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

func (t *TestSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	cat := "vlass"
	return repository.ParquetMastercat{
		ID:   &t.Oid,
		Ra:   &t.Ra,
		Dec:  &t.Dec,
		Cat:  &cat,
		Ipix: &ipix,
	}
}

func (t *TestSchema) ToMetadata() any {
	return t
}

func (t *TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
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
