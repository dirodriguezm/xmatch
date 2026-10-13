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
)

// Converter is what the metadata indexer needs from a decoded catalog row: the
// row builds its own metadata model.
type Converter[Md any] interface {
	ConvertToMetadata() (Md, error)
}

type Indexer[Md any, Raw Converter[Md]] struct {
	writer *actor.Actor[Md]
}

func New[Md any, Raw Converter[Md]](writer *actor.Actor[Md]) *Indexer[Md, Raw] {
	return &Indexer[Md, Raw]{writer: writer}
}

func (ind *Indexer[Md, Raw]) Index(_ *actor.Actor[Raw], msg actor.Message[Raw]) {
	slog.Debug("Metadata Indexer Received Message")
	if msg.Error != nil {
		slog.Error("Metadata Indexer Received Error", "error", msg.Error)
		ind.writer.Send(actor.Message[Md]{Error: msg.Error})
		return
	}

	outputBatch, err := ind.getOutputBatch(msg.Rows)
	if err != nil {
		ind.writer.Send(actor.Message[Md]{Error: err})
		return
	}

	slog.Debug("Metadata Indexer Sending Message", "len", len(outputBatch))
	ind.writer.Send(actor.Message[Md]{Rows: outputBatch})
}

func (ind *Indexer[Md, Raw]) getOutputBatch(rows []Raw) ([]Md, error) {
	outputBatch := make([]Md, len(rows))
	for i := range rows {
		md, err := rows[i].ConvertToMetadata()
		if err != nil {
			return nil, err
		}
		outputBatch[i] = md
	}
	return outputBatch, nil
}
