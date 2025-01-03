package parquet_reader

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/stretchr/testify/require"
)

var metadata = []string{
	"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=ra, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
	"name=dec, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
}

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

func setUpTestFixture_Parquet(t *testing.T) *TestFixture {
	t.Helper()

	fileContent := [][]string{
		{"o1", "1", "1"},
		{"o2", "2", "2"},
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
			WithParquetFiles(metadata, testData),
		reader: AReader(t).WithType("parquet"),
		expectedRows: []TestData{
			{"o1", "1", "1"},
			{"o2", "2", "2"},
		},
	}
}

func TestReadMultipleFiles_Parquet(t *testing.T) {
	// arrange
	fixture := setUpTestFixture_Parquet(t)

	// create source
	source := fixture.source.Build()

	// create r
	r := fixture.reader.WithSource(source).WithParquetMetadata(metadata).Build()

	// act
	r.Start()

	// collect results
	allRows := make([]indexer.Row, 0)
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
		require.Equal(t, expectedData.Oid, row["oid"])
		require.Equal(t, expectedData.Ra, row["ra"])
		require.Equal(t, expectedData.Dec, row["dec"])
	}
}

func TestReadWithDefaultMetadata(t *testing.T) {
	// uses default metadata to write but not to read
	metadata := []string{
		"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY",
		"name=ra, type=DOUBLE",
		"name=dec, type=DOUBLE",
	}

	// create source
	fileContent := [][]string{
		{"o1", "1", "1"},
		{"o2", "2", "2"},
	}
	nFiles := 5
	testData := make([][][]string, nFiles)
	for i := 0; i < nFiles; i++ {
		testData[i] = fileContent
	}
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	src := source.
		ASource(t).
		WithType("parquet").
		WithUrl(url).
		WithParquetFiles(metadata, testData).
		Build()

	// create reader
	// note we don't specify metadata here
	readerBuilder := AReader(t).WithType("parquet").WithSource(src)
	r := readerBuilder.Build()

	// act
	r.Start()

	// collect results
	allRows := make([]indexer.Row, 0)
	var err error
	for msg := range readerBuilder.OutputChannel {
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
	expectedRows := []struct {
		Oid string
		Ra  float64
		Dec float64
	}{{"o1", 1, 1}, {"o2", 2, 2}}
	for i, row := range allRows {
		expectedData := expectedRows[i%2]
		require.Equal(t, expectedData.Oid, row["oid"])
		require.Equal(t, expectedData.Ra, row["ra"])
		require.Equal(t, expectedData.Dec, row["dec"])
	}
}
