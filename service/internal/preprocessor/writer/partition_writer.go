// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	xwaveWriter "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	parquetWriter "github.com/xitongsys/parquet-go/writer"
)

type PartitionWriter struct {
	*xwaveWriter.BaseWriter[repository.InputSchema]
	fs            *filesystemmanager.FileSystemManager
	maxFileSize   int
	currentWriter *parquetWriter.ParquetWriter
	currentFile   *os.File
	dirMap        map[string]int
	currentDir    string
	cfg           *config.PartitionWriterConfig
}

func New(
	cfg *config.PartitionWriterConfig,
	inputChan chan writer.WriterInput[repository.InputSchema],
	doneChan chan struct{},
) (*PartitionWriter, error) {
	if cfg.MaxFileSize <= 0 {
		return nil, fmt.Errorf("MaxFileSize must be greater than 0")
	}
	if cfg.NumPartitions <= 0 {
		return nil, fmt.Errorf("NumPartitions must be greater than 0")
	}
	if cfg.PartitionLevels <= 0 {
		return nil, fmt.Errorf("PartitionLevels must be greater than 0")
	}
	if cfg.BaseDir == "" {
		return nil, fmt.Errorf("BaseDir must not be empty")
	}

	w := &PartitionWriter{
		BaseWriter: &xwaveWriter.BaseWriter[repository.InputSchema]{
			InboxChannel: inputChan,
			DoneChannel:  doneChan,
		},
		fs: &filesystemmanager.FileSystemManager{
			Handler: partition.PartitionHandler{
				NumPartitions:   cfg.NumPartitions,
				PartitionLevels: cfg.PartitionLevels,
			},
			BaseDir: cfg.BaseDir,
		},
		maxFileSize: cfg.MaxFileSize,
		dirMap:      make(map[string]int),
		cfg:         cfg,
	}
	w.Writer = w
	return w, nil
}

func (w *PartitionWriter) Receive(msg xwaveWriter.WriterInput[repository.InputSchema]) {
	if msg.Error != nil {
		panic(msg.Error)
	}

	err := w.write(msg.Rows, w.cfg.Schema)
	if err != nil {
		panic(err)
	}
}

func (w *PartitionWriter) write(rows []repository.InputSchema, schema config.ParquetWriterSchema) error {
	idMap := w.idMap(rows)

	for id := range idMap {
		dir, err := w.fs.GetDirectory(id)
		if err != nil {
			return fmt.Errorf("Error while getting directory for id %s\n%w", id, err)
		}

		dirChanged := w.currentDir != dir
		if dirChanged {
			w.currentDir = dir
			_, err := w.fs.CreatePartition(id)
			if err != nil {
				return fmt.Errorf("Error while creating partition for id %s\n%w", id, err)
			}
		}

		if w.currentFile == nil || w.fs.GetSizeOfFile(w.currentFile) > int64(w.maxFileSize) || dirChanged {
			err := w.updateCurrentWriter(w.currentDir, schema)
			if err != nil {
				return fmt.Errorf("Error while updating current writer and file\n%w", err)
			}
		}

		for _, rowIdx := range idMap[id] {
			if err := w.currentWriter.Write(rows[rowIdx]); err != nil {
				return fmt.Errorf("ParquetWriter could not write object %v\n%w", rows[rowIdx], err)
			}
			if w.fs.GetSizeOfFile(w.currentFile) > int64(w.maxFileSize) {
				err := w.updateCurrentWriter(w.currentDir, schema)
				if err != nil {
					return fmt.Errorf("Error updating current writer and file due to file size\n%w", err)
				}
			}
		}
	}

	return nil
}

func (w *PartitionWriter) Stop() {
	if err := w.closeResources(); err != nil {
		panic(err)
	}
	w.DoneChannel <- struct{}{}
	close(w.DoneChannel)
}

// Create a map of object ids and the indexes of their rows
func (w *PartitionWriter) idMap(rows []repository.InputSchema) map[string][]int {
	idMap := make(map[string][]int)

	for i, row := range rows {
		id := row.GetId()
		if _, ok := idMap[id]; !ok {
			idMap[id] = make([]int, 0)
		}
		idMap[id] = append(idMap[id], i)
	}

	return idMap
}

func (w *PartitionWriter) updateCurrentWriter(dir string, schema config.ParquetWriterSchema) error {
	if err := w.closeResources(); err != nil {
		return fmt.Errorf(
			"Error closing resources while changing to a new file\n%w",
			err,
		)
	}

	updateCurrentDirFileIndex(w.dirMap, dir)

	var err error
	w.currentFile, err = w.fs.CreateNewFile(dir, w.dirMap[dir])
	if err != nil {
		return fmt.Errorf("Error while creating file\n%w", err)
	}

	w.currentWriter, err = w.createParquetWriter(w.currentFile, schema)
	if err != nil {
		return fmt.Errorf("ParquetGo could not create writer from writer\n%w", err)
	}

	return nil
}

func (w *PartitionWriter) createParquetWriter(file *os.File, schema config.ParquetWriterSchema) (*parquetWriter.ParquetWriter, error) {
	var writerSchema any
	switch schema {
	case config.AllwiseSchema:
		writerSchema = new(repository.AllwiseInputSchema)
	case config.TestSchema:
		writerSchema = new(TestInputSchema)
	default:
		panic("Unknown schema")
	}
	writer, err := parquetWriter.NewParquetWriterFromWriter(file, writerSchema, 1)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *PartitionWriter) closeResources() error {
	if w.currentWriter == nil || w.currentFile == nil {
		return nil
	}

	if err := w.currentWriter.WriteStop(); err != nil {
		return fmt.Errorf("ParquetWriter could not stop. %w", err)
	}

	if err := w.currentFile.Close(); err != nil {
		return fmt.Errorf("ParquetWriter could not close parquet file %w", err)
	}

	return nil
}

func updateCurrentDirFileIndex(dirMap map[string]int, dir string) {
	if _, ok := dirMap[dir]; !ok {
		dirMap[dir] = 1
	} else {
		dirMap[dir]++
	}
}
