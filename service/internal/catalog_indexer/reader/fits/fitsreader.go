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

// Package fitsreader provides a reader for FITS files
package fitsreader

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"codeberg.org/astrogo/fitsio"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type FitsReader[T repository.InputSchema] struct {
	currentFileReader io.ReadCloser
	currentFitsRows   *fitsio.Rows
	currentFitsFile   *fitsio.File
	src               *source.Source
	batchSize         int
}

func NewFitsReader[T repository.InputSchema](src *source.Source, opts ...FitsReaderOption[T]) (*FitsReader[T], error) {
	currentFileReader, err := src.Next()
	if err != nil {
		return nil, err
	}

	fits, err := fitsio.Open(currentFileReader)
	if err != nil {
		return nil, err
	}
	table := fits.HDU(1).(*fitsio.Table)
	rows, err := table.Read(0, table.NumRows())
	if err != nil {
		return nil, err
	}

	r := &FitsReader[T]{
		currentFileReader: currentFileReader,
		currentFitsRows:   rows,
		currentFitsFile:   fits,
		src:               src,
		batchSize:         1,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func (r *FitsReader[T]) Read() ([]repository.InputSchema, error) {
	panic("Not implemented")
}

func (r *FitsReader[T]) ReadBatch() ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, r.batchSize)

	currentRows, err := r.ReadBatchSingleFile(r.currentFitsRows, r.batchSize)

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

	if err != nil {
		return nil, fmt.Errorf("could not read batch from csv: %w", err)
	}

	rows = append(rows, currentRows...)
	return rows, nil
}

func (r *FitsReader[T]) closeResources() error {
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

func (r *FitsReader[T]) switchToNewFile() error {
	ioReader, err := r.src.Next()
	if err != nil {
		return err
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

func (r *FitsReader[T]) ReadBatchSingleFile(rowIterator *fitsio.Rows, size int) ([]repository.InputSchema, error) {
	rows := make([]repository.InputSchema, 0, size)
	for range size {
		hasNext := rowIterator.Next()
		if !hasNext {
			return rows, io.EOF
		}

		row := r.createInputSchema(rowIterator)
		rows = append(rows, row)
	}

	return rows, nil
}

func (r *FitsReader[T]) createInputSchema(rowIterator *fitsio.Rows) repository.InputSchema {
	adapter, err := catalog.GetFactory(r.src.CatalogName)
	if err != nil {
		var schema T
		if err := rowIterator.Scan(&schema); err != nil {
			panic(err)
		}
		return schema
	}

	schemaPtr := reflect.New(reflect.TypeOf(adapter.NewInputSchema()))
	if err := rowIterator.Scan(schemaPtr.Interface()); err != nil {
		panic(err)
	}
	return schemaPtr.Elem().Interface().(repository.InputSchema)
}

func (r *FitsReader[T]) Close() error {
	return errors.Join(
		r.currentFitsRows.Close(),
		r.currentFitsFile.Close(),
	)
}
