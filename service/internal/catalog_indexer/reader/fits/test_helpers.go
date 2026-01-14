// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fits_reader

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"codeberg.org/astrogo/fitsio"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func writeFitsFile(t *testing.T, filename string, rows []map[string]any) {
	t.Helper()

	t.Log("Creating file", filename)
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	t.Log("Created file", filename)

	t.Log("Creating fits file")
	fitsFile, err := fitsio.Create(f)
	if err != nil {
		t.Fatal("Could not create fits file")
	}
	defer fitsFile.Close()
	t.Log("Created fits file")

	t.Log("Creating primaryHDU")
	phdu, err := fitsio.NewPrimaryHDU(nil)
	if err != nil {
		t.Fatal("Could not create primaryHDU %w", err)
	}
	t.Log("Created primaryHDU")

	t.Log("Writing primaryHDU to file")
	if err := fitsFile.Write(phdu); err != nil {
		t.Fatal("Could not write primaryHDU to file %w", err)
	}
	t.Log("Wrote primaryHDU to file")

	// create table
	t.Log("Creating table")
	fitsTable, err := CreateFitsTable(rows)
	if err != nil {
		t.Fatal("Could not create table %w", err)
	}
	defer fitsTable.Close()
	t.Log("Created table")

	// Write the table to the file
	t.Log("Writing table to file")
	err = fitsFile.Write(fitsTable)
	if err != nil {
		t.Fatal("Could not write table to file %w", err)
	}
	t.Log("Wrote table to file")

}

func CreateFitsTable(data []map[string]any) (*fitsio.Table, error) {
	// create columns using the keys of the first map
	// the type of the column is determined by the type of the value
	columns, err := createColumns(data)
	if err != nil {
		return nil, err
	}

	data = cleanData(data, columns)

	// create table from columns
	table, err := fitsio.NewTable("results", columns, fitsio.BINARY_TBL)
	if err != nil {
		return nil, err
	}

	// populate the table
	rslice := reflect.ValueOf(data)
	for i := range rslice.Len() {
		row := rslice.Index(i).Addr()
		err := table.Write(row.Interface())
		if err != nil {
			return nil, err
		}
	}

	nrows := table.NumRows()
	if nrows != int64(len(data)) {
		return nil, fmt.Errorf("Error creating table: number of rows written (%d) does not match number of rows in data (%d)", nrows, len(data))
	}

	return table, nil
}

func createColumns(data []map[string]any) ([]fitsio.Column, error) {
	if len(data) == 0 {
		return []fitsio.Column{}, nil
	}

	columns := make([]fitsio.Column, 0, len(data[0]))
	keys := make([]string, 0, len(data[0]))
	for key := range data[0] {
		keys = append(keys, key)
	}
	// sort the keys to ensure consistent column order
	sort.Strings(keys)
	var format string
	for _, key := range keys {
		for i := range data {
			if data[i][key] == nil {
				continue
			}
			value := data[i][key]
			format = getFormat(value)
			if format == "" {
				return nil, fmt.Errorf("Error creating columns: unsupported type %T for key %s", value, key)
			}
			columns = append(columns, fitsio.Column{
				Name:   key,
				Format: format,
			})
			break
		}
	}
	return columns, nil
}

func getFormat(value any) string {
	switch value := value.(type) {
	case int16:
		return "I" // 16-bit integer
	case int32:
		return "J" // 32-bit integer
	case int:
		return "K" // 64-bit integer
	case int64:
		return "K" // 64-bit integer
	case float32:
		return "E" // 32-bit floating point
	case float64:
		return "D" // 64-bit floating point
	case string:
		stringLength := len(value)
		return fmt.Sprintf("%dA", stringLength) // Character string with length
	case bool:
		return "L" // logical
	default:
		return ""
	}
}

func cleanData(data []map[string]any, columns []fitsio.Column) []map[string]any {
	if len(data) == 0 {
		return []map[string]any{}
	}

	cleanedData := make([]map[string]any, len(data))
	for i := range data {
		cleanedData[i] = make(map[string]any, len(columns))
	}
	keys := make([]string, 0, len(data[0]))
	for key := range data[0] {
		keys = append(keys, key)
	}
	// sort the keys to ensure consistent column order
	sort.Strings(keys)
	for _, key := range keys {
		for i := range data {
			if data[i][key] != nil {
				cleanedData[i][key] = data[i][key]
				continue
			}
			for j := range columns {
				if columns[j].Name != key {
					continue
				}
				format := columns[j].Format
				cleanedData[i][key] = getZeroValueForFormat(format)
			}
		}
	}
	return cleanedData
}

func getZeroValueForFormat(format string) any {
	if strings.Contains(format, "A") {
		format = "A"
	}
	switch format {
	case "I":
		return int16(0)
	case "J":
		return int32(0)
	case "K":
		return int64(0)
	case "E":
		return float32(0)
	case "D":
		return float64(0)
	case "A":
		return ""
	case "L":
		return false
	default:
		return nil
	}
}

type TestSchema struct {
	Oid string
	Ra  float64
	Dec float64
}

func (t TestSchema) FillMastercat(dst *repository.Mastercat, ipix int64) {
	dst.ID = t.Oid
	dst.Ra = t.Ra
	dst.Dec = t.Dec
	dst.Cat = "test"
	dst.Ipix = ipix
}

func (t TestSchema) FillMetadata(dst repository.Metadata) {}

func (t TestSchema) GetCoordinates() (float64, float64) {
	return t.Ra, t.Dec
}

func (t TestSchema) GetId() string {
	return t.Oid
}
