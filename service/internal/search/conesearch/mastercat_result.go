package conesearch

import "github.com/dirodriguezm/xmatch/service/internal/repository"

type MastercatResult struct {
	Catalog string                 `json:"catalog"`
	Data    []repository.Mastercat `json:"data"`
}

func ResultFromMastercat(objs []repository.Mastercat) []MastercatResult {
	result := make([]MastercatResult, 0)
	grouped := make(map[string][]repository.Mastercat)
	for _, m := range objs {
		grouped[m.Cat] = append(grouped[m.Cat], m)
	}
	for catalog, data := range grouped {
		result = append(result, MastercatResult{Catalog: catalog, Data: data})
	}
	return result
}
