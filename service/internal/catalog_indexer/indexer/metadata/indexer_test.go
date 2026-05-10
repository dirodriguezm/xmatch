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
	"errors"
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	indexer := New(adapter)
	received := make(chan actor.Message, 1)
	ctx := t.Context()
	testActor := actor.New("receiver", 1, func(a *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, ctx)

	testActor.Start()

	rows := make([]any, 10)
	for i := range 10 {
		id := "test" + strconv.Itoa(i)
		cntr := int64(i)
		row := allwise.InputSchema{Source_id: &id, Cntr: &cntr}
		rows[i] = row
		adapter.EXPECT().ConvertToMetadataFromRaw(row).Return(repository.Allwise{ID: id}, nil).Once()
	}

	indexerActor := actor.New("metadata indexer", 1, indexer.Index, nil, []*actor.Actor{testActor}, ctx)
	indexer.Index(indexerActor, actor.Message{Rows: rows, Error: nil})
	result := <-received
	testActor.Stop()

	require.NoError(t, result.Error)
	require.Len(t, result.Rows, 10)
	assert.Equal(t, repository.Allwise{ID: "test0"}, result.Rows[0])
	assert.Equal(t, repository.Allwise{ID: "test9"}, result.Rows[9])
}

func TestStart_BroadcastsInputErrorAndReturns(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	indexer := New(adapter)
	received := make(chan actor.Message, 1)
	ctx := t.Context()
	testActor := actor.New("receiver", 1, func(a *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, ctx)
	testActor.Start()

	indexerActor := actor.New("metadata indexer", 1, indexer.Index, nil, []*actor.Actor{testActor}, ctx)
	testErr := errors.New("input failed")
	indexer.Index(indexerActor, actor.Message{Rows: []any{allwise.InputSchema{}}, Error: testErr})
	result := <-received
	testActor.Stop()

	assert.Equal(t, testErr, result.Error)
	assert.Nil(t, result.Rows)
}

func TestStart_BroadcastsConversionErrorAndStopsBatch(t *testing.T) {
	adapter := catalog.NewMockCatalogAdapter(t)
	indexer := New(adapter)
	received := make(chan actor.Message, 1)
	ctx := t.Context()
	testActor := actor.New("receiver", 1, func(a *actor.Actor, m actor.Message) {
		received <- m
	}, nil, nil, ctx)
	testActor.Start()

	row := allwise.InputSchema{}
	testErr := errors.New("metadata conversion failed")
	adapter.EXPECT().ConvertToMetadataFromRaw(row).Return(nil, testErr).Once()

	indexerActor := actor.New("metadata indexer", 1, indexer.Index, nil, []*actor.Actor{testActor}, ctx)
	indexer.Index(indexerActor, actor.Message{Rows: []any{row}, Error: nil})
	result := <-received
	testActor.Stop()

	assert.Equal(t, testErr, result.Error)
	assert.Nil(t, result.Rows)
}
