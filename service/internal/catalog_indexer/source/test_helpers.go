package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
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
			CatalogName: "tesc",
			RaCol:       "ra",
			DecCol:      "dec",
			OidCol:      "oid",
		},
	}
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
