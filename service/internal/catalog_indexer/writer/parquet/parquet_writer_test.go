package parquet_writer

import (
	"path"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func TestReceive(t *testing.T) {
	builder := AWriter[TestStruct](t)
	dir := t.TempDir()
	outputFile := path.Join(dir, "output.parquet")
	builder = builder.WithOutputFile(outputFile)
	w := builder.Build()
	rows := []TestStruct{{"oid1", 1, 1}, {"oid2", 2, 2}}

	w.Receive(writer.WriterInput[TestStruct]{Error: nil, Rows: rows})
	err := w.parquetWriter.WriteStop()
	require.NoError(t, err, "can't stop writer")
	w.pfile.Close()

	require.FileExists(t, outputFile)

	readRows := read_helper[TestStruct](t, outputFile)
	require.Len(t, readRows, 2)
	for i := range readRows {
		require.Equal(t, rows[i].Oid, readRows[i].Oid)
		require.Equal(t, rows[i].Ra, readRows[i].Ra)
		require.Equal(t, rows[i].Dec, readRows[i].Dec)
	}
}

func TestStart(t *testing.T) {
	builder := AWriter[TestStruct](t)
	file := path.Join(t.TempDir(), "output.parquet")
	builder = builder.WithOutputFile(file)
	w := builder.Build()

	msg := writer.WriterInput[TestStruct]{
		Error: nil,
		Rows: []TestStruct{
			{"oid1", 1, 1},
			{"oid2", 2, 2},
		},
	}

	w.Start()
	builder.input <- msg
	close(builder.input)
	w.Done()

	require.FileExists(t, file)

	readRows := read_helper[TestStruct](t, file)
	require.Len(t, readRows, 2)
	for i := range readRows {
		require.Equal(t, msg.Rows[i].Oid, readRows[i].Oid)
		require.Equal(t, msg.Rows[i].Ra, readRows[i].Ra)
		require.Equal(t, msg.Rows[i].Dec, readRows[i].Dec)
	}
}

func read_helper[T any](t *testing.T, file string) []T {
	t.Helper()

	fr, err := local.NewLocalFileReader(file)
	require.NoError(t, err, "could not create local file reader")

	pr, err := reader.NewParquetReader(fr, new(T), 4)
	require.NoError(t, err, "could not create parquet reader")

	num := int(pr.GetNumRows())

	rows := make([]T, num)
	err = pr.Read(&rows)
	require.NoError(t, err, "could not read rows")

	pr.ReadStop()
	fr.Close()

	return rows
}
