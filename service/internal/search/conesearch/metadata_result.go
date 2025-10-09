package conesearch

import (
	"encoding/json"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/knn"
)

type MetadataExtended struct {
	repository.Metadata `json:"-"`
	Distance            float64 `json:"distance"`
}

func (m MetadataExtended) MarshalJSON() ([]byte, error) {
	metadataBytes, err := json.Marshal(m.Metadata)
	if err != nil {
		return nil, err
	}

	var metadataMap map[string]any
	if err := json.Unmarshal(metadataBytes, &metadataMap); err != nil {
		return nil, err
	}

	metadataMap["distance"] = m.Distance
	return json.Marshal(metadataMap)
}

type MetadataResult struct {
	Catalog string             `json:"catalog"`
	Data    []MetadataExtended `json:"data"`
}

func ResultFromKnnMetadata(metadata knn.KnnResult[repository.Metadata]) []MetadataResult {
	result := make([]MetadataResult, 0)
	grouped := make(map[string][]MetadataExtended)
	for i, m := range metadata.Data {
		grouped[m.GetCatalog()] = append(grouped[m.GetCatalog()], MetadataExtended{
			Metadata: m,
			Distance: metadata.Distance[i],
		})
	}
	for catalog, data := range grouped {
		result = append(result, MetadataResult{Catalog: catalog, Data: data})
	}
	return result
}
