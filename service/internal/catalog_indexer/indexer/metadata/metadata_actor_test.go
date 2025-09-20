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

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	inbox := make(chan reader.ReaderResult)
	outbox := make(chan writer.WriterInput[repository.Allwise])
	actor := New(inbox, outbox)

	actor.Start()
	rows := make([]repository.InputSchema, 10)
	for i := range 10 {
		id := "test" + strconv.Itoa(i)
		rows[i] = repository.AllwiseInputSchema{
			Source_id: &id,
		}
	}
	inbox <- reader.ReaderResult{
		Rows:  rows,
		Error: nil,
	}
	close(inbox)

	for msg := range outbox {
		require.NoError(t, msg.Error)
		require.Len(t, msg.Rows, 10)
		for i := range 10 {
			require.Equal(t, *rows[i].(repository.AllwiseInputSchema).Source_id, msg.Rows[i].ID)
		}
	}
}
