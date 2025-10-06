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
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Indexer struct {
	catalog string
}

func New(catalog string) *Indexer {
	return &Indexer{catalog: strings.ToLower(catalog)}
}

func (ind *Indexer) Index(a *actor.Actor, msg actor.Message) {
	slog.Debug("Metadata Indexer Received Message")
	if msg.Error != nil {
		slog.Error("Metadata Indexer Received Error", "error", msg.Error)
		a.Broadcast(actor.Message{Error: msg.Error, Rows: nil})
		return
	}

	outputBatch := getOutputBatch(ind.catalog, msg.Rows)

	slog.Debug("Metadata Indexer Sending Message", "len", len(outputBatch))
	a.Broadcast(actor.Message{Rows: outputBatch, Error: nil})
}

func getOutputBatch(catalog string, rows []any) []any {
	outputBatch := make([]any, len(rows))

	switch catalog {
	case "gaia":
		for i := range rows {
			var obj repository.Gaia
			rows[i].(repository.InputSchema).FillMetadata(&obj)
			outputBatch[i] = obj
		}
	case "allwise":
		for i := range rows {
			var obj repository.Allwise
			rows[i].(repository.InputSchema).FillMetadata(&obj)
			outputBatch[i] = obj
		}
	}
	return outputBatch
}
