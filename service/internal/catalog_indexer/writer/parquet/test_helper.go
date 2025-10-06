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

package parquet_writer

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type TestStruct struct {
	Oid string  `parquet:"name=oid, type=BYTE_ARRAY"`
	Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
}

type ParquetWriterBuilder[T any] struct {
	t *testing.T

	cfg *config.WriterConfig
	ctx context.Context
}

func AWriter[T any](t *testing.T) *ParquetWriterBuilder[T] {
	t.Helper()

	return &ParquetWriterBuilder[T]{
		t:   t,
		cfg: &config.WriterConfig{OutputFile: "test.parquet", Schema: config.TestSchema},
		ctx: context.Background(),
	}
}

func (b *ParquetWriterBuilder[T]) WithOutputFile(file string) *ParquetWriterBuilder[T] {
	b.t.Helper()

	b.cfg.OutputFile = file
	return b
}

func (b *ParquetWriterBuilder[T]) WithContext(ctx context.Context) *ParquetWriterBuilder[T] {
	b.t.Helper()

	b.ctx = ctx
	return b
}

func (b *ParquetWriterBuilder[T]) Build() *ParquetWriter[T] {
	b.t.Helper()

	w, err := New[T](b.cfg, b.ctx)
	if err != nil {
		b.t.Fatal(err)
	}
	return w
}
