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
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
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

func TestIndexActor(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)

	rows := []any{
		TestSchema{Ra: 1.0, Dec: 2.0, ID: "id1", Cat: "test"},
		TestSchema{Ra: 3.0, Dec: 4.0, ID: "id2", Cat: "test"},
	}
	for _, row := range rows {
		schema := row.(TestSchema)
		expectedIPix := mapper.PixelAt(healpix.RADec(schema.Ra, schema.Dec))
		adapter.EXPECT().GetCoordinates(row).Return(schema.Ra, schema.Dec, nil).Once()
		adapter.EXPECT().ConvertToMastercat(row, expectedIPix).Return(repository.Mastercat{
			ID:   schema.ID,
			Ipix: expectedIPix,
			Ra:   schema.Ra,
			Dec:  schema.Dec,
			Cat:  "test",
		}, nil).Once()
	}

	indexer, err := New(config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, adapter)
	require.NoError(t, err)
	ctx := t.Context()

	results := make([]any, 0)
	receiver := actor.New(
		"receiver",
		2,
		func(a *actor.Actor, m actor.Message) {
			results = append(results, m.Rows...)
		},
		nil,
		nil,
		ctx,
	)
	a := actor.New("mastercat indexer", 2, indexer.Index, nil, []*actor.Actor{receiver}, ctx)

	receiver.Start()
	a.Start()
	a.Send(actor.Message{Rows: rows, Error: nil})

	a.Stop()
	receiver.Stop()

	require.Len(t, results, 2)
	assert.Equal(t, "id1", results[0].(repository.Mastercat).ID)
	assert.Equal(t, "id2", results[1].(repository.Mastercat).ID)
}

func TestIndexActor_BroadcastsInputErrorAndReturns(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	indexer, err := New(config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, adapter)
	require.NoError(t, err)

	received := make(chan actor.Message, 1)
	receiver := actor.New("receiver", 1, func(_ *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, t.Context())
	receiver.Start()

	a := actor.New("mastercat indexer", 1, indexer.Index, nil, []*actor.Actor{receiver}, t.Context())
	testErr := errors.New("input failed")
	indexer.Index(a, actor.Message{Rows: []any{TestSchema{}}, Error: testErr})
	receiver.Stop()

	msg := <-received
	assert.Equal(t, testErr, msg.Error)
	assert.Nil(t, msg.Rows)
}

func TestIndexActor_BroadcastsCoordinateErrorAndStopsBatch(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	row := TestSchema{Ra: 1.0, Dec: 2.0, ID: "id1"}
	testErr := errors.New("coordinates failed")
	adapter.EXPECT().GetCoordinates(row).Return(0.0, 0.0, testErr).Once()

	indexer, err := New(config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, adapter)
	require.NoError(t, err)

	received := make(chan actor.Message, 1)
	receiver := actor.New("receiver", 1, func(_ *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, t.Context())
	receiver.Start()

	a := actor.New("mastercat indexer", 1, indexer.Index, nil, []*actor.Actor{receiver}, t.Context())
	indexer.Index(a, actor.Message{Rows: []any{row}, Error: nil})
	receiver.Stop()

	msg := <-received
	assert.Equal(t, testErr, msg.Error)
	assert.Nil(t, msg.Rows)
}

func TestIndexActor_BroadcastsConversionErrorAndStopsBatch(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)

	row := TestSchema{Ra: 1.0, Dec: 2.0, ID: "id1"}
	ipix := mapper.PixelAt(healpix.RADec(row.Ra, row.Dec))
	testErr := errors.New("conversion failed")
	adapter.EXPECT().GetCoordinates(row).Return(row.Ra, row.Dec, nil).Once()
	adapter.EXPECT().ConvertToMastercat(row, ipix).Return(repository.Mastercat{}, testErr).Once()

	indexer, err := New(config.IndexerConfig{OrderingScheme: "nested", Nside: 18}, adapter)
	require.NoError(t, err)

	received := make(chan actor.Message, 1)
	receiver := actor.New("receiver", 1, func(_ *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, t.Context())
	receiver.Start()

	a := actor.New("mastercat indexer", 1, indexer.Index, nil, []*actor.Actor{receiver}, t.Context())
	indexer.Index(a, actor.Message{Rows: []any{row}, Error: nil})
	receiver.Stop()

	msg := <-received
	assert.Equal(t, testErr, msg.Error)
	assert.Nil(t, msg.Rows)
}
