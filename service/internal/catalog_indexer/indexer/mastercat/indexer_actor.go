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

package mastercat_indexer

import (
	"log/slog"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type IndexerActor struct {
	source *source.Source
	mapper *healpix.HEALPixMapper
	inbox  chan reader.ReaderResult
	outbox chan writer.WriterInput[any]
}

func New(
	src *source.Source,
	inbox chan reader.ReaderResult,
	outbox chan writer.WriterInput[any],
	cfg *config.IndexerConfig,
) (*IndexerActor, error) {
	slog.Debug("Creating new Indexer")
	orderingScheme := healpix.Ring
	if strings.ToLower(cfg.OrderingScheme) == "nested" {
		orderingScheme = healpix.Nest
	}
	mapper, err := healpix.NewHEALPixMapper(src.Nside, orderingScheme)
	if err != nil {
		return nil, err
	}
	return &IndexerActor{
		source: src,
		mapper: mapper,
		inbox:  inbox,
		outbox: outbox,
	}, nil
}

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
			Rows:  nil,
			Error: msg.Error,
		}
		return
	}

	outputBatch := make([]any, len(msg.Rows))
	for i, row := range msg.Rows {
		ra, dec := row.GetCoordinates()
		point := healpix.RADec(ra, dec)
		ipix := actor.mapper.PixelAt(point)
		outputBatch[i] = row.ToMastercat(ipix)
	}

	slog.Debug("Indexer sending message")
	actor.outbox <- writer.WriterInput[any]{
		Rows:  outputBatch,
		Error: nil,
	}
}
