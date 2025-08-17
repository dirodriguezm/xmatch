package partition_reader

import (
	"os"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
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

func (r TestInputSchema) FillMetadata(dst repository.Metadata) {}

func (r TestInputSchema) FillMastercat(dst *repository.Mastercat, ipix int64) {}

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

func stringPtr(s string) *string {
	return &s
}
