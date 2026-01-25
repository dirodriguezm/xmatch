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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type TestSchema struct {
	Ra  float64
	Dec float64
	ID  string
	Cat string
}

func FillMastercat(schema repository.InputSchema, ipix int64) repository.Mastercat {
	return repository.Mastercat{
		ID:   schema.GetId(),
		Ipix: ipix,
		Ra:   schema.(TestSchema).Ra,
		Dec:  schema.(TestSchema).Dec,
		Cat:  "test",
	}
}

func FillMetadata(schema repository.InputSchema) repository.Metadata {
	return nil
}

func (t TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t TestSchema) GetId() string {
	return t.ID
}

func TestIndexActor(t *testing.T) {
	rows := make([]any, 2)
	rows[0] = TestSchema{Ra: 0.0, Dec: 0.0, ID: "id1", Cat: "test"}
	rows[1] = TestSchema{Ra: 0.0, Dec: 0.0, ID: "id2", Cat: "test"}

	indexer, err := New(config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, FillMastercat)
	require.NoError(t, err)
	ctx := t.Context()

	results := make([]any, 0)
	receiver := actor.New(
		2,
		func(a *actor.Actor, m actor.Message) {
			results = append(results, m.Rows...)
		},
		nil,
		nil,
		ctx,
	)
	a := actor.New(2, indexer.Index, nil, []*actor.Actor{receiver}, ctx)

	receiver.Start()
	a.Start()
	a.Send(actor.Message{Rows: rows, Error: nil})

	a.Stop()
	receiver.Stop()

	require.Len(t, results, 2)
}
