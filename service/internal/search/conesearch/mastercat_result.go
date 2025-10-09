package conesearch

import (
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/knn"
)

type MastercatExtended struct {
	repository.Mastercat
	Distance float64 `json:"distance"`
}

type MastercatResult struct {
	Catalog string              `json:"catalog"`
	Data    []MastercatExtended `json:"data"`
}

func ResultFromKnn(objs knn.KnnResult[repository.Mastercat]) []MastercatResult {
	result := make([]MastercatResult, 0)
	grouped := make(map[string][]MastercatExtended)
	for i, m := range objs.Data {
		grouped[m.Cat] = append(grouped[m.Cat], MastercatExtended{
			Mastercat: m,
			Distance:  objs.Distance[i],
		})
	}
	for catalog, data := range grouped {
		result = append(result, MastercatResult{Catalog: catalog, Data: data})
	}
	return result
}
