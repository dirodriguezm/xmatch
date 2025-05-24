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
)

type Reducer struct {
	workers []*Worker
	wg      sync.WaitGroup
}

func NewReducer(workers []*Worker) *Reducer {
	return &Reducer{
		workers: workers,
	}
}

func (r *Reducer) Start() {
	slog.Debug("Reducer starting")

	r.wg = sync.WaitGroup{}
	r.wg.Add(len(r.workers))
	for _, worker := range r.workers {
		go worker.Start(&r.wg)
	}
}

func (r *Reducer) Done() {
	defer slog.Debug("Reducer done")

	r.wg.Wait()
	close(r.workers[0].outCh)
}
