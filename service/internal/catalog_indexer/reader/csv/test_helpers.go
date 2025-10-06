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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type ReaderBuilder struct {
	ReaderConfig *config.ReaderConfig
	t            *testing.T
	Source       *source.Source
}

func AReader(t *testing.T) *ReaderBuilder {
	return &ReaderBuilder{
		t: t,
		ReaderConfig: &config.ReaderConfig{
			Type:            "csv",
			FirstLineHeader: true,
			BatchSize:       1,
		},
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

func (builder *ReaderBuilder) WithSource(src *source.Source) *ReaderBuilder {
	builder.t.Helper()

	builder.Source = src
	return builder
}

func (builder *ReaderBuilder) Build() reader.Reader {
	builder.t.Helper()

	r, err := NewCsvReader(
		builder.Source,
		WithHeader(builder.ReaderConfig.Header),
		WithFirstLineHeader(builder.ReaderConfig.FirstLineHeader),
	)
	require.NoError(builder.t, err)
	return r
}

type TestSchema struct {
	Oid string
	Ra  float64
	Dec float64
}

func (t *TestSchema) FillMastercat(dst *repository.Mastercat, ipix int64) {
	dst.ID = t.Oid
	dst.Ra = t.Ra
	dst.Dec = t.Dec
	dst.Cat = "test"
	dst.Ipix = ipix
}

func (t *TestSchema) FillMetadata(dst repository.Metadata) {}

func (t *TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t *TestSchema) GetId() string {
	return t.Oid
}
