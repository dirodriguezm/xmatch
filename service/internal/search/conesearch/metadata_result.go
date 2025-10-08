package conesearch

import "github.com/dirodriguezm/xmatch/service/internal/repository"

type MetadataResult struct {
	Catalog string                `json:"catalog"`
	Data    []repository.Metadata `json:"data"`
}

func ResultFromMetadata(metadata []repository.Metadata) []MetadataResult {
	result := make([]MetadataResult, 0)
	grouped := make(map[string][]repository.Metadata)
	for _, m := range metadata {
		grouped[m.GetCatalog()] = append(grouped[m.GetCatalog()], m)
	}
	for catalog, data := range grouped {
		result = append(result, MetadataResult{Catalog: catalog, Data: data})
	}
	return result
}
