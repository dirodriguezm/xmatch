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
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// Source is what the mastercat indexer needs from a decoded catalog row: its
// coordinates and the ID it is known by in the mastercat table.
type Source interface {
	GetCoordinates() (float64, float64)
	MastercatID() string
}

type Indexer[Raw Source] struct {
	mapper      *healpix.HEALPixMapper
	catalogName string
	writer      *actor.Actor[repository.Mastercat]
}

func New[Raw Source](
	cfg config.IndexerConfig,
	catalogName string,
	writer *actor.Actor[repository.Mastercat],
) (*Indexer[Raw], error) {
	slog.Debug("Creating new Mastercat Indexer")
	orderingScheme := healpix.Ring
	if strings.ToLower(cfg.OrderingScheme) == "nested" {
		orderingScheme = healpix.Nest
	}
	mapper, err := healpix.NewHEALPixMapper(cfg.Nside, orderingScheme)
	if err != nil {
		return nil, err
	}
	return &Indexer[Raw]{
		mapper:      mapper,
		catalogName: strings.ToLower(catalogName),
		writer:      writer,
	}, nil
}

func (ind *Indexer[Raw]) Index(_ *actor.Actor[Raw], msg actor.Message[Raw]) {
	slog.Debug("Mastercat Indexer Received Message")
	if msg.Error != nil {
		ind.writer.Send(actor.Message[repository.Mastercat]{Error: msg.Error})
		return
	}

	outputBatch := make([]repository.Mastercat, len(msg.Rows))
	for i := range msg.Rows {
		ra, dec := msg.Rows[i].GetCoordinates()
		ipix := ind.mapper.PixelAt(healpix.RADec(ra, dec))
		outputBatch[i] = repository.Mastercat{
			ID:   msg.Rows[i].MastercatID(),
			Ipix: ipix,
			Ra:   ra,
			Dec:  dec,
			Cat:  ind.catalogName,
		}
	}

	slog.Debug("Mastercat Indexer sending message", "len", len(outputBatch))
	ind.writer.Send(actor.Message[repository.Mastercat]{Rows: outputBatch})

	msg.Rows = nil
}
