package conesearch

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// selectionBenchmarkObjects generates the same candidate distribution as the
// temporary knn baseline benchmark (seed 42, center 179.5/14.5, ±0.5° spread,
// every candidate inside the 3600″ radius) so the two runs are directly
// comparable.
func selectionBenchmarkObjects(n int) []repository.Mastercat {
	rng := rand.New(rand.NewSource(42))
	objects := make([]repository.Mastercat, n)
	for i := range objects {
		objects[i] = repository.Mastercat{
			ID:   strconv.Itoa(i),
			Ipix: 1,
			Ra:   179.5 + (rng.Float64() - 0.5),
			Dec:  14.5 + (rng.Float64() - 0.5),
			Cat:  "bench",
		}
	}
	return objects
}

// BenchmarkSelectNearest measures the linear filter-and-select helper over the
// same candidate/N matrix as BenchmarkNearestNeighborSearch.
func BenchmarkSelectNearest(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 10000, 100000} {
		objects := selectionBenchmarkObjects(count)
		for _, nneighbor := range []int{1, 10} {
			b.Run(fmt.Sprintf("candidates=%d/nneighbor=%d", count, nneighbor), func(b *testing.B) {
				for b.Loop() {
					_ = selectNearest(objects, 179.5, 14.5, 3600, nneighbor, mastercatCoordinates)
				}
			})
		}
	}
}
