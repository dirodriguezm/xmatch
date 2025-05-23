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

package parquet_reader

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type TestInputSchema struct {
	Oid string
	Ra  float64
	Dec float64
}

func (t *TestInputSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	catalog := "test"
	return repository.ParquetMastercat{
		ID:   &t.Oid,
		Ra:   &t.Ra,
		Dec:  &t.Dec,
		Cat:  &catalog,
		Ipix: &ipix,
	}
}

func (t *TestInputSchema) ToMetadata() any {
	return t
}

func (t *TestInputSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t *TestInputSchema) SetField(name string, val interface{}) {
	switch name {
	case "Ra":
		if v, ok := val.(float64); ok {
			t.Ra = v
		}
	case "Dec":
		if v, ok := val.(float64); ok {
			t.Dec = v
		}
	case "Oid":
		t.Oid = val.(string)
	}
}

func (t *TestInputSchema) GetId() string {
	return t.Oid
}

type ReaderBuilder[T any] struct {
	ReaderConfig  *config.ReaderConfig
	t             *testing.T
	Source        *source.Source
	OutputChannel []chan reader.ReaderResult
}

func AReader[T any](t *testing.T) *ReaderBuilder[T] {
	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
	return &ReaderBuilder[T]{
		t: t,
		ReaderConfig: &config.ReaderConfig{
			Type:            "csv",
			FirstLineHeader: true,
			BatchSize:       1,
		},
		OutputChannel: outputs,
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

func (builder *ReaderBuilder[T]) WithOutputChannels(ch []chan reader.ReaderResult) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.OutputChannel = ch
	return builder
}

func (builder *ReaderBuilder[T]) WithSource(src *source.Source) *ReaderBuilder[T] {
	builder.t.Helper()

	builder.Source = src
	return builder
}

func (builder *ReaderBuilder[T]) Build() reader.Reader {
	builder.t.Helper()

	r, err := NewParquetReader(
		builder.Source,
		builder.OutputChannel,
		WithParquetBatchSize[T](builder.ReaderConfig.BatchSize),
	)
	require.NoError(builder.t, err)
	return r
}
