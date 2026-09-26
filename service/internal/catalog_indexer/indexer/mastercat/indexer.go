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
	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type Indexer struct {
	mapper  *healpix.HEALPixMapper
	adapter catalog.CatalogAdapter
}

func New(cfg config.IndexerConfig, adapter catalog.CatalogAdapter) (*Indexer, error) {
	slog.Debug("Creating new Mastercat Indexer")
	orderingScheme := healpix.Ring
	if strings.ToLower(cfg.OrderingScheme) == "nested" {
		orderingScheme = healpix.Nest
	}
	mapper, err := healpix.NewHEALPixMapper(cfg.Nside, orderingScheme)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		mapper:  mapper,
		adapter: adapter,
	}, nil
}

func (ind Indexer) Index(a *actor.Actor, msg actor.Message) {
	slog.Debug("Mastercat Indexer Received Message")
	if msg.Error != nil {
		a.Broadcast(actor.Message{
			Error: msg.Error,
			Rows:  nil,
		})
		return
	}

	outputBatch := make([]any, len(msg.Rows))
	for i := range msg.Rows {
		ra, dec, err := ind.adapter.GetCoordinates(msg.Rows[i])
		if err != nil {
			a.Broadcast(actor.Message{Error: err, Rows: nil})
			return
		}
		ipix := ind.mapper.PixelAt(healpix.RADec(ra, dec))
		mastercat, err := ind.adapter.ConvertToMastercat(msg.Rows[i], ipix)
		if err != nil {
			a.Broadcast(actor.Message{Error: err, Rows: nil})
			return
		}
		outputBatch[i] = mastercat
	}

	slog.Debug("Mastercat Indexer sending message", "len", len(outputBatch))
	a.Broadcast(actor.Message{
		Rows:  outputBatch,
		Error: nil,
	})

	msg.Rows = nil
}
