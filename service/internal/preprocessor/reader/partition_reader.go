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
	"os"
	"path/filepath"
	"sync"
)

type PartitionReader struct {
	dirChannel chan<- string
	workers    []*Worker
	baseDir    string
}

func NewPartitionReader(dirChannel chan<- string, workers []*Worker, baseDir string) *PartitionReader {
	return &PartitionReader{
		dirChannel: dirChannel,
		workers:    workers,
		baseDir:    baseDir,
	}
}

// Recursively traverse the directory tree
// sending each leaf directory to the output channel
func (pr *PartitionReader) TraversePartitions(baseDir string) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pr.TraversePartitions(filepath.Join(baseDir, entry.Name()))
		} else {
			pr.dirChannel <- baseDir
			return
		}
	}
}

// Start the partition reader
// This will traverse the directory tree and send each leaf directory to the output channel
// The workers will read from the directory channel and send the records to the output channel
// The directory channel and worker output channel will be closed when the workers are done
func (pr *PartitionReader) Start() {
	go pr.TraversePartitions(pr.baseDir)

	wg := sync.WaitGroup{}
	wg.Add(len(pr.workers))
	for _, worker := range pr.workers {
		go worker.Start(&wg)
	}
	wg.Wait()

	close(pr.dirChannel)
	close(pr.workers[0].output)
}
