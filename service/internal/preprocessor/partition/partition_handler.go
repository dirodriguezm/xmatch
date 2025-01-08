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
		levels[level] = int(levelValue) % handler.NumPartitions

		// Shift right by bits used
		currentHash = currentHash >> levelBits
	}

	partition.Levels = levels

	return partition, nil
}
