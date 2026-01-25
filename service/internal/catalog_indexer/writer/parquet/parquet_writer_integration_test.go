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

package parquet_writer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/app"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var cfg config.Config

func TestMain(m *testing.M) {
	// create a config file
	tmpDir, err := os.MkdirTemp("", "parquet_writer_integration_test_*")
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
catalog_indexer:
  source:
    url: "buffer:"
    type: "csv"
    catalog_name: "vlass"
  reader:
    batch_size: 500
    type: "csv"
  database:
    url: "file:"
  indexer:
    ordering_scheme: "nested"
  indexer_writer:
    type: "parquet"
    output_file: "%s"
    schema: 1
`
	outputFile := filepath.Join(tmpDir, "test.parquet")
	config = fmt.Sprintf(config, outputFile)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		panic(err)
	}

	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}

	cfg, err = app.Config(getenv)
	if err != nil {
		panic(err)
	}

	m.Run()
	os.Remove(configPath)
	os.Remove(outputFile)
	os.Remove(tmpDir)
}

func TestWrite(t *testing.T) {
	writer, err := parquet_writer.New[repository.Mastercat](cfg.CatalogIndexer.IndexerWriter, t.Context())
	require.NoError(t, err)

	writer.Write(nil, actor.Message{
		Rows: []any{
			repository.Mastercat{ID: "oid1", Ipix: 1, Ra: 1, Dec: 1, Cat: "vlass"},
			repository.Mastercat{ID: "oid2", Ipix: 2, Ra: 2, Dec: 2, Cat: "vlass"},
		},
		Error: nil,
	})
	writer.Stop(nil)

	result := read_helper[repository.Mastercat](t, cfg.CatalogIndexer.IndexerWriter.OutputFile)

	require.Equal(t, 2, len(result))
	for i, row := range result {
		require.Equal(t, fmt.Sprintf("oid%d", i+1), row.ID)
		require.Equal(t, int64(i+1), row.Ipix)
		require.Equal(t, float64(i+1), row.Ra)
		require.Equal(t, float64(i+1), row.Dec)
		require.Equal(t, "vlass", row.Cat)
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
	err = fr.Close()
	require.NoError(t, err, "could not close file reader")

	return rows
}
