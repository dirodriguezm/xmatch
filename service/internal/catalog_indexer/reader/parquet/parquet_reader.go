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
	"log/slog"
	"reflect"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/xitongsys/parquet-go-source/local"
	preader "github.com/xitongsys/parquet-go/reader"
	psource "github.com/xitongsys/parquet-go/source"
)

type ParquetReader[T any] struct {
	*reader.BaseReader
	parquetReaders []*preader.ParquetReader
	src            *source.Source
	currentReader  int
	outbox         []chan reader.ReaderResult

	fileReaders []psource.ParquetFile
}

func NewParquetReader[T any](
	src *source.Source,
	channel []chan reader.ReaderResult,
	opts ...ParquetReaderOption[T],
) (*ParquetReader[T], error) {
	readers := []*preader.ParquetReader{}
	fileReaders := []psource.ParquetFile{}

	for _, srcReader := range src.Reader {
		fr, err := local.NewLocalFileReader(srcReader.Url)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewLocalFileReader\n%w", err)
		}
		fileReaders = append(fileReaders, fr)
		schema := new(T)
		pr, err := preader.NewParquetReader(fr, schema, 4)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
		}
		readers = append(readers, pr)
	}

	newReader := &ParquetReader[T]{
		BaseReader: &reader.BaseReader{
			Src:    src,
			Outbox: channel,
		},
		src:            src,
		currentReader:  0,
		outbox:         channel,
		fileReaders:    fileReaders,
		parquetReaders: readers,
	}

	for _, opt := range opts {
		opt(newReader)
	}

	newReader.Reader = newReader
	return newReader, nil
}

func (r *ParquetReader[T]) ReadSingleFile(currentReader *preader.ParquetReader) ([]repository.InputSchema, error) {
	defer currentReader.ReadStop()

	nrows := currentReader.GetNumRows()
	records := make([]T, nrows)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(r.currentReader, err, r.src, "Failed to read parquet")
	}

	parsedRecords := convertToInputSchema(records, r.src.CatalogName)

	return parsedRecords, nil
}

func (r *ParquetReader[T]) Read() ([]repository.InputSchema, error) {
	defer func() {
		for i := range r.fileReaders {
			r.fileReaders[i].Close()
		}
	}()

	rows := make([]repository.InputSchema, 0, 0)
	for _, currentReader := range r.parquetReaders {
		currentRows, err := r.ReadSingleFile(currentReader)
		if err != nil {
			return nil, err
		}
		rows = append(rows, currentRows...)
	}
	return rows, nil
}

func (r *ParquetReader[T]) ReadBatchSingleFile(currentReader *preader.ParquetReader) ([]repository.InputSchema, error) {
	records := make([]T, r.BatchSize)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(r.currentReader, err, r.src, "Failed to read parquet in batch")
	}

	parsedRecords := convertToInputSchema(records, r.src.CatalogName)

	if isZeroValueSlice(parsedRecords) {
		// finished reading
		currentReader.ReadStop()
		return nil, io.EOF
	}

	return parsedRecords, nil
}

func (r *ParquetReader[T]) ReadBatch() ([]repository.InputSchema, error) {
	currentReader := r.parquetReaders[r.currentReader]
	slog.Debug("ParquetReader reading batch")
	currentRows, err := r.ReadBatchSingleFile(currentReader)
	if err != nil {
		if err == io.EOF && r.currentReader < len(r.parquetReaders)-1 {
			// only finished current file, but more files remain
			r.currentReader += 1
			return currentRows, nil
		}
		if err == io.EOF {
			// finished reading all files
			r.currentReader += 1
			for i := range r.fileReaders {
				r.fileReaders[i].Close()
			}
			return currentRows, err
		}
		return nil, err
	}
	return currentRows, nil
}

func convertToInputSchema[T any](records []T, catalogName string) []repository.InputSchema {
	inputSchemas := make([]repository.InputSchema, len(records))
	for i := 0; i < len(records); i++ {
		elem := reflect.ValueOf(records[i])
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		// Check if it's a struct
		if elem.Kind() != reflect.Struct {
			panic(fmt.Errorf("expected struct, got %v", elem.Kind()))
		}

		switch catalogName {
		case "allwise":
			inputSchemas[i] = &repository.AllwiseInputSchema{}
			inputSchemas[i].SetField("Source_id", elem.FieldByName("Source_id").Interface())
			inputSchemas[i].SetField("Ra", elem.FieldByName("Ra").Interface())
			inputSchemas[i].SetField("Dec", elem.FieldByName("Dec").Interface())
			inputSchemas[i].SetField("W1mpro", elem.FieldByName("W1mpro").Interface())
			inputSchemas[i].SetField("W1sigmpro", elem.FieldByName("W1sigmpro").Interface())
			inputSchemas[i].SetField("W2mpro", elem.FieldByName("W2mpro").Interface())
			inputSchemas[i].SetField("W2sigmpro", elem.FieldByName("W2sigmpro").Interface())
			inputSchemas[i].SetField("W3mpro", elem.FieldByName("W3mpro").Interface())
			inputSchemas[i].SetField("W3sigmpro", elem.FieldByName("W3sigmpro").Interface())
			inputSchemas[i].SetField("W4mpro", elem.FieldByName("W4mpro").Interface())
			inputSchemas[i].SetField("W4sigmpro", elem.FieldByName("W4sigmpro").Interface())
			inputSchemas[i].SetField("J_m_2mass", elem.FieldByName("J_m_2mass").Interface())
			inputSchemas[i].SetField("H_m_2mass", elem.FieldByName("H_m_2mass").Interface())
			inputSchemas[i].SetField("K_m_2mass", elem.FieldByName("K_m_2mass").Interface())
			inputSchemas[i].SetField("J_msig_2mass", elem.FieldByName("J_msig_2mass").Interface())
			inputSchemas[i].SetField("H_msig_2mass", elem.FieldByName("H_msig_2mass").Interface())
			inputSchemas[i].SetField("K_msig_2mass", elem.FieldByName("K_msig_2mass").Interface())
		default:
			inputSchemas[i] = &TestInputSchema{}
			inputSchemas[i].SetField("Oid", elem.FieldByName("Oid").Interface())
			inputSchemas[i].SetField("Ra", elem.FieldByName("Ra").Interface())
			inputSchemas[i].SetField("Dec", elem.FieldByName("Dec").Interface())
		}

	}
	return inputSchemas
}

func isZeroValueSlice(s []repository.InputSchema) bool {
	for i := 0; i < len(s); i++ {
		if !isZeroValueInputSchema(s[i]) {
			return false
		}
	}
	return true
}

func isZeroValueInputSchema(schema repository.InputSchema) bool {
	elem := reflect.ValueOf(schema)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	// Check if it's a struct
	if elem.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected struct, got %v", elem.Kind()))
	}

	return isZeroValue(elem.Interface())
}

// Helper function to check if an interface{} value is zero
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch value := v.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(value).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(value).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(value).Float() == 0
	case bool:
		return !value
	case string:
		return value == ""
	case []interface{}:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	default:
		// For struct types, compare with their zero value
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}
