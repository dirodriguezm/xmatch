package reader

import (
	"encoding/csv"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
)

type CsvReader struct {
	Header          []string
	FirstLineHeader bool
	csvReaders      []*csv.Reader
	currentReader   int
	src             *source.Source
	BatchSize       int
	outbox          chan indexer.ReaderResult
}

func NewCsvReader(src *source.Source, channel chan indexer.ReaderResult, opts ...CsvReaderOption) (*CsvReader, error) {
	readers := []*csv.Reader{}
	for _, reader := range src.Reader {
		readers = append(readers, csv.NewReader(reader))
	}
	reader := CsvReader{
		csvReaders:    readers,
		currentReader: 0,
		src:           src,
		outbox:        channel,
	}
	for _, opt := range opts {
		opt(&reader)
	}
	return &reader, nil
}

func (r *CsvReader) Start() {
	go func() {
		defer close(r.outbox)
		eof := false
		for !eof {
			rows, err := r.ReadBatch()
			if err != nil && err != io.EOF {
				readResult := indexer.ReaderResult{
					Rows:  nil,
					Error: err,
				}
				r.outbox <- readResult
				return
			}
			eof = err == io.EOF
			readResult := indexer.ReaderResult{
				Rows:  rows,
				Error: nil,
			}
			r.outbox <- readResult
		}
	}()
}

func (r *CsvReader) ReadSingleFile(currentReader *csv.Reader) ([]indexer.Row, error) {
	rows := make([]indexer.Row, 0, 0)
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
		row := make(indexer.Row)
		for i, h := range r.Header {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (r *CsvReader) Read() ([]indexer.Row, error) {
	rows := make([]indexer.Row, 0, 0)
	for _, currentReader := range r.csvReaders {
		currentRows, err := r.ReadSingleFile(currentReader)
		if err != nil {
			return nil, err
		}
		rows = append(rows, currentRows...)
	}
	return rows, nil
}

func (r *CsvReader) ReadBatchSingleFile(currentReader *csv.Reader) ([]indexer.Row, error) {
	rows := make([]indexer.Row, 0, 0)
	if r.Header == nil {
		header, err := currentReader.Read()
		if err != nil {
			slog.Error("Could not read header from csv.", "reader", r.csvReaders)
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
		row := make(indexer.Row)
		for i, h := range r.Header {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (r *CsvReader) ReadBatch() ([]indexer.Row, error) {
	rows := make([]indexer.Row, 0, 0)
	currentReader := r.csvReaders[r.currentReader]
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

func (r *CsvReader) ObjectIdCol() string {
	return r.src.OidCol
}

func (r *CsvReader) RaCol() string {
	return r.src.RaCol
}

func (r *CsvReader) DecCol() string {
	return r.src.DecCol
}

func (r *CsvReader) Catalog() string {
	return r.src.CatalogName
}
