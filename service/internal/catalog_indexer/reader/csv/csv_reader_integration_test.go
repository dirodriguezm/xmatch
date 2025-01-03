package csv_reader

import (
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/stretchr/testify/require"
)

type TestData struct {
	Oid string
	Ra  string
	Dec string
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
			{"o1", "1", "1"},
			{"o2", "2", "2"},
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
	allRows := make([]indexer.Row, 0)
	var err error
	for msg := range fixture.reader.OutputChannel {
		if msg.Error != nil {
			err = msg.Error
			break
		}
		allRows = append(allRows, msg.Rows...)
	}
	require.NoError(t, err)

	// assert
	require.Len(t, allRows, 10)
	for i, row := range allRows {
		expectedData := fixture.expectedRows[i%2]
		require.Equal(t, expectedData.Oid, row["oid"])
		require.Equal(t, expectedData.Ra, row["ra"])
		require.Equal(t, expectedData.Dec, row["dec"])
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
	allRows := make([]indexer.Row, 0)
	for msg := range fixture.reader.OutputChannel {
		for _, row := range msg.Rows {
			allRows = append(allRows, row)
		}
	}

	// assert
	require.Len(t, allRows, 20)
	for i, row := range allRows {
		expectedData := fixture.expectedRows[i%2]
		require.Equal(t, expectedData.Oid, row["oid"])
		require.Equal(t, expectedData.Ra, row["ra"])
		require.Equal(t, expectedData.Dec, row["dec"])
	}
}
