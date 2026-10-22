package conesearch

import (
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type MastercatExtended struct {
	repository.Mastercat
	Distance float64 `json:"distance"`
}

type MastercatResult struct {
	Catalog        string              `json:"catalog"`
	Data           []MastercatExtended `json:"data"`
	Index          int                 `json:"index"`
	Total          int                 `json:"total"`
	TotalInCatalog int                 `json:"total_in_catalog"`
}

func ResultFromSelection(objs selectionResult[repository.Mastercat], index int) []MastercatResult {
	result := make([]MastercatResult, 0)
	grouped := make(map[string][]MastercatExtended)
	for i, m := range objs.Data {
		grouped[m.Cat] = append(grouped[m.Cat], MastercatExtended{
			Mastercat: m,
			Distance:  objs.Distance[i],
		})
	}
	for catalog, data := range grouped {
		result = append(result, MastercatResult{
			Catalog:        catalog,
			Data:           data,
			Index:          index,
			Total:          objs.Total,
			TotalInCatalog: objs.CatalogCounts[catalog],
		})
	}
	return result
}
