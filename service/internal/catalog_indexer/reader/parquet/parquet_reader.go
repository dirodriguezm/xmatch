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
	"os"
	"reflect"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/xitongsys/parquet-go-source/local"
	preader "github.com/xitongsys/parquet-go/reader"
	psource "github.com/xitongsys/parquet-go/source"
)

type ParquetReader[T repository.InputSchema] struct {
	currentParquetReader *preader.ParquetReader
	currentFileReader    psource.ParquetFile
	currentFileName      string
	src                  *source.Source
	batchSize            int
}

func NewParquetReader[T repository.InputSchema](src *source.Source, opts ...ParquetReaderOption[T]) (*ParquetReader[T], error) {
	currentReader, err := src.Next()
	if err != nil {
		return nil, fmt.Errorf("could not get next source: %w", err)
	}
	currentFileName := currentReader.(*os.File).Name()
	currentReader.(*os.File).Close()

	fr, err := local.NewLocalFileReader(currentFileName)
	if err != nil {
		return nil, fmt.Errorf("Could not create NewLocalFileReader\n%w", err)
	}

	schema := new(T)
	pr, err := preader.NewParquetReader(fr, schema, 1)
	if err != nil {
		return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
	}

	newReader := &ParquetReader[T]{
		currentParquetReader: pr,
		currentFileReader:    fr,
		currentFileName:      currentFileName,
		src:                  src,
	}

	for _, opt := range opts {
		opt(newReader)
	}

	return newReader, nil
}

func (r *ParquetReader[T]) ReadSingleFile(src *source.Source, currentReader *preader.ParquetReader) ([]repository.InputSchema, error) {
	defer currentReader.ReadStop()

	nrows := currentReader.GetNumRows()
	records := make([]T, nrows)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(err, src, "Failed to read parquet")
	}

	return convertToInputSchema(records), nil
}

func convertToInputSchema[T repository.InputSchema](records []T) []repository.InputSchema {
	converted := make([]repository.InputSchema, len(records))
	for i := range records {
		converted[i] = records[i]
	}
	return converted
}

func (r *ParquetReader[T]) Read() ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.currentParquetReader.GetNumRows())
	eof := false

	for !eof {
		// Read the current file completely
		currentRows, err := r.ReadSingleFile(r.src, r.currentParquetReader)
		if err != nil {
			return nil, fmt.Errorf("Could not read file: %s. %w", r.currentFileName, err)
		}
		rows = append(rows, currentRows...)

		// Get next file if any
		nextReader, err := r.src.Next()
		if err == io.EOF {
			//no more files
			return rows, io.EOF
		}
		if err != nil {
			return nil, fmt.Errorf("Could not get next source: %w", err)
		}

		nextFileName := nextReader.(*os.File).Name()
		nextReader.(*os.File).Close()
		// If no more files remaining, the returned error will be EOF
		eof = err == io.EOF

		// Create a new local file reader from the next file in the Source
		newFileReader, err := local.NewLocalFileReader(nextFileName)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewLocalFileReader\n%w", err)
		}

		// Create a new Parquet File Reader
		schema := new(T)
		newParquetReader, err := preader.NewParquetReader(newFileReader, schema, 1)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
		}

		// Update the current file reader and close the previous
		r.currentParquetReader.ReadStop()
		r.currentFileReader.Close()
		r.currentFileReader = newFileReader
		r.currentParquetReader = newParquetReader
		r.currentFileName = nextFileName
	}

	return rows, nil
}

// Reads batch from the passed reader
// The closing should be handling by the caller
// Returns io.EOF if there are no more rows to read
func (r *ParquetReader[T]) ReadBatchSingleFile(currentReader *preader.ParquetReader) ([]repository.InputSchema, error) {
	records := make([]T, r.batchSize)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(err, r.src, "Failed to read parquet in batch")
	}

	if isZeroValueSlice(records) {
		// finished reading
		return nil, io.EOF
	}

	return convertToInputSchema(records), nil
}

func (r *ParquetReader[T]) ReadBatch() ([]repository.InputSchema, error) {
	currentRows, err := r.ReadBatchSingleFile(r.currentParquetReader)

	// Read did not finish successfully
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Error while reading batch from current parquet reader: %s. %w", r.currentFileName, err)
	}

	// EOF means we only finished reading the current file, but more files could remain
	if err == io.EOF {
		// Get the next file in the source
		nextReader, nextError := r.src.Next()
		// If error is not EOF it means there has been an error retrieving the next file from the Source
		if nextError != nil && nextError != io.EOF {
			return currentRows, fmt.Errorf("Error getting the next file from the Source: %w", nextError)
		}
		// If error is EOF, it means that there are no more files to read
		if nextError == io.EOF {
			return currentRows, io.EOF
		}

		// Now we update the current reader and close the previous
		//
		// first get the next file name
		nextFileName := nextReader.(*os.File).Name()
		nextReader.(*os.File).Close()

		// Create a new local file reader from the next file in the Source
		newFileReader, err := local.NewLocalFileReader(nextFileName)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewLocalFileReader\n%w", err)
		}

		// Create a new Parquet File Reader
		schema := new(T)
		newParquetReader, err := preader.NewParquetReader(newFileReader, schema, 1)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
		}

		// update the current values
		r.currentParquetReader.ReadStop()
		r.currentFileReader.Close()
		r.currentFileReader = newFileReader
		r.currentParquetReader = newParquetReader
		r.currentFileName = nextFileName

		// return the latest rows
		return currentRows, nil
	}

	// Here we are still reading the current file.
	return currentRows, nil
}

func isZeroValueSlice[T repository.InputSchema](s []T) bool {
	for i := range s {
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
func isZeroValue(v any) bool {
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
	case []any:
		return len(value) == 0
	case map[string]any:
		return len(value) == 0
	default:
		// For struct types, compare with their zero value
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}

func (r *ParquetReader[T]) Close() error {
	r.currentParquetReader.ReadStop()
	return r.currentFileReader.Close()
}
