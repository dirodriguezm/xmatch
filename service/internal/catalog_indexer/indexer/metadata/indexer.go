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

package metadata

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
)

type Indexer struct {
	adapter catalog.CatalogAdapter
}

func New(adapter catalog.CatalogAdapter) *Indexer {
	return &Indexer{adapter: adapter}
}

func (ind *Indexer) Index(a *actor.Actor, msg actor.Message) {
	slog.Debug("Metadata Indexer Received Message")
	if msg.Error != nil {
		slog.Error("Metadata Indexer Received Error", "error", msg.Error)
		a.Broadcast(actor.Message{Error: msg.Error, Rows: nil})
		return
	}

	outputBatch, err := ind.getOutputBatch(msg.Rows)
	if err != nil {
		a.Broadcast(actor.Message{Error: err, Rows: nil})
		return
	}

	slog.Debug("Metadata Indexer Sending Message", "len", len(outputBatch))
	a.Broadcast(actor.Message{Rows: outputBatch, Error: nil})
}

func (ind *Indexer) getOutputBatch(rows []any) ([]any, error) {
	outputBatch := make([]any, len(rows))
	for i := range rows {
		md, err := ind.adapter.ConvertToMetadataFromRaw(rows[i])
		if err != nil {
			return nil, err
		}
		outputBatch[i] = md
	}
	return outputBatch, nil
}
