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
			partition, err := handler.GetPartition("0438p015_ac51-018218")

			require.NoError(t, err)
			require.Len(t, partition.Levels, j, "Case NumPartitions %d, Levels %d", i, j)

			for level := range partition.Levels {
				require.Less(t, partition.Levels[level], i)
				require.GreaterOrEqual(t, partition.Levels[level], 0)
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

func TestPartitionHandler_GetPartition_SameOid(t *testing.T) {
	handler := PartitionHandler{NumPartitions: 2, PartitionLevels: 2}
	partition1, err := handler.GetPartition("FirstId")
	require.NoError(t, err)

	for i := 0; i < 1000; i++ {
		partition2, err := handler.GetPartition("FirstId")
		require.NoError(t, err)
		require.Equal(t, partition1, partition2)
	}
}
