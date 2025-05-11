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
)

type PartitionReader struct {
	dirChannel chan<- string
}

func NewPartitionReader(dirChannel chan<- string) *PartitionReader {
	return &PartitionReader{
		dirChannel: dirChannel,
	}
}

// Recursively traverse the directory tree
// sending each leaf directory to the output channel
func (pr *PartitionReader) TraversePartitions(baseDir string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err := pr.TraversePartitions(filepath.Join(baseDir, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			pr.dirChannel <- baseDir
			return nil
		}
	}
	return nil
}
