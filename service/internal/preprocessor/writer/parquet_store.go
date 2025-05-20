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

package partition_writer

import (
	"fmt"
	"os"
	"path"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
)

type ParquetStore struct {
	fs      *filesystemmanager.FileSystemManager
	maxSize int
	schema  config.ParquetSchema
}

func (store *ParquetStore) Write(rowsToWrite map[string][]repository.InputSchema, dirMap map[string]int) (map[string]int, error) {
	for dir, rows := range rowsToWrite {
		if _, ok := dirMap[dir]; !ok {
			dirMap[dir] = 1
			_, err := store.createNewFile(dir, dirMap[dir], rows, false)
			if err != nil {
				return dirMap, fmt.Errorf("Error while creating new file %s\n%w", dir, err)
			}
			continue
		}

		file, err := store.fs.GetFile(dir, dirMap[dir])
		if err != nil {
			return dirMap, fmt.Errorf("Error while getting file\n%w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(fmt.Errorf("File could not close. %w", err))
			}
		}()

		if store.fs.GetSizeOfFile(file) < int64(store.maxSize) {
			if err := store.reuseFile(file, rows, dirMap); err != nil {
				return dirMap, fmt.Errorf("Error while reusing file\n%w", err)
			}
		} else {
			dirMap[dir]++
			if fname, err := store.createNewFile(dir, dirMap[dir], rows, false); err != nil {
				return dirMap, fmt.Errorf("Error while creating new file %s\n%w", fname, err)
			}
		}
	}

	return dirMap, nil
}

func (store *ParquetStore) reuseFile(file *os.File, rows []repository.InputSchema, dirMap map[string]int) error {
	records, err := store.readFile(file.Name())
	if err != nil {
		return err
	}

	records = append(records, rows...)

	dir, _ := path.Split(file.Name())
	newFile, err := store.createNewFile(dir, dirMap[dir], records, true)
	if err != nil {
		return err
	}

	if err := os.Remove(file.Name()); err != nil {
		return err
	}

	if err := os.Rename(newFile, file.Name()); err != nil {
		return err
	}

	return nil
}

func (store *ParquetStore) createNewFile(dir string, number int, rows []repository.InputSchema, tmp bool) (string, error) {
	var file *os.File
	var err error
	if tmp {
		file, err = store.fs.CreateNewTmpFile(dir, number)
	} else {
		file, err = store.fs.CreateNewFile(dir, number)
	}
	if err != nil {
		return "", fmt.Errorf("Error while creating file\n%w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("File could not close. %w", err))
		}
	}()

	pw, err := writer.NewParquetWriterFromWriter(file, getParquetSchema(store.schema), 4)
	if err != nil {
		return "", fmt.Errorf("ParquetGo could not create writer from writer\n%w", err)
	}

	defer func() {
		if err := pw.WriteStop(); err != nil {
			panic(fmt.Errorf("ParquetWriter could not stop. %w", err))
		}
	}()

	for _, row := range rows {
		if err := pw.Write(row); err != nil {
			return "", fmt.Errorf("ParquetWriter could not write object %v\n%w", row, err)
		}
	}

	return file.Name(), nil
}

func (store *ParquetStore) readFile(fname string) ([]repository.InputSchema, error) {
	fr, err := local.NewLocalFileReader(fname)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := fr.Close(); err != nil {
			panic(fmt.Errorf("File could not close. %w", err))
		}
	}()

	readerSchema := getParquetSchema(store.schema)

	pr, err := reader.NewParquetReader(fr, readerSchema, 4)
	if err != nil {
		return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
	}
	defer pr.ReadStop()

	nrows := pr.GetNumRows()

	switch store.schema {
	case config.AllwiseSchema:
		return store.readAllwiseSchema(pr, nrows)
	case config.VlassSchema:
		return store.readVlassSchema(pr, nrows)
	case config.TestSchema:
		return store.readTestSchema(pr, nrows)
	default:
		panic("Unknown schema")
	}
}

func (store *ParquetStore) readAllwiseSchema(pr *reader.ParquetReader, nrows int64) ([]repository.InputSchema, error) {
	records := make([]repository.AllwiseInputSchema, nrows)
	defer func() {
		records = nil
	}()

	if err := pr.Read(&records); err != nil {
		return nil, err
	}

	result := make([]repository.InputSchema, len(records))
	for i, r := range records {
		result[i] = &r
	}

	return result, nil
}

func (store *ParquetStore) readVlassSchema(pr *reader.ParquetReader, nrows int64) ([]repository.InputSchema, error) {
	records := make([]repository.VlassInputSchema, nrows)
	defer func() {
		records = nil
	}()

	if err := pr.Read(&records); err != nil {
		return nil, err
	}

	result := make([]repository.InputSchema, len(records))
	for i, r := range records {
		result[i] = &r
	}

	return result, nil
}

func (store *ParquetStore) readTestSchema(pr *reader.ParquetReader, nrows int64) ([]repository.InputSchema, error) {
	records := make([]TestInputSchema, nrows)
	defer func() {
		records = nil
	}()

	if err := pr.Read(&records); err != nil {
		return nil, err
	}

	result := make([]repository.InputSchema, len(records))
	for i, r := range records {
		result[i] = r
	}

	return result, nil
}

func getParquetSchema(schema config.ParquetSchema) any {
	switch schema {
	case config.AllwiseSchema:
		return new(repository.AllwiseInputSchema)
	case config.TestSchema:
		return new(TestInputSchema)
	case config.VlassSchema:
		return new(repository.VlassInputSchema)
	default:
		panic("Unknown schema")
	}
}
