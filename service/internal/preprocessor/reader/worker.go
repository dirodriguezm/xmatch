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

package partition_reader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

type Records = []repository.InputSchema

type Worker struct {
	dirChannel <-chan string
	schema     config.ParquetSchema
	output     chan<- Records
}

func NewWorker(dirChannel <-chan string, schema config.ParquetSchema, output chan<- Records) *Worker {
	return &Worker{
		dirChannel: dirChannel,
		schema:     schema,
		output:     output,
	}
}

func (w *Worker) Start() {
	for dir := range w.dirChannel {
		records, err := w.readDirectory(dir)
		if err != nil {
			panic(err)
		}

		groupedRecords := w.groupByOid(records)
		for _, records := range groupedRecords {
			w.output <- records
		}
	}
	close(w.output)
}

func (w *Worker) readDirectory(dir string) (Records, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	records := make(Records, 0)
	for _, file := range files {
		if file.IsDir() {
			panic(fmt.Errorf("found directory in partition: %s", file.Name()))
		}

		newRecords, err := w.readFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, newRecords...)
	}

	return records, nil
}

func (w *Worker) readFile(file string) (Records, error) {
	fr, err := local.NewLocalFileReader(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", file, err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, getParquetSchema(w.schema), 1)
	if err != nil {
		return nil, fmt.Errorf("error creating parquet reader: %w", err)
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	switch w.schema {
	case config.AllwiseSchema:
		records := make([]repository.AllwiseInputSchema, num)
		if err = pr.Read(&records); err != nil {
			return nil, fmt.Errorf("error reading parquet file: %w", err)
		}
		return convertAllwiseToInputSchema(records), nil
	case config.TestSchema:
		records := make([]TestInputSchema, num)
		if err = pr.Read(&records); err != nil {
			return nil, fmt.Errorf("error reading parquet file: %w", err)
		}
		return convertTestToInputSchema(records), nil
	default:
		panic("Unknown schema")
	}
}

func (w *Worker) groupByOid(records Records) map[string]Records {
	result := make(map[string]Records)
	for i := 0; i < len(records); i++ {
		result[records[i].GetId()] = append(result[records[i].GetId()], records[i])
	}
	return result
}

func getParquetSchema(schema config.ParquetSchema) any {
	switch schema {
	case config.AllwiseSchema:
		return new(repository.AllwiseInputSchema)
	case config.TestSchema:
		return new(TestInputSchema)
	default:
		panic("Unknown schema")
	}
}

func convertAllwiseToInputSchema(records []repository.AllwiseInputSchema) []repository.InputSchema {
	result := make([]repository.InputSchema, len(records))
	for i, r := range records {
		result[i] = &r
	}
	return result
}

func convertTestToInputSchema(records []TestInputSchema) []repository.InputSchema {
	result := make([]repository.InputSchema, len(records))
	for i, r := range records {
		result[i] = &r
	}
	return result
}
