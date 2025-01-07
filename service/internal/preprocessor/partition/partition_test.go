package partition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ToStringTestCase struct {
	Input         []int
	Expected      string
	NumPartitions int
}

func TestPartition_ToString(t *testing.T) {
	cases := []ToStringTestCase{
		{[]int{10, 88}, "010/088", 256},
		{[]int{10, 88}, "10/88", 64},
		{[]int{10, 88}, "0010/0088", 1024},
		{[]int{0, 24}, "00/24", 32},
		{[]int{0, 0}, "00/00", 32},
		{[]int{0, 1}, "0/1", 8},
	}

	for i := 0; i < len(cases); i++ {
		part := Partition{Levels: cases[i].Input}
		result := part.LevelsToString(cases[i].NumPartitions)

		require.Equal(t, cases[i].Expected, result, "Case %d. %v", i, cases[i])
	}
}
