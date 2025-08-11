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

package csv_reader

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type CsvReader struct {
	*reader.BaseReader
	Header          []string
	FirstLineHeader bool
	currentReader   *csv.Reader
}

func NewCsvReader(src *source.Source, channel []chan reader.ReaderResult, opts ...CsvReaderOption) (*CsvReader, error) {
	currentReader, err := src.Next()
	if err != nil {
		return nil, fmt.Errorf("could not get next source: %w", err)
	}

	r := &CsvReader{
		BaseReader: &reader.BaseReader{
			Src:    src,
			Outbox: channel,
		},
		currentReader: csv.NewReader(currentReader),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.Reader = r

	return r, nil
}

func (r *CsvReader) ReadSingleFile(currentReader *csv.Reader) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.BatchSize)

	// Read the header if not already read
	if r.Header == nil {
		header, err := currentReader.Read()
		if err != nil {
			return nil, fmt.Errorf("could not read header from csv: %w", err)
		}
		r.Header = header
	}

	// Read the rest of the file
	records, err := currentReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read all records from csv: %w", err)
	}

	// Transform data into the correct schema
	for _, record := range records {
		row := r.createInputSchema(r.Src.CatalogName)
		for i, h := range r.Header {
			row.SetField(h, record[i])
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func (r *CsvReader) Read() ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0)

	eof := false
	for !eof {
		currentRows, readErr := r.ReadSingleFile(r.currentReader)
		if readErr != nil {
			return nil, readErr
		}

		rows = append(rows, currentRows...)

		ioReader, err := r.Src.Next()
		eof = err == io.EOF
		if err != nil && !eof {
			return nil, fmt.Errorf("Could not get next source: %w", err)
		}

		r.currentReader = csv.NewReader(ioReader)
	}

	return rows, nil
}

func (r *CsvReader) ReadBatchSingleFile(currentReader *csv.Reader) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.BatchSize)

	if r.Header == nil {
		header, err := currentReader.Read()
		if err != nil {
			return nil, fmt.Errorf("Could not read header from csv: %w", err)
		}
		r.Header = header
	}

	for range r.BatchSize {
		record, err := currentReader.Read()
		if err == io.EOF {
			return rows, err
		}
		if err != nil {
			return nil, err
		}

		row := r.createInputSchema(r.Src.CatalogName)

		for j, h := range r.Header {
			row.SetField(h, record[j])
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// ReadBatch reads a batch of records from the current CSV reader.
// It processes records in batches according to the BatchSize.
// If the end of the file is reached, it retrieves the next source if available.
// Returns the processed rows or an error, including EOF if the end of the last file is reached.
func (r *CsvReader) ReadBatch() ([]repository.InputSchema, error) {
	// Initialize the result slice. Right now, the last batch of a file could have
	// less than BatchSize rows. Maybe the slice could have a fixed size, but it's not too important currently.
	rows := make([]repository.InputSchema, 0, r.BatchSize)

	// Create CSV reader from file reader and read a batch
	currentRows, err := r.ReadBatchSingleFile(r.currentReader)

	// If the error is EOF, we get the next reader from the Source.
	// And if there is no next reader, we return the rows we have so far.
	if err == io.EOF {
		rows = append(rows, currentRows...)
		if r.FirstLineHeader {
			r.Header = nil
		}

		var nextErr error
		ioReader, nextErr := r.Src.Next()
		if nextErr != nil {
			return rows, nextErr // This error can potentially be EOF, handled by the caller.
		}
		r.currentReader = csv.NewReader(ioReader)

		return rows, nil
	}

	// If the error is not EOF, it's a real error.
	if err != nil {
		return nil, fmt.Errorf("could not read batch from csv: %w", err)
	}

	// Read batch successfully and more to read
	rows = append(rows, currentRows...)
	return rows, nil
}

func (r *CsvReader) createInputSchema(catalogName string) repository.InputSchema {
	switch catalogName {
	case "allwise":
		return &repository.AllwiseInputSchema{}
	case "vlass":
		return &repository.VlassInputSchema{}
	default:
		return &TestSchema{}
	}
}
