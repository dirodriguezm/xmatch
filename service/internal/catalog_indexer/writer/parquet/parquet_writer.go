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
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	pwriter "github.com/xitongsys/parquet-go/writer"
)

type ParquetWriter[T any] struct {
	parquetWriter *pwriter.ParquetWriter
	pfile         *os.File
}

func New[T any](cfg *config.WriterConfig, ctx context.Context) (*ParquetWriter[T], error) {
	slog.Debug("Creating new ParquetWriter")

	file, err := os.Create(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("ParquetWriter could not create file %s\n%w", cfg.OutputFile, err)
	}

	var schema any
	switch cfg.Schema {
	case config.AllwiseSchema:
		schema = new(repository.Allwise)
	case config.MastercatSchema:
		schema = new(repository.Mastercat)
	case config.VlassSchema:
		schema = new(repository.VlassObjectSchema)
	case config.TestSchema:
		schema = new(TestStruct)
	default:
		return nil, fmt.Errorf("Schema %v not supported", cfg.Schema)
	}

	parquetWriter, err := pwriter.NewParquetWriterFromWriter(file, schema, 1)
	if err != nil {
		return nil, fmt.Errorf("ParquetWriter could not create writer %w", err)
	}

	w := &ParquetWriter[T]{parquetWriter: parquetWriter, pfile: file}
	return w, nil
}

func (w *ParquetWriter[T]) Write(a *actor.Actor, msg actor.Message) {
	slog.Debug("ParquetWriter received message")
	if msg.Error != nil {
		slog.Error("ParquetWriter received error message")
		panic(msg.Error)
	}

	for i := range msg.Rows {
		obj := msg.Rows[i].(T)

		if reflect.DeepEqual(obj, *new(T)) {
			continue // skip empty objects
		}

		slog.Debug("ParquetWriter writing messages", "len", len(msg.Rows))
		if err := w.parquetWriter.Write(obj); err != nil {
			panic(fmt.Errorf("ParquetWriter could not write object %v\n%w", obj, err))
		}
	}
}

func (w *ParquetWriter[T]) Stop(a *actor.Actor) {
	if err := w.parquetWriter.WriteStop(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not stop. Error: %w", err))
	}
	if err := w.pfile.Close(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not close parquet file %w", err))
	}
}
