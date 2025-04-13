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
	"math"
	"strings"
)

type Partition struct {
	Levels []int
}

func (p *Partition) LevelsToString(maxPartitions int) (string, error) {
	if maxPartitions < 1 {
		return "", fmt.Errorf("Could not parse levels to string with negative partitions: %d", maxPartitions)
	}

	result := make([]string, len(p.Levels))
	for i := 0; i < len(p.Levels); i++ {
		if p.Levels[i] > maxPartitions {
			return "", fmt.Errorf("Can't parse levels to string if the number of partitions is lower than level %d", p.Levels[i])
		}

		if p.Levels[i] < 0 {
			return "", fmt.Errorf("Can't parse negative level to string %d", p.Levels[i])
		}

		format := fmt.Sprintf("%%0%dd", getNumberWidth(maxPartitions))
		result[i] = fmt.Sprintf(format, p.Levels[i])
	}

	return strings.Join(result, "/"), nil
}

func getNumberWidth(n int) int {
	return int(math.Ceil(math.Log10(float64(n))))
}
