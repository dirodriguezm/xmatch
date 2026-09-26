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
	pwriter "github.com/xitongsys/parquet-go/writer"
)

type ParquetWriter struct {
	parquetWriter *pwriter.ParquetWriter
	pfile         *os.File
	zeroValue     any
}

func New(cfg config.WriterConfig, ctx context.Context, schema any) (*ParquetWriter, error) {
	slog.Debug("Creating new ParquetWriter")

	schemaType := reflect.TypeOf(schema)
	if schemaType == nil {
		return nil, fmt.Errorf("ParquetWriter schema cannot be nil")
	}
	if schemaType.Kind() == reflect.Pointer {
		schemaType = schemaType.Elem()
	}
	writerSchema := reflect.New(schemaType).Interface()
	zeroValue := reflect.Zero(schemaType).Interface()

	file, err := os.Create(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("ParquetWriter could not create file %s\n%w", cfg.OutputFile, err)
	}

	parquetWriter, err := pwriter.NewParquetWriterFromWriter(file, writerSchema, 1)
	if err != nil {
		return nil, fmt.Errorf("ParquetWriter could not create writer %w", err)
	}

	w := &ParquetWriter{parquetWriter: parquetWriter, pfile: file, zeroValue: zeroValue}
	return w, nil
}

func (w *ParquetWriter) Write(a *actor.Actor, msg actor.Message) {
	slog.Debug("ParquetWriter received message")
	if msg.Error != nil {
		slog.Error("ParquetWriter received error message")
		panic(msg.Error)
	}

	for i := range msg.Rows {
		obj := msg.Rows[i]
		if reflect.DeepEqual(obj, w.zeroValue) {
			continue // skip empty objects
		}

		slog.Debug("ParquetWriter writing messages", "len", len(msg.Rows))
		if err := w.parquetWriter.Write(obj); err != nil {
			panic(fmt.Errorf("ParquetWriter could not write object %v\n%w", obj, err))
		}
	}
}

func (w *ParquetWriter) Stop(a *actor.Actor) {
	if err := w.parquetWriter.WriteStop(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not stop. Error: %w", err))
	}
	if err := w.pfile.Close(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not close parquet file %w", err))
	}
}
