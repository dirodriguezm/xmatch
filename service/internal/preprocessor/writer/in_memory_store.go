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
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type InMemoryStore struct {
	// map of partition directory to rows
	store map[string][]repository.InputSchema
	// max number of rows for each partition
	maxPartitionSize int
	// filesystem manager to get partitions for each oid
	fs *filesystemmanager.FileSystemManager
}

func NewInMemoryStore(maxPartitionSize int, fs *filesystemmanager.FileSystemManager) *InMemoryStore {
	return &InMemoryStore{
		store:            make(map[string][]repository.InputSchema),
		maxPartitionSize: maxPartitionSize,
		fs:               fs,
	}
}

// Writes rows to the store
//
// Returns the rows that are ready to be flushed
func (s *InMemoryStore) Write(rowsToWrite []repository.InputSchema) (map[string][]repository.InputSchema, error) {
	toFlush := make(map[string][]repository.InputSchema)

	for _, row := range rowsToWrite {
		dir, err := s.fs.GetDirectory(row.GetId())
		if err != nil {
			return nil, err
		}

		if !s.canWrite(dir) {
			toFlush[dir] = append(toFlush[dir], row)
			continue
		}

		s.store[dir] = append(s.store[dir], row)
	}

	return s.flush(toFlush), nil
}

func (s *InMemoryStore) canWrite(dir string) bool {
	return len(s.store[dir]) < s.maxPartitionSize
}

func (s *InMemoryStore) flush(toFlush map[string][]repository.InputSchema) map[string][]repository.InputSchema {
	for dir := range toFlush {
		toFlush[dir] = append(s.store[dir], toFlush[dir]...)
		s.store[dir] = make([]repository.InputSchema, 0)
	}

	return toFlush
}
