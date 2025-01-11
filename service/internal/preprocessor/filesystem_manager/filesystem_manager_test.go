package filesystemmanager

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
	"github.com/stretchr/testify/require"
)

func TestFileSystemManager_TestCreatePartitionDirectory(t *testing.T) {
	dir := t.TempDir()

	part := partition.Partition{Levels: []int{0, 0}}
	handler := partition.PartitionHandler{PartitionLevels: 2, NumPartitions: 16}
	manager := FileSystemManager{
		handler: handler,
		baseDir: dir,
	}

	_, err := manager.createPartitionDirectory(part)
	require.NoError(t, err)

	require.DirExists(t, path.Join(dir, "00"))
	require.DirExists(t, path.Join(dir, "00", "00"))

	dir = t.TempDir()
	manager.baseDir = dir
	part = partition.Partition{Levels: []int{1, 0}}

	_, err = manager.createPartitionDirectory(part)
	require.NoError(t, err)

	require.DirExists(t, path.Join(dir, "01"))
	require.DirExists(t, path.Join(dir, "01", "00"))

	dir = t.TempDir()
	manager.baseDir = dir
	part = partition.Partition{Levels: []int{11, 12}}

	_, err = manager.createPartitionDirectory(part)
	require.NoError(t, err)

	require.DirExists(t, path.Join(dir, "11"))
	require.DirExists(t, path.Join(dir, "11", "12"))
}

func TestFileSystemManager_TestCreatePartitionDirectory_DirectoryAlreadyExistDoesNotFail(t *testing.T) {
	dir := t.TempDir()
	err := os.Mkdir(path.Join(dir, "00"), 0777)
	require.NoError(t, err)

	part := partition.Partition{Levels: []int{0, 0}}
	handler := partition.PartitionHandler{PartitionLevels: 2, NumPartitions: 16}
	manager := FileSystemManager{
		handler: handler,
		baseDir: dir,
	}

	_, err = manager.createPartitionDirectory(part)

	require.NoError(t, err)
	require.DirExists(t, path.Join(dir, "00"))
	require.DirExists(t, path.Join(dir, "00", "00"))
}

func TestFileSystemManager_TestAssignPartition(t *testing.T) {
	manager := AFileSystemManager(t).Build()

	part, err := manager.assignPartition("objectid")

	require.NoError(t, err)
	require.Len(t, part.Levels, 2)
	for _, level := range part.Levels {
		require.GreaterOrEqual(t, level, 0)
		require.Less(t, level, 16)
	}
}

func TestFileSystemManager_TestGetDirectory(t *testing.T) {
	manager := AFileSystemManager(t).WithBaseDir("test").Build()

	dir, err := manager.GetDirectory("objectid")

	require.NoError(t, err)
	require.Contains(t, dir, "test/")
	require.Len(t, strings.Split(dir, "/"), 3)
}

func TestFileSystemManager_TestCreatePartition(t *testing.T) {
	manager := AFileSystemManager(t).Build()

	createdDir, err := manager.CreatePartition("objectid")
	require.NoError(t, err)

	assignedDir, err := manager.GetDirectory("objectid")
	require.NoError(t, err)
	require.Equal(t, assignedDir, createdDir)
	require.DirExists(t, createdDir)
}
