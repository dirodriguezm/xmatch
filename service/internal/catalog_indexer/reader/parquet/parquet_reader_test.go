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

package parquet_reader

import (
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
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

func (o Object) GetId() string {
	return o.Oid
}

func (o Object) GetCoordinates() (float64, float64) {
	return o.Ra, o.Dec
}

func (o Object) FillMastercat(dst *repository.Mastercat, ipix int64) {
	dst.ID = o.Oid
	dst.Ra = o.Ra
	dst.Dec = o.Dec
}
func (o Object) FillMetadata(dst repository.Metadata) {}

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
	for i := range nrows {
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
	source, err := source.NewSource(&config.SourceConfig{
		Url:         "file:" + filePath,
		Type:        "parquet",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	parquetReader, err := NewParquetReader(source, WithParquetBatchSize[Object](2))
	if err != nil {
		t.Fatal(err)
	}

	rows, err := parquetReader.Read()
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	require.Equal(t, 10, len(rows))
	expectedOids := []string{"o0", "o1", "o2", "o3", "o4", "o5", "o6", "o7", "o8", "o9"}
	receivedOids := make([]string, 10)
	for i, row := range rows {
		fmt.Println(row)
		mastercat := repository.Mastercat{}
		row.FillMastercat(&mastercat, 0)
		receivedOids[i] = mastercat.ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_single_file(t *testing.T) {
	filePath := Write(t, 10)
	source, err := source.NewSource(&config.SourceConfig{
		Url:         "file:" + filePath,
		Type:        "parquet",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	parquetReader, err := NewParquetReader(source, WithParquetBatchSize[Object](2))
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
			mastercat := repository.Mastercat{}
			row.FillMastercat(&mastercat, 0)
			receivedOids = append(receivedOids, mastercat.ID)
		}
	}
	require.Equal(t, 6, batches) // reader reads one extra batch with zero value
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_single_file_with_empty_batches(t *testing.T) {
	filePath := Write(t, 7)
	source, err := source.NewSource(&config.SourceConfig{
		Url:         "file:" + filePath,
		Type:        "parquet",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	parquetReader, err := NewParquetReader(source, WithParquetBatchSize[Object](2))
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
			mastercat := repository.Mastercat{}
			row.FillMastercat(&mastercat, 0)
			receivedOids = append(receivedOids, mastercat.ID)
		}
	}
	require.Equal(t, 5, batches) // reader reads one extra batch with zero value
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadParquet_read_batch_larger_than_rows(t *testing.T) {
	filePath := Write(t, 2)
	source, err := source.NewSource(&config.SourceConfig{
		Url:         "file:" + filePath,
		Type:        "parquet",
		CatalogName: "test",
		Nside:       18,
		Metadata:    false,
	})
	require.NoError(t, err)

	parquetReader, err := NewParquetReader(source, WithParquetBatchSize[Object](2))
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
			mastercat := repository.Mastercat{}
			row.FillMastercat(&mastercat, 0)
			receivedOids = append(receivedOids, mastercat.ID)
		}
	}
	require.Equal(t, 2, batches)
	require.Equal(t, expectedOids, receivedOids)
}
