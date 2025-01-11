package parquet_writer

import (
	"path"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func TestConvertMapToStruct(t *testing.T) {
	type TestStruct struct {
		Oid string
		Ra  float64
		Dec float64
	}

	data := indexer.Row{
		"ra":  float64(1),
		"dec": float64(2),
		"oid": "oid1",
	}
	result := convertMapToStruct[TestStruct](data)
	require.Equal(t, data["oid"], result.Oid)
	require.Equal(t, data["ra"], result.Ra)
	require.Equal(t, data["dec"], result.Dec)
}

func TestConvertMapToStruct_SnakeCase(t *testing.T) {
	type TestStruct struct {
		ObjectIdLong string
		Ra           float64
		Dec          float64
	}

	data := indexer.Row{
		"ra":             float64(1),
		"dec":            float64(2),
		"object_id_long": "oid1",
	}
	result := convertMapToStruct[TestStruct](data)
	require.Equal(t, data["object_id_long"], result.ObjectIdLong)
	require.Equal(t, data["ra"], result.Ra)
	require.Equal(t, data["dec"], result.Dec)
}

func TestReceive(t *testing.T) {
	type TestStruct struct {
		Oid string  `parquet:"name=oid, type=BYTE_ARRAY"`
		Ra  float64 `parquet:"name=ra, type=DOUBLE"`
		Dec float64 `parquet:"name=dec, type=DOUBLE"`
	}
	builder := AWriter[TestStruct](t)
	dir := t.TempDir()
	outputFile := path.Join(dir, "output.parquet")
	builder = builder.WithOutputFile(outputFile)
	w := builder.Build()
	rows := []indexer.Row{
		{
			"oid": "oid1",
			"ra":  float64(1),
			"dec": float64(1),
		},
		{
			"oid": "oid2",
			"ra":  float64(2),
			"dec": float64(2),
		},
	}

	w.Receive(indexer.WriterInput{Error: nil, Rows: rows})
	err := w.parquetWriter.WriteStop()
	require.NoError(t, err, "can't stop writer")
	w.pfile.Close()

	require.FileExists(t, outputFile)

	readRows := read_helper[TestStruct](t, outputFile)
	require.Len(t, readRows, 2)
	for i := range readRows {
		require.Equal(t, rows[i]["oid"], readRows[i].Oid)
		require.Equal(t, rows[i]["ra"], readRows[i].Ra)
		require.Equal(t, rows[i]["dec"], readRows[i].Dec)
	}
}

func TestStart(t *testing.T) {
	type TestStruct struct {
		Oid string  `parquet:"name=oid, type=BYTE_ARRAY"`
		Ra  float64 `parquet:"name=ra, type=DOUBLE"`
		Dec float64 `parquet:"name=dec, type=DOUBLE"`
	}
	builder := AWriter[TestStruct](t)
	file := path.Join(t.TempDir(), "output.parquet")
	builder = builder.WithOutputFile(file)
	w := builder.Build()

	msg := indexer.WriterInput{
		Error: nil,
		Rows: []indexer.Row{
			{
				"oid": "oid1",
				"ra":  float64(1),
				"dec": float64(1),
			},
			{
				"oid": "oid2",
				"ra":  float64(2),
				"dec": float64(2),
			},
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
		require.Equal(t, msg.Rows[i]["oid"], readRows[i].Oid)
		require.Equal(t, msg.Rows[i]["ra"], readRows[i].Ra)
		require.Equal(t, msg.Rows[i]["dec"], readRows[i].Dec)
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
