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

package filesystemmanager

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
)

type FileSystemManagerBuilder struct {
	numPartitions int
	numLevels     int
	baseDir       string

	t *testing.T
}

func AFileSystemManager(t *testing.T) *FileSystemManagerBuilder {
	t.Helper()

	return &FileSystemManagerBuilder{
		numPartitions: 16,
		numLevels:     2,
		baseDir:       t.TempDir(),
		t:             t,
	}
}

func (b *FileSystemManagerBuilder) Build() FileSystemManager {
	b.t.Helper()

	handler := partition.PartitionHandler{
		NumPartitions:   b.numPartitions,
		PartitionLevels: b.numLevels,
	}

	return FileSystemManager{
		Handler: handler,
		BaseDir: b.baseDir,
	}
}

func (b *FileSystemManagerBuilder) WithBaseDir(baseDir string) *FileSystemManagerBuilder {
	b.t.Helper()

	b.baseDir = baseDir
	return b
}

func (b *FileSystemManagerBuilder) WithNumPartitions(n int) *FileSystemManagerBuilder {
	b.t.Helper()

	b.numPartitions = n
	return b
}

func (b *FileSystemManagerBuilder) WithNumLevels(n int) *FileSystemManagerBuilder {
	b.t.Helper()

	b.numLevels = n
	return b
}
