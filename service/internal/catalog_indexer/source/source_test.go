package source

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSourceReader_File(t *testing.T) {

}

func TestSourceReader_Buffer(t *testing.T) {

}

func TestSourceReader_Files(t *testing.T) {
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	source := ASource(t).WithUrl(url).WithCsvFiles([]string{"", ""}).Build()

	require.Len(t, source.Reader, 2)
}

func TestSourceReader_NestedFiles(t *testing.T) {
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	source := ASource(t).WithUrl(url).WithNestedCsvFiles([]string{""}, []string{""}).Build()

	require.Len(t, source.Reader, 2)
}

func TestSourceReader_ParquetFiles(t *testing.T) {
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	metadata := []string{"name=Col, type=INT64"}
	source := ASource(t).WithUrl(url).WithParquetFiles(metadata, [][][]string{{{"1"}}, {{"2"}}}).Build()

	require.Len(t, source.Reader, 2)
}
