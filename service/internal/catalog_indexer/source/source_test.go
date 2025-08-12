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

	require.Len(t, source.Sources, 2)
}

func TestSourceReader_NestedFiles(t *testing.T) {
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	source := ASource(t).WithUrl(url).WithNestedCsvFiles([]string{""}, []string{""}).Build()

	require.Len(t, source.Sources, 2)
}

func TestSourceReader_ParquetFiles(t *testing.T) {
	dir := t.TempDir()
	url := fmt.Sprintf("files:%s", dir)
	metadata := []string{"name=Col, type=INT64"}
	source := ASource(t).WithUrl(url).WithParquetFiles(metadata, [][][]string{{{"1"}}, {{"2"}}}).Build()

	require.Len(t, source.Sources, 2)
}
