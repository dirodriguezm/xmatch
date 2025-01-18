package partition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ToStringTestCase struct {
	Input         []int
	Expected      string
	NumPartitions int
	Error         bool
}

func TestPartition_ToString(t *testing.T) {
	cases := []ToStringTestCase{
		{[]int{10, 88}, "010/088", 256, false},
		{[]int{10, 50}, "10/50", 64, false},
		{[]int{10, 88}, "0010/0088", 1024, false},
		{[]int{0, 24}, "00/24", 32, false},
		{[]int{0, 0}, "00/00", 32, false},
		{[]int{0, 1}, "0/1", 8, false},
		{[]int{-1, -1}, "", 2, true},
		{[]int{1, 1}, "", -1, true},
		{[]int{10, 10}, "", 1, true},
	}

	for i := 0; i < len(cases); i++ {
		part := Partition{Levels: cases[i].Input}
		result, err := part.LevelsToString(cases[i].NumPartitions)

		if cases[i].Error {
			require.Error(t, err, "Case %d", i)
		} else {
			require.NoError(t, err, "Case %d", i)
		}
		require.Equal(t, cases[i].Expected, result, "Case %d. %v", i, cases[i])
	}
}
