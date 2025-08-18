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
	"context"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type PartitionWriter struct {
	*writer.BaseWriter[repository.InputSchema, repository.InputSchema]
	cfg           *config.PartitionWriterConfig
	parquetStore  *ParquetStore
	inMemoryStore *InMemoryStore

	dirMap map[string]int
}

func New(
	cfg *config.PartitionWriterConfig,
	inputChan chan writer.WriterInput[repository.InputSchema],
	doneChan chan struct{},
	ctx context.Context,
) *PartitionWriter {
	fs := filesystemmanager.FileSystemManager{
		BaseDir: cfg.BaseDir,
		Handler: partition.PartitionHandler{
			NumPartitions:   cfg.NumPartitions,
			PartitionLevels: cfg.PartitionLevels,
		},
	}
	w := &PartitionWriter{
		BaseWriter: &writer.BaseWriter[repository.InputSchema, repository.InputSchema]{
			InboxChannel: inputChan,
			DoneChannel:  doneChan,
			Ctx:          ctx,
		},
		cfg: cfg,
		parquetStore: &ParquetStore{
			fs:      &fs,
			maxSize: cfg.MaxFileSize,
			schema:  cfg.Schema,
		},
		inMemoryStore: NewInMemoryStore(cfg.InMemoryMaxPartitionSize, &fs),
		dirMap:        make(map[string]int),
	}
	w.Writer = w
	return w
}

func (w *PartitionWriter) Receive(msg writer.WriterInput[repository.InputSchema]) {
	if msg.Error != nil {
		panic(msg.Error)
	}

	rowsToFlush, err := w.inMemoryStore.Write(msg.Rows)
	if err != nil {
		panic(err)
	}

	if len(rowsToFlush) > 0 {
		if dirMap, err := w.parquetStore.Write(rowsToFlush, w.dirMap); err != nil {
			panic(err)
		} else {
			w.dirMap = dirMap
		}
	}
}

func (w *PartitionWriter) Stop() {
	w.Flush()
	w.DoneChannel <- struct{}{}
	close(w.DoneChannel)
}

func (w *PartitionWriter) Flush() {
	w.parquetStore.Write(w.inMemoryStore.store, w.dirMap)
}
