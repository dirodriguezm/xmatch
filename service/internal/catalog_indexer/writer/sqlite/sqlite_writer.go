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

type SqliteWriter[I, O any] struct {
	*writer.BaseWriter[I, O]
	repository conesearch.Repository
	ctx        context.Context
	src        *source.Source
	parser     ParamsParser[I, O]
	bulkWriter ParamsWriter[O]
}

func NewSqliteWriter[I, O any](
	repository conesearch.Repository,
	ch chan writer.WriterInput[I],
	done chan struct{},
	ctx context.Context,
	src *source.Source,
	parser ParamsParser[I, O],
	bulkWriter ParamsWriter[O],
) *SqliteWriter[I, O] {
	slog.Debug("Creating new SqliteWriter")
	w := &SqliteWriter[I, O]{
		BaseWriter: &writer.BaseWriter[I, O]{
			Ctx:          ctx,
			DoneChannel:  done,
			InboxChannel: ch,
		},
		repository: repository,
		ctx:        ctx,
		src:        src,
		parser:     parser,
		bulkWriter: bulkWriter,
	}
	w.Writer = w
	return w
}

func (w *SqliteWriter[I, O]) Receive(msg writer.WriterInput[I]) {
	slog.Debug("Writer received message")
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(msg.Error)
	}

	params := make([]O, len(msg.Rows))
	for i, object := range msg.Rows {
		// convert the received row to insert params needed by the repository
		p := w.parser.Parse(object)
		params[i] = p
	}
	// insert converted rows
	err := w.bulkWriter.BulkWrite(params)
	if err != nil {
		slog.Error("SqliteWriter could not write objects to database")
		panic(err)
	}
}

func (w *SqliteWriter[I, O]) Stop() {
	w.DoneChannel <- struct{}{}
}
