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

package reducer

import (
	"log/slog"
	"sync"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	partition_reader "github.com/dirodriguezm/xmatch/service/internal/preprocessor/reader"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Worker struct {
	inCh   <-chan partition_reader.Records
	outCh  chan<- repository.InputSchema
	schema config.ParquetSchema
}

func NewWorker(inCh <-chan partition_reader.Records, outCh chan<- repository.InputSchema, schema config.ParquetSchema) *Worker {
	return &Worker{
		inCh:   inCh,
		outCh:  outCh,
		schema: schema,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	slog.Debug("Reducer Worker starting")
	defer wg.Done()

	for records := range w.inCh {
		slog.Debug("Reducer sending message")
		w.outCh <- w.getObject(records)
	}
}

func (w *Worker) getObject(records partition_reader.Records) repository.InputSchema {
	switch w.schema {
	case config.VlassSchema:
		return w.getVlassObject(records)
	default:
		panic("Unknown schema")
	}
}

func (w *Worker) getVlassObject(records partition_reader.Records) repository.InputSchema {
	meanRa := float64(0)
	meanDec := float64(0)
	errRa := float64(0)
	errDec := float64(0)
	meanFlux := float64(0)
	meanFluxErr := float64(0)

	for _, r := range records {
		ra, dec := r.GetCoordinates()
		era := r.(*repository.VlassInputSchema).ERA
		edec := r.(*repository.VlassInputSchema).EDEC
		flux := r.(*repository.VlassInputSchema).TotalFlux
		eflux := r.(*repository.VlassInputSchema).ETotalFlux

		meanRa += ra
		meanDec += dec
		if era != nil {
			errRa += *era
		}
		if edec != nil {
			errDec += *edec
		}
		if flux != nil {
			meanFlux += *flux
		}
		if eflux != nil {
			meanFluxErr += *eflux
		}
	}
	meanRa /= float64(len(records))
	meanDec /= float64(len(records))
	errRa /= float64(len(records))
	errDec /= float64(len(records))
	meanFlux /= float64(len(records))
	meanFluxErr /= float64(len(records))

	id := records[0].GetId()
	return &repository.VlassObjectSchema{
		Id:    &id,
		Ra:    &meanRa,
		Dec:   &meanDec,
		Era:   &errRa,
		Edec:  &errDec,
		Flux:  &meanFlux,
		EFlux: &meanFluxErr,
	}
}
