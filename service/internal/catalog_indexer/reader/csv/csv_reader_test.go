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

package csv_reader

import (
	"io"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`

	source, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:" + csv,
		Type:        "csv",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	csvReader, err := NewCsvReader(source, WithCsvBatchSize(2))
	require.NoError(t, err)

	rows, err := csvReader.Read()
	require.NoError(t, err)

	require.Equal(t, 3, len(rows))

	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3)
	for i, row := range rows {
		mastercat := repository.Mastercat{}
		row.(repository.InputSchema).FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}

	require.Equal(t, expectedOids, receivedOids)
}

func TestReadWithHeader(t *testing.T) {
	csv := `o1,1,1
o2,2,2
o3,3,3
`

	source, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:" + csv,
		Type:        "csv",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	csvReader, err := NewCsvReader(source, WithHeader([]string{"oid", "ra", "dec"}), WithCsvBatchSize(2))
	require.NoError(t, err)

	rows, err := csvReader.Read()
	require.NoError(t, err)

	require.Equal(t, 3, len(rows))

	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3)
	for i, row := range rows {
		mastercat := repository.Mastercat{}
		row.(repository.InputSchema).FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}

	require.Equal(t, expectedOids, receivedOids)
}

func TestReadWithHeader_Error(t *testing.T) {
	source := source.Source{
		Sources:     []string{"buffer:"},
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, WithCsvBatchSize(2))
	require.NoError(t, err)

	rows, err := csvReader.Read()
	require.Error(t, err)
	require.Nil(t, rows)
}

func TestReadBatch(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
o4,4,4
`

	source := source.Source{
		Sources:     []string{"buffer:" + csv},
		CatalogName: "test",
	}

	csvReader, err := NewCsvReader(
		&source,
		WithCsvBatchSize(2),
		WithFirstLineHeader(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	var rows []repository.InputSchema

	for {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}

		if err == io.EOF {
			t.Log("EOF")
			break
		}

		require.Len(t, batch, 2)

		for _, row := range batch {
			rows = append(rows, row.(repository.InputSchema))
		}
	}

	require.Equal(t, 4, len(rows))
	expectedOids := []string{"o1", "o2", "o3", "o4"}
	receivedOids := make([]string, 4)
	for i, row := range rows {
		mastercat := repository.Mastercat{}
		row.FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadBatchWithDelta(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`

	source := source.Source{
		Sources:     []string{"buffer:" + csv},
		CatalogName: "test",
	}

	csvReader, err := NewCsvReader(
		&source,
		WithFirstLineHeader(true),
		WithCsvBatchSize(2),
	)
	if err != nil {
		t.Fatal(err)
	}

	var rows []repository.InputSchema

	eof := false
	for !eof {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			eof = true
		}
		for _, row := range batch {
			rows = append(rows, row.(repository.InputSchema))
		}
	}

	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3)
	for i, row := range rows {
		mastercat := repository.Mastercat{}
		row.FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadMultipleReaders(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`

	source := source.Source{
		Sources:     []string{"buffer:" + csv, "buffer:" + csv},
		CatalogName: "test",
	}

	csvReader, err := NewCsvReader(
		&source,
		WithFirstLineHeader(true),
		WithCsvBatchSize(2),
	)
	require.NoError(t, err)

	var rows []repository.InputSchema

	eof := false
	for !eof {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			eof = true
		}

		for _, row := range batch {
			rows = append(rows, row.(repository.InputSchema))
		}
	}

	require.Equal(t, 6, len(rows))

	expectedOids := []string{"o1", "o2", "o3", "o1", "o2", "o3"}
	receivedOids := make([]string, 6)
	for i, row := range rows {
		mastercat := repository.Mastercat{}
		row.FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}

	require.Equal(t, expectedOids, receivedOids)
}
