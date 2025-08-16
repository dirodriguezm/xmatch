package partition_writer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
)

type TestInputSchema struct {
	Id      *string `parquet:"name=id, type=BYTE_ARRAY"`
	Column1 *int    `parquet:"name=column1, type=INT64"`
	Column2 *string `parquet:"name=column2, type=BYTE_ARRAY"`
}

func (r TestInputSchema) GetCoordinates() (float64, float64) {
	return 0, 0
}

func (r TestInputSchema) ToMetadata() any {
	return nil
}

func (r TestInputSchema) ToMastercat(ipix int64) repository.Mastercat {
	return repository.Mastercat{}
}

func (r TestInputSchema) GetId() string {
	return *r.Id
}

func writeParquet(t *testing.T, file *os.File, rows []TestInputSchema) {
	t.Helper()

	pr, err := writer.NewParquetWriterFromWriter(file, new(TestInputSchema), 1)
	if err != nil {
		t.Fatal(err)
	}

	for _, row := range rows {
		if err := pr.Write(row); err != nil {
			t.Fatal(err)
		}
	}

	if err := pr.WriteStop(); err != nil {
		t.Fatal(err)
	}
}

func readParquet(t *testing.T, file string) []TestInputSchema {
	t.Helper()

	fr, err := local.NewLocalFileReader(file)
	require.NoError(t, err)
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(TestInputSchema), 1)
	require.NoError(t, err)
	defer pr.ReadStop()

	nrows := pr.GetNumRows()

	records := make([]TestInputSchema, nrows)
	if err := pr.Read(&records); err != nil {
		t.Fatal(err)
	}

	return records
}

func CreateTestConfig(outputFile string) string {
	// create a config file
	tmpDir, err := os.MkdirTemp("", "partition_writer_integration_test_*")
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
reader:
  batch_size: 500
preprocessor:
  source:
    url: "buffer:"
    type: "csv"
    catalog_name: "test"
  reader:
    batch_size: 500
    type: "csv"
  partition_writer:
    max_file_size: 104857600
    num_partitions: 1
    partition_levels: 1
    base_dir: "%s"
  partition_reader:
    num_workers: 1
  reducer_writer:
    batch_size: 100
    type: "parquet"
    output_file: "%s"
`
	config = fmt.Sprintf(config, tmpDir, outputFile)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		panic(err)
	}
	return configPath
}
