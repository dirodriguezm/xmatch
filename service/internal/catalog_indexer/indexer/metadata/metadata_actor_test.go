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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type TestInputSchema struct {
	Id  string
	Ra  float64
	Dec float64
}

func (t TestInputSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t TestInputSchema) ToMetadata() any {
	return t
}

func (t TestInputSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	return repository.ParquetMastercat{}
}

func (t TestInputSchema) SetField(string, any) {}

func (t TestInputSchema) GetId() string {
	return t.Id
}

func TestStart(t *testing.T) {
	inbox := make(chan reader.ReaderResult)
	outbox := make(chan writer.WriterInput[any])
	actor := New(inbox, outbox)

	actor.Start()
	rows := make([]repository.InputSchema, 10)
	for i := 0; i < 10; i++ {
		rows[i] = TestInputSchema{
			Id:  "test",
			Ra:  float64(i),
			Dec: float64(i + 1),
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
		for i := 0; i < 10; i++ {
			require.Equal(t, rows[i].(TestInputSchema).Id, msg.Rows[i].(TestInputSchema).Id)
			require.Equal(t, rows[i].(TestInputSchema).Ra, msg.Rows[i].(TestInputSchema).Ra)
			require.Equal(t, rows[i].(TestInputSchema).Dec, msg.Rows[i].(TestInputSchema).Dec)
		}
	}
}
