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
	"os"
	"reflect"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"

	"github.com/xitongsys/parquet-go-source/local"
	preader "github.com/xitongsys/parquet-go/reader"
	psource "github.com/xitongsys/parquet-go/source"
)

type rawRecordFactory interface {
	NewRawRecord() any
}

type ParquetReader struct {
	currentParquetReader *preader.ParquetReader
	currentFileReader    psource.ParquetFile
	currentFileName      string
	src                  *source.Source
	adapter              rawRecordFactory
	recordType           reflect.Type
	batchSize            int
}

func NewParquetReader(src *source.Source, adapter rawRecordFactory, opts ...ParquetReaderOption) (*ParquetReader, error) {
	if adapter == nil {
		return nil, fmt.Errorf("parquet reader requires a raw record factory")
	}
	recordType := reflect.TypeOf(adapter.NewRawRecord())
	if recordType == nil || recordType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("raw record factory must return a struct")
	}

	currentReader, err := src.Next()
	if err != nil {
		return nil, fmt.Errorf("could not get next source: %w", err)
	}
	currentFileName := currentReader.(*os.File).Name()
	closeErr := currentReader.(*os.File).Close()
	if closeErr != nil {
		slog.Error("could not close current file reader", "error", closeErr)
	}

	fr, err := local.NewLocalFileReader(currentFileName)
	if err != nil {
		return nil, fmt.Errorf("could not create NewLocalFileReader\n%w", err)
	}

	pr, err := preader.NewParquetReader(fr, reflect.New(recordType).Interface(), 1)
	if err != nil {
		return nil, fmt.Errorf("could not create NewParquetReader\n%w", err)
	}

	newReader := &ParquetReader{
		currentParquetReader: pr,
		currentFileReader:    fr,
		currentFileName:      currentFileName,
		src:                  src,
		adapter:              adapter,
		recordType:           recordType,
	}

	for _, opt := range opts {
		opt(newReader)
	}

	return newReader, nil
}

func (r *ParquetReader) newSchema() any {
	return reflect.New(r.recordType).Interface()
}

func (r *ParquetReader) newRecordSlice(size int) any {
	sliceType := reflect.SliceOf(r.recordType)
	records := reflect.New(sliceType)
	records.Elem().Set(reflect.MakeSlice(sliceType, size, size))
	return records.Interface()
}

func (r *ParquetReader) ReadSingleFile(src *source.Source, currentReader *preader.ParquetReader) ([]any, error) {
	defer currentReader.ReadStop()

	records := r.newRecordSlice(int(currentReader.GetNumRows()))
	if err := currentReader.Read(records); err != nil {
		return nil, reader.NewReadError(err, src, "Failed to read parquet")
	}

	return convertToInputSchema(records), nil
}

func convertToInputSchema(records any) []any {
	value := reflect.ValueOf(records)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	converted := make([]any, value.Len())
	for i := range value.Len() {
		converted[i] = value.Index(i).Interface()
	}
	return converted
}

func (r *ParquetReader) Read() ([]any, error) {
	rows := make([]any, 0, int(r.currentParquetReader.GetNumRows()))
	eof := false

	for !eof {
		// Read the current file completely
		currentRows, err := r.ReadSingleFile(r.src, r.currentParquetReader)
		if err != nil {
			return nil, fmt.Errorf("could not read file: %s. %w", r.currentFileName, err)
		}
		rows = append(rows, currentRows...)

		// Get next file if any
		nextReader, err := r.src.Next()
		if err == io.EOF {
			// no more files
			return rows, io.EOF
		}
		if err != nil {
			return nil, fmt.Errorf("could not get next source: %w", err)
		}

		nextFileName := nextReader.(*os.File).Name()
		closeErr := nextReader.(*os.File).Close()
		if closeErr != nil {
			slog.Error("could not close current file reader", "error", closeErr)
		}
		// If no more files remaining, the returned error will be EOF
		eof = err == io.EOF

		// Create a new local file reader from the next file in the Source
		newFileReader, err := local.NewLocalFileReader(nextFileName)
		if err != nil {
			return nil, fmt.Errorf("could not create NewLocalFileReader\n%w", err)
		}

		// Create a new Parquet File Reader
		newParquetReader, err := preader.NewParquetReader(newFileReader, r.newSchema(), 1)
		if err != nil {
			return nil, fmt.Errorf("could not create NewParquetReader\n%w", err)
		}

		// Update the current file reader and close the previous
		r.currentParquetReader.ReadStop()
		closeErr = r.currentFileReader.Close()
		if closeErr != nil {
			slog.Error("could not close current file reader", "error", closeErr)
		}
		r.currentFileReader = newFileReader
		r.currentParquetReader = newParquetReader
		r.currentFileName = nextFileName
	}

	return rows, nil
}

// Reads batch from the passed reader
// The closing should be handling by the caller
// Returns io.EOF if there are no more rows to read
func (r *ParquetReader) ReadBatchSingleFile(currentReader *preader.ParquetReader) ([]any, error) {
	records := r.newRecordSlice(r.batchSize)

	if err := currentReader.Read(records); err != nil {
		return nil, reader.NewReadError(err, r.src, "Failed to read parquet in batch")
	}

	if isZeroValueSlice(records) {
		// finished reading
		return nil, io.EOF
	}

	return convertToInputSchema(records), nil
}

func (r *ParquetReader) ReadBatch() ([]any, error) {
	currentRows, err := r.ReadBatchSingleFile(r.currentParquetReader)

	// Read did not finish successfully
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error while reading batch from current parquet reader: %s. %w", r.currentFileName, err)
	}

	// EOF means we only finished reading the current file, but more files could remain
	if err == io.EOF {
		// Get the next file in the source
		nextReader, nextError := r.src.Next()
		// If error is not EOF it means there has been an error retrieving the next file from the Source
		if nextError != nil && nextError != io.EOF {
			return currentRows, fmt.Errorf("error getting the next file from the Source: %w", nextError)
		}
		// If error is EOF, it means that there are no more files to read
		if nextError == io.EOF {
			return currentRows, io.EOF
		}

		// Now we update the current reader and close the previous
		//
		// first get the next file name
		nextFileName := nextReader.(*os.File).Name()
		closeErr := nextReader.(*os.File).Close()
		if closeErr != nil {
			slog.Error("could not close current file reader", "error", closeErr)
		}

		// Create a new local file reader from the next file in the Source
		newFileReader, err := local.NewLocalFileReader(nextFileName)
		if err != nil {
			return nil, fmt.Errorf("could not create NewLocalFileReader\n%w", err)
		}

		// Create a new Parquet File Reader
		newParquetReader, err := preader.NewParquetReader(newFileReader, r.newSchema(), 1)
		if err != nil {
			return nil, fmt.Errorf("could not create NewParquetReader\n%w", err)
		}

		// update the current values
		r.currentParquetReader.ReadStop()
		closeErr = r.currentFileReader.Close()
		if closeErr != nil {
			slog.Error("could not close current file reader", "error", closeErr)
		}
		r.currentFileReader = newFileReader
		r.currentParquetReader = newParquetReader
		r.currentFileName = nextFileName

		// return the latest rows
		return currentRows, nil
	}

	// Here we are still reading the current file.
	return currentRows, nil
}

func isZeroValueSlice(records any) bool {
	value := reflect.ValueOf(records)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for i := range value.Len() {
		if !isZeroValueInputSchema(value.Index(i).Interface()) {
			return false
		}
	}
	return true
}

func isZeroValueInputSchema(schema any) bool {
	elem := reflect.ValueOf(schema)
	if elem.Kind() == reflect.Pointer {
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

func (r *ParquetReader) Close() error {
	r.currentParquetReader.ReadStop()
	return r.currentFileReader.Close()
}
