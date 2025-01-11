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
		handler: handler,
		baseDir: b.baseDir,
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
