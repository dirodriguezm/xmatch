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
	"fmt"
	"os"
	"path"

	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/partition"
)

type FileSystemManager struct {
	handler partition.PartitionHandler
	baseDir string
}

func (manager FileSystemManager) CreatePartition(oid string) (string, error) {
	part, err := manager.assignPartition(oid)
	if err != nil {
		return "", err
	}
	return manager.createPartitionDirectory(part)
}

func (manager FileSystemManager) createPartitionDirectory(part partition.Partition) (string, error) {
	partitionDir, err := part.LevelsToString(manager.handler.NumPartitions)
	if err != nil {
		return "", fmt.Errorf("Could not get partition directory for partition %v.\n%w", part, err)
	}

	dir := path.Join(manager.baseDir, partitionDir)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("Could not create directory %s.\n%w", dir, err)
	}
	return dir, nil
}

func (manager FileSystemManager) assignPartition(oid string) (partition.Partition, error) {
	part, err := manager.handler.GetPartition(oid)
	if err != nil {
		err = fmt.Errorf("Could not assign partition to oid %s\n%w", oid, err)
	}
	return part, err
}

func (manager FileSystemManager) GetDirectory(oid string) (string, error) {
	// get the partition for the object id
	part, err := manager.assignPartition(oid)
	if err != nil {
		return "", err
	}

	// get the nested directory for the assigned partition
	levels, err := part.LevelsToString(manager.handler.NumPartitions)
	if err != nil {
		return "", err
	}

	// join the base dir with the nested directory
	dir := path.Join(manager.baseDir, levels)
	return dir, nil
}
