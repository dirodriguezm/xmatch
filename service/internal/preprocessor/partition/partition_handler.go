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

package partition

import (
	"fmt"

	"github.com/spaolacci/murmur3"
)

type PartitionHandler struct {
	NumPartitions   int
	PartitionLevels int
}

func (handler *PartitionHandler) GetPartition(oid string) (Partition, error) {
	if handler.PartitionLevels <= 0 {
		return Partition{}, fmt.Errorf("Can't get partition for negative levels %d", handler.PartitionLevels)
	}
	if handler.NumPartitions <= 0 {
		return Partition{}, fmt.Errorf("Can't get partition for negative number of partitions %d", handler.NumPartitions)
	}

	partition := Partition{}

	hasher := murmur3.New64()
	hasher.Write([]byte(oid))
	hashValue := hasher.Sum64()

	bitsPerLevel := 64 / handler.PartitionLevels
	remainingBits := 64 % handler.PartitionLevels

	levels := make([]int, handler.PartitionLevels)
	currentHash := hashValue
	for level := 0; level < handler.PartitionLevels; level++ {
		// Add an extra bit to early levels if we have remaining bits
		levelBits := bitsPerLevel
		if level < remainingBits {
			levelBits++
		}

		// Create mask for this level's bits
		mask := uint64((1 << levelBits) - 1)

		// Extract this level's partition
		levelValue := currentHash & mask
		levels[level] = int(levelValue % uint64(handler.NumPartitions))

		// Shift right by bits used
		currentHash = currentHash >> levelBits
	}

	partition.Levels = levels

	return partition, nil
}
