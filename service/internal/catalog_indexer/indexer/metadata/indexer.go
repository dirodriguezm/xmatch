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
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Indexer struct {
	fillMetadata func(repository.InputSchema) repository.Metadata
}

func New(fillMetadata func(repository.InputSchema) repository.Metadata) *Indexer {
	return &Indexer{fillMetadata: fillMetadata}
}

func (ind *Indexer) Index(a *actor.Actor, msg actor.Message) {
	slog.Debug("Metadata Indexer Received Message")
	if msg.Error != nil {
		slog.Error("Metadata Indexer Received Error", "error", msg.Error)
		a.Broadcast(actor.Message{Error: msg.Error, Rows: nil})
		return
	}

	outputBatch := ind.getOutputBatch(msg.Rows)

	slog.Debug("Metadata Indexer Sending Message", "len", len(outputBatch))
	a.Broadcast(actor.Message{Rows: outputBatch, Error: nil})
}

func (ind *Indexer) getOutputBatch(rows []any) []any {
	outputBatch := make([]any, len(rows))
	for i := range rows {
		outputBatch[i] = ind.fillMetadata(rows[i].(repository.InputSchema))
	}
	return outputBatch
}
