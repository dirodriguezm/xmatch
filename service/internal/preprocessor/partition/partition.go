package partition

import (
	"fmt"
	"math"
	"strings"
)

type Partition struct {
	Levels []int
}

func (p *Partition) LevelsToString(maxPartitions int) string {
	result := make([]string, len(p.Levels))
	for i := 0; i < len(p.Levels); i++ {
		format := fmt.Sprintf("%%0%dd", getNumberWidth(maxPartitions))
		result[i] = fmt.Sprintf(format, p.Levels[i])
	}
	return strings.Join(result, "/")
}

func getNumberWidth(n int) int {
	return int(math.Ceil(math.Log10(float64(n))))
}
