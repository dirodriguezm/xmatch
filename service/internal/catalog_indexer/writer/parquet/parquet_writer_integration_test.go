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
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

var ctr container.Container

func TestMain(m *testing.M) {
	// create a config file
	tmpDir, err := os.MkdirTemp("", "parquet_writer_integration_test_*")
	if err != nil {
		slog.Error("could not make temp dir")
		panic(err)
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
catalog_indexer:
  source:
    url: "buffer:"
    type: "csv"
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
`
	outputFile := filepath.Join(tmpDir, "test.parquet")
	config = fmt.Sprintf(config, outputFile)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		slog.Error("could not write config file")
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
	ctx := context.Background()
	stdout := &strings.Builder{}

	// build DI container
	ctr = di.BuildIndexerContainer(ctx, getenv, stdout)

	m.Run()
	os.Remove(configPath)
	os.Remove(outputFile)
	os.Remove(tmpDir)
}

func TestActor(t *testing.T) {
	var w writer.Writer[any]
	err := ctr.NamedResolve(&w, "indexer_writer")
	require.NoError(t, err)

	w.Start()
	oids := []string{"oid1", "oid2"}
	ras := []float64{1, 2}
	decs := []float64{1, 2}
	ipixs := []int64{1, 2}
	cat := "vlass"
	w.(*parquet_writer.ParquetWriter[any]).InboxChannel <- writer.WriterInput[any]{
		Rows: []any{
			repository.ParquetMastercat{ID: &oids[0], Ipix: &ipixs[0], Ra: &ras[0], Dec: &decs[0], Cat: &cat},
			repository.ParquetMastercat{ID: &oids[1], Ipix: &ipixs[1], Ra: &ras[1], Dec: &decs[1], Cat: &cat},
		},
	}
	close(w.(*parquet_writer.ParquetWriter[any]).InboxChannel)
	w.Done()

	// check the parquet file
	var cfg *config.Config
	err = ctr.Resolve(&cfg)
	require.NoError(t, err)
	result := read_helper[repository.ParquetMastercat](t, cfg.CatalogIndexer.IndexerWriter.OutputFile)

	require.Equal(t, 2, len(result))
	for i, row := range result {
		require.Equal(t, fmt.Sprintf("oid%d", i+1), *row.ID)
		require.Equal(t, int64(i+1), *row.Ipix)
		require.Equal(t, float64(i+1), *row.Ra)
		require.Equal(t, float64(i+1), *row.Dec)
		require.Equal(t, "vlass", *row.Cat)
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
