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

package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

type SourceBuilder struct {
	SourceConfig *config.SourceConfig
	t            *testing.T
}

func ASource(t *testing.T) *SourceBuilder {
	t.Helper()

	return &SourceBuilder{
		t: t,
		SourceConfig: &config.SourceConfig{
			Type:        "csv",
			CatalogName: "test",
			RaCol:       "ra",
			DecCol:      "dec",
			OidCol:      "oid",
			Nside:       18,
		},
	}
}

func (builder *SourceBuilder) WithType(t string) *SourceBuilder {
	builder.t.Helper()

	builder.SourceConfig = &config.SourceConfig{
		Type:        t,
		CatalogName: "test",
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		Nside:       18,
	}
	return builder
}

func (builder *SourceBuilder) Build() *Source {
	builder.t.Helper()

	src, err := NewSource(builder.SourceConfig)
	require.NoError(builder.t, err)
	return src
}

func (builder *SourceBuilder) WithUrl(url string) *SourceBuilder {
	builder.t.Helper()

	builder.SourceConfig.Url = url
	return builder
}

func (builder *SourceBuilder) WithCsvFiles(data []string) *SourceBuilder {
	builder.t.Helper()

	dir := strings.Split(builder.SourceConfig.Url, "files:")[1]
	for i, fileContent := range data {
		fileName := fmt.Sprintf("file%d.csv", i)
		filePath := filepath.Join(dir, fileName)
		err := os.WriteFile(filePath, []byte(fileContent), 0644)
		require.NoError(builder.t, err)
	}
	return builder
}

func (builder *SourceBuilder) WithNestedCsvFiles(args ...[]string) *SourceBuilder {
	builder.t.Helper()

	builder.WithEmptyDir()
	dir := strings.Split(builder.SourceConfig.Url, "files:")[1]
	for nestedLevel, files := range args {
		if nestedLevel > 0 {
			nestedDir := fmt.Sprintf("nested%d", nestedLevel)
			dir = filepath.Join(dir, nestedDir)
			err := os.Mkdir(dir, 0777)
			require.NoError(builder.t, err)
		}
		for i, fileContent := range files {
			fileName := fmt.Sprintf("file%d.csv", i)
			filePath := filepath.Join(dir, fileName)
			err := os.WriteFile(filePath, []byte(fileContent), 0644)
			require.NoError(builder.t, err)
		}
	}
	return builder
}

func (builder *SourceBuilder) WithParquetFiles(metadata []string, data [][][]string) *SourceBuilder {
	builder.t.Helper()

	dir := strings.Split(builder.SourceConfig.Url, "files:")[1]
	for i, fileData := range data {
		fname := fmt.Sprintf("file%d.parquet", i)
		fw, err := local.NewLocalFileWriter(filepath.Join(dir, fname))
		require.NoError(builder.t, err)

		pw, err := writer.NewCSVWriter(metadata, fw, 1)
		require.NoError(builder.t, err)

		for rowNum := range fileData {
			row := make([]*string, len(fileData[rowNum]))
			for j := range row {
				row[j] = &fileData[rowNum][j]
			}
			err = pw.WriteString(row)
			require.NoError(builder.t, err)
		}

		err = pw.WriteStop()
		require.NoError(builder.t, err)

		err = fw.Close()
		require.NoError(builder.t, err)

		// verify file
		require.FileExistsf(builder.t, filepath.Join(dir, fname), "Parquet file does not exist %s", filepath.Join(dir, fname))
		file, err := os.Open(filepath.Join(dir, fname))
		require.NoError(builder.t, err)
		b := make([]byte, 100)
		file.Read(b)
		require.NotEmpty(builder.t, b)
	}

	return builder
}

func (builder *SourceBuilder) WithEmptyDir() *SourceBuilder {
	builder.t.Helper()

	dir := strings.Split(builder.SourceConfig.Url, "files:")[1]
	err := cleanDir(dir)
	require.NoError(builder.t, err)
	return builder
}

func cleanDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return cleanDir(filepath.Join(dir, entry.Name()))
		}
		err := os.Remove(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
