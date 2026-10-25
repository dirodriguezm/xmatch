package conesearch

import (
	"encoding/json"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

// TestMetadataExtended_MarshalJSON_PublishesMastercatCoordinates pins the
// conesearch metadata payload: every catalog reports the envelope's mastercat
// ra/dec, even when the catalog table has no coordinates (allwise/gaia) or
// carries its own nullable pair (erosita).
func TestMetadataExtended_MarshalJSON_PublishesMastercatCoordinates(t *testing.T) {
	testCases := map[string]any{
		"allwise has no table coordinates": repository.Allwise{ID: "allwise-1"},
		"gaia has no table coordinates":    repository.Gaia{ID: "gaia-1"},
		"erosita table coordinates are overridden": repository.Erosita{
			ID:  "erosita-1",
			Ra:  new(100.0),
			Dec: new(50.0),
		},
	}

	for name, object := range testCases {
		t.Run(name, func(t *testing.T) {
			payload, err := json.Marshal(MetadataExtended{
				Metadata: repository.Metadata{
					ID:      "object-1",
					Catalog: "catalog",
					Ra:      12.5,
					Dec:     -3.25,
					Object:  object,
				},
				Distance: 1.25,
			})
			require.NoError(t, err)

			var decoded map[string]any
			require.NoError(t, json.Unmarshal(payload, &decoded))
			require.Equal(t, 12.5, decoded["ra"])
			require.Equal(t, -3.25, decoded["dec"])
			require.Equal(t, 1.25, decoded["distance"])
			require.NotContains(t, decoded, "ra_2")
			require.NotContains(t, decoded, "dec_2")
		})
	}
}

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
