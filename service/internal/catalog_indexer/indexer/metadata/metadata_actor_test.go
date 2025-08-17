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

type TestMetadata struct {
	Id  string
	Ra  float64
	Dec float64
}

func (t TestInputSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t TestInputSchema) FillMetadata(dst repository.Metadata) {
	dst.(*TestMetadata).Id = t.Id
	dst.(*TestMetadata).Ra = t.Ra
	dst.(*TestMetadata).Dec = t.Dec
}

func (t TestInputSchema) FillMastercat(dst *repository.Mastercat, ipix int64) {
	dst.ID = t.Id
	dst.Ra = t.Ra
	dst.Dec = t.Dec
	dst.Cat = "test"
	dst.Ipix = ipix
}

func (t TestInputSchema) GetId() string {
	return t.Id
}

func (t TestMetadata) GetId() string {
	return t.Id
}

type TestMetadataParser struct{}

func (p TestMetadataParser) Parse(in repository.InputSchema) repository.Metadata {
	testMetadata := TestMetadata{}
	in.FillMetadata(&testMetadata)
	return testMetadata
}

func TestStart(t *testing.T) {
	inbox := make(chan reader.ReaderResult)
	outbox := make(chan writer.WriterInput[repository.Metadata])
	actor := New(inbox, outbox, TestMetadataParser{})

	actor.Start()
	rows := make([]repository.InputSchema, 10)
	for i := range 10 {
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
		for i := range 10 {
			require.Equal(t, rows[i].(TestInputSchema).Id, msg.Rows[i].(TestMetadata).Id)
			require.Equal(t, rows[i].(TestInputSchema).Ra, msg.Rows[i].(TestMetadata).Ra)
			require.Equal(t, rows[i].(TestInputSchema).Dec, msg.Rows[i].(TestMetadata).Dec)
		}
	}
}
