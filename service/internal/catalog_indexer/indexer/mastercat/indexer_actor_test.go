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

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
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

// implement the interface
func (t *TestSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	return repository.ParquetMastercat{
		ID:   &t.ID,
		Ra:   &t.Ra,
		Dec:  &t.Dec,
		Cat:  &t.Cat,
		Ipix: &ipix,
	}
}

// implement the interface
func (t *TestSchema) ToMetadata() any {
	return t
}

func (t *TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t *TestSchema) SetField(name string, val interface{}) {}

func (t *TestSchema) GetId() string {
	return t.ID
}

func TestIndexActor(t *testing.T) {
	inbox := make(chan reader.ReaderResult)
	outbox := make(chan writer.WriterInput[any])
	rows := make([]repository.InputSchema, 2)
	rows[0] = &TestSchema{Ra: 0.0, Dec: 0.0, ID: "id1", Cat: "test"}
	rows[1] = &TestSchema{Ra: 0.0, Dec: 0.0, ID: "id2", Cat: "test"}

	src, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:",
		Type:        "csv",
		CatalogName: "catalog",
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
	})
	require.NoError(t, err)
	indexerActor, err := New(src, inbox, outbox, &config.IndexerConfig{OrderingScheme: "nested", Nside: 18})
	require.NoError(t, err)

	indexerActor.Start()
	inbox <- reader.ReaderResult{Rows: rows, Error: nil}
	close(inbox)
	results := make([]repository.ParquetMastercat, 2)
	for msg := range outbox {
		for i, obj := range msg.Rows {
			results[i] = obj.(repository.ParquetMastercat)
		}
	}

	require.Len(t, results, 2)
}
