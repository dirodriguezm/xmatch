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
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	indexer := Indexer{
		fillMetadata: repository.FillAllwiseMetadata,
	}
	result := make([]actor.Message, 0)
	ctx := t.Context()
	testActor := actor.New(1, func(a *actor.Actor, m actor.Message) {
		result = append(result, m)
	}, nil, nil, ctx)
	indexerActor := actor.New(1, indexer.Index, nil, []*actor.Actor{testActor}, ctx)

	testActor.Start()
	indexerActor.Start()

	rows := make([]any, 10)
	for i := range 10 {
		id := "test" + strconv.Itoa(i)
		cntr := int64(i)
		rows[i] = repository.AllwiseInputSchema{Source_id: &id, Cntr: &cntr}
	}
	indexerActor.Send(actor.Message{Rows: rows, Error: nil})

	indexerActor.Stop()
	testActor.Stop()

	require.Len(t, result, 1)
}
