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
		Handler: handler,
		BaseDir: dir,
	}

	_, err := manager.createPartitionDirectory(part)
	require.NoError(t, err)

	require.DirExists(t, path.Join(dir, "00"))
	require.DirExists(t, path.Join(dir, "00", "00"))

	dir = t.TempDir()
	manager.BaseDir = dir
	part = partition.Partition{Levels: []int{1, 0}}

	_, err = manager.createPartitionDirectory(part)
	require.NoError(t, err)

	require.DirExists(t, path.Join(dir, "01"))
	require.DirExists(t, path.Join(dir, "01", "00"))

	dir = t.TempDir()
	manager.BaseDir = dir
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
		Handler: handler,
		BaseDir: dir,
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
