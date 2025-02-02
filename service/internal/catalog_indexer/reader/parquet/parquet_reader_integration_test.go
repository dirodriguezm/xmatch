package parquet_reader

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

var metadata = []string{
	"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=ra, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=dec, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=mag, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
}

var metadata2 = []string{
	"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=ra, type=DOUBLE",
	"name=dec, type=DOUBLE",
	"name=mag, type=DOUBLE",
}

type TestData struct {
	Oid string `parquet:"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ra  string `parquet:"name=ra, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Dec string `parquet:"name=dec, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

type TestData2 struct {
	Oid string  `parquet:"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
}

type TestFixture struct {
	source       *source.SourceBuilder
	reader       *ReaderBuilder[TestData2]
	expectedRows []TestData2
}

func setUpTestFixture_Parquet(t *testing.T) *TestFixture {
	t.Helper()

	fileContent := [][]string{
		{"o1", "1", "1", "1"},
		{"o2", "2", "2", "2"},
	}
	nFiles := 5
	testData := make([][][]string, nFiles)
	for i := 0; i < nFiles; i++ {
		testData[i] = fileContent
	}
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	return &TestFixture{
		source: source.
			ASource(t).
			WithType("parquet").
			WithUrl(url).
			WithParquetFiles(metadata2, testData),
		reader: AReader[TestData2](t).WithType("parquet"),
		expectedRows: []TestData2{
			{"o1", 1, 1},
			{"o2", 2, 2},
		},
	}
}

func TestReadMultipleFiles_Parquet(t *testing.T) {
	// arrange
	fixture := setUpTestFixture_Parquet(t)

	// create source
	source := fixture.source.Build()

	// create r
	r := fixture.reader.WithSource(source).Build()

	// act
	r.Start()

	// collect results
	allRows := make([]repository.InputSchema, 0)
	var err error
	for msg := range fixture.reader.OutputChannel {
		if msg.Error != nil {
			err = msg.Error
			if errors.As(err, &reader.ReadError{}) {
				slog.Error("Error reading parquet", "source", err.(reader.ReadError).Source)
			}
			break
		}
		allRows = append(allRows, msg.Rows...)
	}
	require.NoError(t, err)

	// assert
	require.Len(t, allRows, 10)
	for i, row := range allRows {
		expectedData := fixture.expectedRows[i%2]
		require.Equal(t, expectedData.Oid, *row.ToMastercat().ID)
		require.Equal(t, expectedData.Ra, *row.ToMastercat().Ra)
		require.Equal(t, expectedData.Dec, *row.ToMastercat().Dec)
	}
}
