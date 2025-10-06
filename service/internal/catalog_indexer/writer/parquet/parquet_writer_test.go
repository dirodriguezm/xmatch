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

package parquet_writer

import (
	"path"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/stretchr/testify/require"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func TestWrite(t *testing.T) {
	builder := AWriter[TestStruct](t)
	dir := t.TempDir()
	outputFile := path.Join(dir, "output.parquet")
	builder = builder.WithOutputFile(outputFile)
	w := builder.Build()
	rows := []any{TestStruct{"oid1", 1, 1}, TestStruct{"oid2", 2, 2}}

	w.Write(nil, actor.Message{Error: nil, Rows: rows})
	err := w.parquetWriter.WriteStop()
	require.NoError(t, err, "can't stop writer")
	w.pfile.Close()

	require.FileExists(t, outputFile)

	readRows := read_helper[TestStruct](t, outputFile)
	require.Len(t, readRows, 2)
	for i := range readRows {
		require.Equal(t, rows[i].(TestStruct).Oid, readRows[i].Oid)
		require.Equal(t, rows[i].(TestStruct).Ra, readRows[i].Ra)
		require.Equal(t, rows[i].(TestStruct).Dec, readRows[i].Dec)
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
	fr.Close()

	return rows
}
