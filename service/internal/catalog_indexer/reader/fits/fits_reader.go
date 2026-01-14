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

package fits_reader

import (
	"errors"
	"fmt"
	"io"

	"codeberg.org/astrogo/fitsio"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type FitsReader struct {
	currentFileReader io.ReadCloser
	currentFitsRows   *fitsio.Rows
	currentFitsFile   *fitsio.File
	src               *source.Source
	batchSize         int
}

func NewFitsReader(src *source.Source, opts ...FitsReaderOption) (FitsReader, error) {
	currentFileReader, err := src.Next()
	if err != nil {
		return FitsReader{}, err
	}

	fits, err := fitsio.Open(currentFileReader)
	if err != nil {
		return FitsReader{}, err
	}
	table := fits.HDU(1).(*fitsio.Table)
	rows, err := table.Read(0, table.NumRows())
	if err != nil {
		return FitsReader{}, err
	}

	r := FitsReader{
		currentFileReader: currentFileReader,
		currentFitsRows:   rows,
		currentFitsFile:   fits,
		src:               src,
		batchSize:         1,
	}

	for _, opt := range opts {
		r = opt(r)
	}

	return r, nil
}

func (r *FitsReader) Read() ([]repository.InputSchema, error) {
	panic("Not implemented")
}

func (r *FitsReader) ReadBatch() ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.batchSize)

	currentRows, err := r.ReadBatchSingleFile(r.currentFitsRows, r.batchSize, r.src.CatalogName)

	// If the error is EOF, we get the next reader from the Source.
	// And if there is no next reader, we return the rows we have so far.
	if err == io.EOF {
		rows = append(rows, currentRows...)

		err = r.closeResources()
		if err != nil {
			return rows, err
		}

		err = r.switchToNewFile()
		if err != nil {
			return rows, err
		}

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

func (r *FitsReader) closeResources() error {
	err := r.currentFitsRows.Close()
	if err != nil {
		return fmt.Errorf("failed to close FITS rows: %w", err)
	}

	err = r.currentFitsFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close FITS file: %w", err)
	}

	err = r.currentFileReader.Close()
	if err != nil {
		return fmt.Errorf("failed to close file reader: %w", err)
	}

	return nil
}

func (r *FitsReader) switchToNewFile() error {
	ioReader, err := r.src.Next()
	if err != nil {
		return err // This error can potentially be EOF, handled by the caller.
	}
	r.currentFileReader = ioReader

	r.currentFitsFile, err = fitsio.Open(ioReader)
	if err != nil {
		return err
	}

	table := r.currentFitsFile.HDU(0).(*fitsio.Table)
	r.currentFitsRows, err = table.Read(0, table.NumRows())
	if err != nil {
		return err
	}

	return nil
}

func (r *FitsReader) ReadBatchSingleFile(rowIterator *fitsio.Rows, size int, name string) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, size)
	for range size {
		hasNext := rowIterator.Next()
		if !hasNext {
			return rows, io.EOF
		}

		row := r.createInputSchema(name, rowIterator)
		rows = append(rows, row)
	}

	return rows, nil
}

func (r *FitsReader) createInputSchema(name string, rowIterator *fitsio.Rows) repository.InputSchema {
	switch name {
	case "allwise":
		schema := repository.AllwiseInputSchema{}
		err := rowIterator.Scan(&schema)
		if err != nil {
			panic(err)
		}
		return schema
	default:
		schema := TestSchema{}
		err := rowIterator.Scan(&schema)
		if err != nil {
			panic(err)
		}
		return schema
	}
}

func (r *FitsReader) Close() error {
	return errors.Join(
		r.currentFitsRows.Close(),
		r.currentFitsFile.Close(),
	)
}
