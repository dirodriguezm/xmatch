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

package sqlite_writer

import (
	"context"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter[T any] struct {
	*writer.BaseWriter[T]
	repository conesearch.Repository
	ctx        context.Context
	src        *source.Source
	bulkWriter ParamsWriter[T]
}

func NewSqliteWriter[T any](
	repository conesearch.Repository,
	ch chan writer.WriterInput[T],
	done chan struct{},
	ctx context.Context,
	src *source.Source,
	bulkWriter ParamsWriter[T],
) *SqliteWriter[T] {
	slog.Debug("Creating new SqliteWriter")
	w := &SqliteWriter[T]{
		BaseWriter: &writer.BaseWriter[T]{
			Ctx:          ctx,
			DoneChannel:  done,
			InboxChannel: ch,
		},
		repository: repository,
		ctx:        ctx,
		src:        src,
		bulkWriter: bulkWriter,
	}
	w.Writer = w
	return w
}

func (w *SqliteWriter[T]) Receive(msg writer.WriterInput[T]) {
	slog.Debug("Writer received message")
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(msg.Error)
	}

	// insert rows
	err := w.bulkWriter.BulkWrite(msg.Rows)
	if err != nil {
		slog.Error("SqliteWriter could not write objects to database")
		panic(err)
	}
}

func (w *SqliteWriter[T]) Stop() {
	w.DoneChannel <- struct{}{}
}
