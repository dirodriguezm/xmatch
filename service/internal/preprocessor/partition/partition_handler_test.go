package partition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPartitionHandler_GetPartition_SimpleCase(t *testing.T) {
	expectedPartition := Partition{Levels: []int{0, 0, 0}}
	handler := PartitionHandler{NumPartitions: 1, PartitionLevels: 3}

	partition, err := handler.GetPartition("objectid")

	require.NoError(t, err)
	require.Equal(t, expectedPartition, partition)
}

func TestPartitionHandler_GetPartition_Cases(t *testing.T) {
	for i := 1; i <= 256; i++ {
		for j := 1; j <= 3; j++ {
			handler := PartitionHandler{NumPartitions: i, PartitionLevels: j}
			partition, err := handler.GetPartition("objectid")

			require.NoError(t, err)
			require.Len(t, partition.Levels, j, "Case NumPartitions %d, Levels %d", i, j)

			for level := range partition.Levels {
				require.Less(t, partition.Levels[level], i)
			}
		}
	}
}

func TestPartitionHandler_GetPartition_OidHash(t *testing.T) {
	handler := PartitionHandler{NumPartitions: 2, PartitionLevels: 2}
	partition1, err := handler.GetPartition("FirstId")
	require.NoError(t, err)
	partition2, err := handler.GetPartition("SecondId")
	require.NoError(t, err)

	part1Str, _ := partition1.LevelsToString(handler.PartitionLevels)
	part2Str, _ := partition2.LevelsToString(handler.PartitionLevels)
	require.NotEqual(t, part1Str, part2Str)

	partition3, err := handler.GetPartition("FirstId")
	require.NoError(t, err)
	part3Str, _ := partition3.LevelsToString(handler.PartitionLevels)
	require.Equal(t, part1Str, part3Str)
}
