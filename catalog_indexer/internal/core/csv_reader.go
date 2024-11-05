package core

import (
	"encoding/csv"
	"io"
	"log/slog"
)

type CsvReader struct {
	Header    []string
	csvReader *csv.Reader
	src       Source
	BatchSize int
}

type CsvReaderOption func(r *CsvReader)

func WithHeader(header []string) CsvReaderOption {
	return func(r *CsvReader) {
		r.Header = header
	}
}

func WithBatchSize(size int) CsvReaderOption {
	return func(r *CsvReader) {
		if size <= 0 {
			size = 1
		}
		r.BatchSize = size
	}
}

func NewCsvReader(src Source, opts ...CsvReaderOption) (*CsvReader, error) {
	reader := CsvReader{
		csvReader: csv.NewReader(src.Reader),
		src:       src,
	}
	for _, opt := range opts {
		opt(&reader)
	}
	return &reader, nil
}

func (r *CsvReader) Read() ([]Row, error) {
	rows := make([]Row, 0, 0)
	if r.Header == nil {
		header, err := r.csvReader.Read()
		if err != nil {
			slog.Error("Could not read header from csv.", "reader", r.csvReader)
			return nil, err
		}
		r.Header = header
	}
	records, err := r.csvReader.ReadAll()
	if err != nil {
		slog.Error("Could not read contents from csv.", "reader", r.csvReader)
		return nil, err
	}
	for _, record := range records {
		row := make(Row)
		for i, h := range r.Header {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (r *CsvReader) ReadBatch() ([]Row, error) {
	rows := make([]Row, 0, 0)
	if r.Header == nil {
		header, err := r.csvReader.Read()
		if err != nil {
			slog.Error("Could not read header from csv.", "reader", r.csvReader)
			return nil, err
		}
		r.Header = header
	}
	for i := 0; i < r.BatchSize; i++ {
		record, err := r.csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, err
			}
			return nil, err
		}
		row := make(Row)
		for i, h := range r.Header {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}
