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
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type TestData struct {
	Oid string
	Ra  float64
	Dec float64
}

type TestFixture struct {
	source       *source.SourceBuilder
	reader       *ReaderBuilder
	expectedRows []TestData
}

func setUpTestFixture(t *testing.T) *TestFixture {
	t.Helper()

	fileContent := `
oid,ra,dec
o1,1,1
o2,2,2
`
	nFiles := 5
	testData := make([]string, nFiles, nFiles)
	for i := range nFiles {
		testData[i] = fileContent
	}
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	return &TestFixture{
		source: source.ASource(t).WithUrl(url).WithCsvFiles(testData),
		reader: AReader(t),
		expectedRows: []TestData{
			{"o1", 1, 1},
			{"o2", 2, 2},
		},
	}
}

func TestReadMultipleFiles_Csv(t *testing.T) {
	// arrange
	fixture := setUpTestFixture(t)

	// create source
	source := fixture.source.Build()

	// create reader
	reader := fixture.reader.WithSource(source).Build()

	// act
	reader.Start()

	// collect results
	allRows := make([]repository.InputSchema, 0)
	var err error
	for i := range fixture.reader.OutputChannel {
		for msg := range fixture.reader.OutputChannel[i] {
			if msg.Error != nil {
				err = msg.Error
				break
			}
			allRows = append(allRows, msg.Rows...)
		}
		require.NoError(t, err)
	}

	// assert
	require.Len(t, allRows, 10)
	for i, row := range allRows {
		expectedData := fixture.expectedRows[i%2]
		require.Equal(t, expectedData.Oid, *row.ToMastercat(0).ID)
		require.Equal(t, expectedData.Ra, *row.ToMastercat(0).Ra)
		require.Equal(t, expectedData.Dec, *row.ToMastercat(0).Dec)
	}
}

func TestReadNestedFiles_Csv(t *testing.T) {
	// arrange
	fixture := setUpTestFixture(t)

	// create source
	fileContent := `
oid,ra,dec
o1,1,1
o2,2,2
`
	nFiles := 5
	testData := make([]string, nFiles, nFiles)
	for i := range nFiles {
		testData[i] = fileContent
	}
	source := fixture.source.WithNestedCsvFiles(testData, testData).Build()
	require.Len(t, source.Reader, 10)

	// create reader
	reader := fixture.reader.WithSource(source).Build()

	// act
	reader.Start()

	// collect results
	allRows := make([]repository.InputSchema, 0)
	for i := range fixture.reader.OutputChannel {
		for msg := range fixture.reader.OutputChannel[i] {
			require.NoError(t, msg.Error)
			for _, row := range msg.Rows {
				allRows = append(allRows, row)
			}
		}
	}

	// assert
	require.Len(t, allRows, 20)
	for i, row := range allRows {
		expectedData := fixture.expectedRows[i%2]
		require.Equal(t, expectedData.Oid, *row.ToMastercat(0).ID)
		require.Equal(t, expectedData.Ra, *row.ToMastercat(0).Ra)
		require.Equal(t, expectedData.Dec, *row.ToMastercat(0).Dec)
	}
}
