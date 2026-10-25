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
	"github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingSchema is a row whose conversion always fails.
type failingSchema struct{ err error }

func (f failingSchema) ConvertToMetadata() (repository.Allwise, error) {
	return repository.Allwise{}, f.err
}

func TestStart(t *testing.T) {
	ctx := t.Context()
	received := make(chan actor.Message[repository.Allwise], 1)
	writer := actor.New("writer", 1, func(_ *actor.Actor[repository.Allwise], m actor.Message[repository.Allwise]) {
		received <- m
	}, nil, ctx)
	writer.Start()

	indexer := New[repository.Allwise, allwise.InputSchema](writer)
	indexerActor := actor.New("metadata indexer", 1, indexer.Index, nil, ctx)

	rows := make([]allwise.InputSchema, 10)
	for i := range 10 {
		id := "test" + strconv.Itoa(i)
		cntr := int64(i)
		rows[i] = allwise.InputSchema{Source_id: &id, Cntr: &cntr}
	}

	indexerActor.Start()
	indexerActor.Send(actor.Message[allwise.InputSchema]{Rows: rows})
	indexerActor.Stop()

	result := <-received
	writer.Stop()

	require.NoError(t, result.Error)
	require.Len(t, result.Rows, 10)
	assert.Equal(t, "test0", result.Rows[0].ID)
	assert.Equal(t, int64(0), result.Rows[0].Cntr)
	assert.Equal(t, "test9", result.Rows[9].ID)
	assert.Equal(t, int64(9), result.Rows[9].Cntr)
}

func TestStart_ForwardsInputErrorAndReturns(t *testing.T) {
	ctx := t.Context()
	received := make(chan actor.Message[repository.Allwise], 1)
	writer := actor.New("writer", 1, func(_ *actor.Actor[repository.Allwise], m actor.Message[repository.Allwise]) {
		received <- m
	}, nil, ctx)
	writer.Start()

	indexer := New[repository.Allwise, allwise.InputSchema](writer)
	testErr := errors.New("input failed")
	indexer.Index(nil, actor.Message[allwise.InputSchema]{Rows: []allwise.InputSchema{{}}, Error: testErr})

	result := <-received
	writer.Stop()
	assert.Equal(t, testErr, result.Error)
	assert.Nil(t, result.Rows)
}

func TestStart_ForwardsConversionErrorAndStopsBatch(t *testing.T) {
	ctx := t.Context()
	received := make(chan actor.Message[repository.Allwise], 1)
	writer := actor.New("writer", 1, func(_ *actor.Actor[repository.Allwise], m actor.Message[repository.Allwise]) {
		received <- m
	}, nil, ctx)
	writer.Start()

	testErr := errors.New("metadata conversion failed")
	indexer := New[repository.Allwise, failingSchema](writer)

	indexer.Index(nil, actor.Message[failingSchema]{Rows: []failingSchema{{err: testErr}}})

	result := <-received
	writer.Stop()
	assert.Equal(t, testErr, result.Error)
	assert.Nil(t, result.Rows)
}
