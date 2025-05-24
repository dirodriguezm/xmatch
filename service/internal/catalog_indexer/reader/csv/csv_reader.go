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
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type CsvReader struct {
	*reader.BaseReader
	Header          []string
	FirstLineHeader bool
	csvReaders      []*csv.Reader
	currentReader   int
}

func NewCsvReader(src *source.Source, channel []chan reader.ReaderResult, opts ...CsvReaderOption) (*CsvReader, error) {
	readers := []*csv.Reader{}
	for _, reader := range src.Reader {
		readers = append(readers, csv.NewReader(reader))
	}
	r := &CsvReader{
		csvReaders:    readers,
		currentReader: 0,
		BaseReader: &reader.BaseReader{
			Src:    src,
			Outbox: channel,
		},
	}
	for _, opt := range opts {
		opt(r)
	}
	r.Reader = r
	return r, nil
}

func (r *CsvReader) ReadSingleFile(currentReader *csv.Reader) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.BatchSize)
	if r.Header == nil {
		header, err := currentReader.Read()
		if err != nil {
			slog.Error("Could not read header from csv.", "reader", r.csvReaders)
			return nil, err
		}
		r.Header = header
	}
	records, err := currentReader.ReadAll()
	if err != nil {
		slog.Error("Could not read contents from csv.", "reader", r.csvReaders)
		return nil, err
	}
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
	rows := make([]repository.InputSchema, 0, 0)
	for _, currentReader := range r.csvReaders {
		currentRows, err := r.ReadSingleFile(currentReader)
		if err != nil {
			return nil, err
		}
		rows = append(rows, currentRows...)
	}
	return rows, nil
}

func (r *CsvReader) ReadBatchSingleFile(currentReader *csv.Reader) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.BatchSize)
	if r.Header == nil {
		header, err := currentReader.Read()
		if err != nil {
			slog.Error("Could not read header from csv.", "reader", r.csvReaders, "error", err)
			return nil, err
		}
		r.Header = header
	}
	for i := 0; i < r.BatchSize; i++ {
		record, err := currentReader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, err
			}
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

func (r *CsvReader) ReadBatch() ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, 0)
	currentReader := r.csvReaders[r.currentReader]
	slog.Debug("CsvReader reading batch")
	currentRows, err := r.ReadBatchSingleFile(currentReader)
	if err != nil {
		if err == io.EOF && r.currentReader < len(r.csvReaders)-1 {
			rows = append(rows, currentRows...)
			r.currentReader += 1
			if r.FirstLineHeader {
				r.Header = nil
			}
			return rows, nil
		}
		if err == io.EOF {
			rows = append(rows, currentRows...)
			r.currentReader += 1
			return rows, err
		}
		return nil, err
	}
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
