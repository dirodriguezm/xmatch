package parquet_reader

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type ObjectWrite struct {
	Oid string  `parquet:"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
	Mag float64 `parquet:"name=mag, type=DOUBLE"`
}

type Object struct {
	Oid string  `parquet:"name=oid, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ra  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
}

func Write(t *testing.T, nrows int) string {
	dir := t.TempDir()
	parquetFile := filepath.Join(dir, "test.parquet")
	var err error
	fw, err := local.NewLocalFileWriter(parquetFile)
	if err != nil {
		t.Fatal("Can't create local file", err)
	}

	//write
	pw, err := writer.NewParquetWriter(fw, new(Object), 4)
	if err != nil {
		t.Fatal("Can't create parquet writer", err)
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY
	for i := 0; i < nrows; i++ {
		obj := ObjectWrite{
			Oid: fmt.Sprintf("o%d", i),
			Ra:  float64(i),
			Dec: float64(i),
			Mag: float64(i * 2),
		}
		if err = pw.Write(obj); err != nil {
			t.Log("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		t.Fatal("WriteStop error", err)
	}
	t.Log("Write Finished")
	fw.Close()

	return parquetFile
}

func TestReadParquet_read_all_file(t *testing.T) {
	filePath := Write(t, 10)
	source := source.Source{
		Reader:      []source.SourceReader{{Url: filePath}},
		CatalogName: "test",
	}

	outputs := make([]chan indexer.ReaderResult, 1)
	for i := range outputs {
		outputs[i] = make(chan indexer.ReaderResult)
	}
	parquetReader, err := NewParquetReader[Object](&source, outputs)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := parquetReader.Read()
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 10, len(rows))
	expectedOids := []string{"o0", "o1", "o2", "o3", "o4", "o5", "o6", "o7", "o8", "o9"}
	receivedOids := make([]string, 10, 10)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat().ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_single_file(t *testing.T) {
	filePath := Write(t, 10)
	source := source.Source{
		Reader:      []source.SourceReader{{Url: filePath}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
		CatalogName: "test",
	}

	outputs := make([]chan indexer.ReaderResult, 1)
	for i := range outputs {
		outputs[i] = make(chan indexer.ReaderResult)
	}
	parquetReader, err := NewParquetReader(&source, outputs, WithParquetBatchSize[Object](2))
	if err != nil {
		t.Fatal(err)
	}

	expectedOids := []string{"o0", "o1", "o2", "o3", "o4", "o5", "o6", "o7", "o8", "o9"}
	receivedOids := []string{}
	batches := 0
	var readErr error
	var rows []repository.InputSchema
	for {
		rows, readErr = parquetReader.ReadBatch()
		batches += 1
		if readErr != nil {
			break
		}

		for _, row := range rows {
			receivedOids = append(receivedOids, *row.ToMastercat().ID)
		}
	}
	require.Equal(t, 6, batches) // reader reads one extra batch with zero value
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_single_file_with_empty_batches(t *testing.T) {
	filePath := Write(t, 7)
	source := source.Source{
		Reader:      []source.SourceReader{{Url: filePath}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
		CatalogName: "test",
	}

	outputs := make([]chan indexer.ReaderResult, 1)
	for i := range outputs {
		outputs[i] = make(chan indexer.ReaderResult)
	}
	parquetReader, err := NewParquetReader(&source, outputs, WithParquetBatchSize[Object](2))
	if err != nil {
		t.Fatal(err)
	}

	expectedOids := []string{"o0", "o1", "o2", "o3", "o4", "o5", "o6"}
	receivedOids := []string{}
	batches := 0
	var readErr error
	var rows []repository.InputSchema
	for {
		rows, readErr = parquetReader.ReadBatch()
		batches += 1
		if readErr != nil {
			break
		}

		for _, row := range rows {
			receivedOids = append(receivedOids, *row.ToMastercat().ID)
		}
	}
	require.Equal(t, 5, batches) // reader reads one extra batch with zero value
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_larger_than_rows(t *testing.T) {
	filePath := Write(t, 2)
	source := source.Source{
		Reader:      []source.SourceReader{{Url: filePath}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
		CatalogName: "test",
	}

	outputs := make([]chan indexer.ReaderResult, 1)
	for i := range outputs {
		outputs[i] = make(chan indexer.ReaderResult)
	}
	parquetReader, err := NewParquetReader(&source, outputs, WithParquetBatchSize[Object](5))
	if err != nil {
		t.Fatal(err)
	}

	expectedOids := []string{"o0", "o1"}
	receivedOids := []string{}
	batches := 0
	var readErr error
	var rows []repository.InputSchema
	for {
		rows, readErr = parquetReader.ReadBatch()
		batches += 1
		if readErr != nil {
			break
		}

		for _, row := range rows {
			receivedOids = append(receivedOids, *row.ToMastercat().ID)
		}
	}
	require.Equal(t, 2, batches)
	require.Equal(t, expectedOids, receivedOids)
}
