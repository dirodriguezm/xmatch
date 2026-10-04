package conesearch

import (
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func objectIDs(objects []repository.Mastercat) []string {
	ids := make([]string, 0, len(objects))
	for _, obj := range objects {
		ids = append(ids, obj.ID)
	}
	return ids
}

func TestSelectNearest(t *testing.T) {
	t.Run("more in-radius objects than requested neighbors returns the closest", func(t *testing.T) {
		objects := make([]repository.Mastercat, 0, 10)
		for i := range 10 {
			objects = append(objects, repository.Mastercat{
				ID:  fmt.Sprintf("obj-%02d", i),
				Ra:  float64(i) * 0.1,
				Dec: 0,
				Cat: "test",
			})
		}

		result := selectNearest(objects, 0, 0, 3600, 3, mastercatCoordinates, mastercatCatalog)

		require.Len(t, result.Data, 3)
		require.Equal(t, []string{"obj-00", "obj-01", "obj-02"}, objectIDs(result.Data))
		require.Len(t, result.Distance, 3)
		for i := 1; i < len(result.Distance); i++ {
			require.GreaterOrEqual(t, result.Distance[i], result.Distance[i-1])
		}

		require.Equal(t, 10, result.Total)
		require.Equal(t, map[string]int{"test": 10}, result.CatalogCounts)
	})

	t.Run("fewer in-radius objects than requested neighbors returns all of them", func(t *testing.T) {
		objects := []repository.Mastercat{
			{ID: "inside-near", Ra: 0.1, Dec: 0, Cat: "test"},
			{ID: "outside", Ra: 2, Dec: 0, Cat: "test"},
			{ID: "inside-far", Ra: 0.5, Dec: 0, Cat: "test"},
		}

		result := selectNearest(objects, 0, 0, 3600, 5, mastercatCoordinates, mastercatCatalog)

		require.Equal(t, []string{"inside-near", "inside-far"}, objectIDs(result.Data))
		require.Len(t, result.Distance, 2)

		require.Equal(t, 2, result.Total)
		require.Equal(t, map[string]int{"test": 2}, result.CatalogCounts)
	})

	t.Run("results are ordered nearest-first regardless of input order", func(t *testing.T) {
		objects := []repository.Mastercat{
			{ID: "far", Ra: 0.3, Dec: 0, Cat: "test"},
			{ID: "near", Ra: 0.1, Dec: 0, Cat: "test"},
			{ID: "middle", Ra: 0.2, Dec: 0, Cat: "test"},
		}

		result := selectNearest(objects, 0, 0, 3600, 2, mastercatCoordinates, mastercatCatalog)

		require.Equal(t, []string{"near", "middle"}, objectIDs(result.Data))
		require.Len(t, result.Distance, 2)
		for i := 1; i < len(result.Distance); i++ {
			require.GreaterOrEqual(t, result.Distance[i], result.Distance[i-1])
		}

		require.Equal(t, 3, result.Total)
	})

	t.Run("no candidates in radius returns an empty result", func(t *testing.T) {
		objects := []repository.Mastercat{
			{ID: "outside", Ra: 2, Dec: 0, Cat: "test"},
		}

		result := selectNearest(objects, 0, 0, 3600, 5, mastercatCoordinates, mastercatCatalog)

		require.Empty(t, result.Data)
		require.Empty(t, result.Distance)
		require.Zero(t, result.Total)
		require.Empty(t, result.CatalogCounts)
	})

	t.Run("counts are per catalog and exclude out-of-radius candidates", func(t *testing.T) {
		objects := []repository.Mastercat{
			{ID: "alpha-near", Ra: 0.1, Dec: 0, Cat: "alpha"},
			{ID: "alpha-outside", Ra: 2, Dec: 0, Cat: "alpha"},
			{ID: "beta-near", Ra: 0.2, Dec: 0, Cat: "beta"},
			{ID: "beta-farther", Ra: 0.3, Dec: 0, Cat: "beta"},
		}

		result := selectNearest(objects, 0, 0, 3600, 1, mastercatCoordinates, mastercatCatalog)

		require.Equal(t, []string{"alpha-near"}, objectIDs(result.Data))
		require.Equal(t, 3, result.Total)
		require.Equal(t, map[string]int{"alpha": 1, "beta": 2}, result.CatalogCounts)
	})
}

func TestSelectNearest_RAWrapAround(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "farther", Ra: 0.5, Dec: 0, Cat: "test"},
		{ID: "wrap", Ra: 359.95, Dec: 0, Cat: "test"},
	}

	result := selectNearest(objects, 0.05, 0, 3600, 1, mastercatCoordinates, mastercatCatalog)

	require.Equal(t, []string{"wrap"}, objectIDs(result.Data))
	require.InDelta(t, 0.1*3600, result.Distance[0], 1e-6)
	require.Equal(t, 2, result.Total)
}

func TestSelectNearest_HighDeclination(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "across-pole", Ra: 180, Dec: 89.95, Cat: "test"},
		{ID: "right-angle", Ra: 90, Dec: 89.8, Cat: "test"},
		{ID: "near-pole", Ra: 0, Dec: 89.95, Cat: "test"},
	}

	result := selectNearest(objects, 0, 89.9, 3600, 2, mastercatCoordinates, mastercatCatalog)

	require.Equal(t, []string{"near-pole", "across-pole"}, objectIDs(result.Data))
	require.InDelta(t, 0.05*3600, result.Distance[0], 1e-6)
	require.InDelta(t, 0.15*3600, result.Distance[1], 1e-6)
	require.Equal(t, 3, result.Total)
}
