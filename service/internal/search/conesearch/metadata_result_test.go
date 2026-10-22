package conesearch

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestResultFromSelectionMetadata_MatchCounts(t *testing.T) {
	t.Run("groups report the global total and their catalog count", func(t *testing.T) {
		selection := selectionResult[repository.Metadata]{
			Data: []repository.Metadata{
				{ID: "a1", Catalog: "alpha"},
				{ID: "b1", Catalog: "beta"},
				{ID: "b2", Catalog: "beta"},
			},
			Distance:      []float64{1, 2, 3},
			Total:         5,
			CatalogCounts: map[string]int{"alpha": 1, "beta": 4},
		}

		result := ResultFromSelectionMetadata(selection)

		require.Len(t, result, 2)
		byCatalog := make(map[string]MetadataResult, len(result))
		for _, group := range result {
			require.Equal(t, selection.Total, group.Total)
			require.LessOrEqual(t, len(group.Data), group.TotalInCatalog)
			byCatalog[group.Catalog] = group
		}

		require.Equal(t, 1, byCatalog["alpha"].TotalInCatalog)
		require.Equal(t, 4, byCatalog["beta"].TotalInCatalog)

		sum := 0
		for _, group := range result {
			sum += group.TotalInCatalog
		}
		require.Equal(t, selection.Total, sum)
	})

	t.Run("a catalog fully cut by nneighbor produces no group", func(t *testing.T) {
		selection := selectionResult[repository.Metadata]{
			Data: []repository.Metadata{
				{ID: "a1", Catalog: "alpha"},
			},
			Distance:      []float64{1},
			Total:         3,
			CatalogCounts: map[string]int{"alpha": 1, "gamma": 2},
		}

		result := ResultFromSelectionMetadata(selection)

		require.Len(t, result, 1)
		require.Equal(t, "alpha", result[0].Catalog)
		require.Equal(t, selection.Total, result[0].Total)
		require.Equal(t, 1, result[0].TotalInCatalog)
	})
}
