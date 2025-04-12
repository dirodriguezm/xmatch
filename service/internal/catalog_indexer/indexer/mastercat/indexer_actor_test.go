package mastercat_indexer

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type TestSchema struct {
	Ra  float64
	Dec float64
	ID  string
	Cat string
}

// implement the interface
func (t *TestSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	return repository.ParquetMastercat{
		ID:   &t.ID,
		Ra:   &t.Ra,
		Dec:  &t.Dec,
		Cat:  &t.Cat,
		Ipix: &ipix,
	}
}

// implement the interface
func (t *TestSchema) ToMetadata() any {
	return t
}

func (t *TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

// implement the interface
func (t *TestSchema) SetField(name string, val interface{}) {}

func TestIndexActor(t *testing.T) {
	inbox := make(chan reader.ReaderResult)
	outbox := make(chan writer.WriterInput[any])
	rows := make([]repository.InputSchema, 2)
	rows[0] = &TestSchema{Ra: 0.0, Dec: 0.0, ID: "id1", Cat: "test"}
	rows[1] = &TestSchema{Ra: 0.0, Dec: 0.0, ID: "id2", Cat: "test"}

	src, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:",
		Type:        "csv",
		CatalogName: "catalog",
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
	})
	require.NoError(t, err)
	indexerActor, err := New(src, inbox, outbox, &config.IndexerConfig{OrderingScheme: "nested", Nside: 18})
	require.NoError(t, err)

	indexerActor.Start()
	inbox <- reader.ReaderResult{Rows: rows, Error: nil}
	close(inbox)
	results := make([]repository.ParquetMastercat, 2)
	for msg := range outbox {
		for i, obj := range msg.Rows {
			results[i] = obj.(repository.ParquetMastercat)
		}
	}

	require.Len(t, results, 2)
}
