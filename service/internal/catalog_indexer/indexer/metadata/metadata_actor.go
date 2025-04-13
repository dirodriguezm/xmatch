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

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
)

type IndexerActor struct {
	inbox  chan reader.ReaderResult
	outbox chan writer.WriterInput[any]
}

func New(
	inbox chan reader.ReaderResult,
	outbox chan writer.WriterInput[any],
) *IndexerActor {
	slog.Debug("Creating new Indexer")
	return &IndexerActor{
		inbox:  inbox,
		outbox: outbox,
	}
}

// Starts the Actor goroutine
//
// Listen for messages in the inbox channel and process them
// sending the results to the outbox channel
func (actor *IndexerActor) Start() {
	slog.Debug("Starting Indexer")
	go func() {
		defer func() {
			close(actor.outbox)
			slog.Debug("Closing Indexer")
		}()
		for msg := range actor.inbox {
			actor.receive(msg)
		}
	}()
}

func (actor *IndexerActor) receive(msg reader.ReaderResult) {
	slog.Debug("Indexer Received Message")
	if msg.Error != nil {
		actor.outbox <- writer.WriterInput[any]{
			Error: msg.Error,
			Rows:  nil,
		}
		return
	}

	objects := make([]any, len(msg.Rows))
	for i := 0; i < len(msg.Rows); i++ {
		object := msg.Rows[i]
		objects[i] = object.ToMetadata()
	}

	actor.outbox <- writer.WriterInput[any]{
		Error: nil,
		Rows:  objects,
	}
}
