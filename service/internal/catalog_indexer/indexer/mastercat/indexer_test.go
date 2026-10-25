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
	"errors"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSchema struct {
	Ra  float64
	Dec float64
	ID  string
	Cat string
}

func (r TestSchema) GetCoordinates() (float64, float64) { return r.Ra, r.Dec }

func (r TestSchema) MastercatID() string { return r.ID }

func TestIndexActor(t *testing.T) {
	ctx := t.Context()

	rows := []TestSchema{
		{Ra: 1.0, Dec: 2.0, ID: "id1", Cat: "test"},
		{Ra: 3.0, Dec: 4.0, ID: "id2", Cat: "test"},
	}

	received := make(chan actor.Message[repository.Mastercat], 1)
	receiver := actor.New("receiver", 2, func(_ *actor.Actor[repository.Mastercat], m actor.Message[repository.Mastercat]) {
		received <- m
	}, nil, ctx)
	receiver.Start()

	indexer, err := New[TestSchema](config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, "TeSt", receiver)
	require.NoError(t, err)

	a := actor.New("mastercat indexer", 2, indexer.Index, nil, ctx)
	a.Start()
	a.Send(actor.Message[TestSchema]{Rows: rows, Error: nil})
	a.Stop()

	msg := <-received
	receiver.Stop()

	require.NoError(t, msg.Error)
	require.Len(t, msg.Rows, 2)

	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)
	for i, row := range rows {
		assert.Equal(t, row.ID, msg.Rows[i].ID)
		assert.Equal(t, row.Ra, msg.Rows[i].Ra)
		assert.Equal(t, row.Dec, msg.Rows[i].Dec)
		assert.Equal(t, mapper.PixelAt(healpix.RADec(row.Ra, row.Dec)), msg.Rows[i].Ipix)
		assert.Equal(t, "test", msg.Rows[i].Cat, "cat should be the lowercased configured catalog name")
	}
}

func TestIndexActor_ForwardsInputErrorAndReturns(t *testing.T) {
	ctx := t.Context()

	received := make(chan actor.Message[repository.Mastercat], 1)
	receiver := actor.New("receiver", 1, func(_ *actor.Actor[repository.Mastercat], m actor.Message[repository.Mastercat]) {
		received <- m
	}, nil, ctx)
	receiver.Start()

	indexer, err := New[TestSchema](config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, "test", receiver)
	require.NoError(t, err)

	testErr := errors.New("input failed")
	indexer.Index(nil, actor.Message[TestSchema]{Rows: []TestSchema{{}}, Error: testErr})

	msg := <-received
	receiver.Stop()
	assert.Equal(t, testErr, msg.Error)
	assert.Nil(t, msg.Rows)
}
