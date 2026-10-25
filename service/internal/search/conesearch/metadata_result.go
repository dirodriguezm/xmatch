package conesearch

import (
	"encoding/json"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type MetadataExtended struct {
	repository.Metadata `json:"-"`
	Distance            float64 `json:"distance"`
}

func (m MetadataExtended) MarshalJSON() ([]byte, error) {
	payload := m.Object
	if payload == nil {
		payload = m.Metadata
	}
	metadataBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var metadataMap map[string]any
	if err := json.Unmarshal(metadataBytes, &metadataMap); err != nil {
		return nil, err
	}

	// The catalog object does not carry the coordinates used for selection:
	// allwise/gaia have no ra/dec columns, and erosita's own nullable pair is
	// not necessarily the mastercat pair. Publish the envelope pair so every
	// catalog reports the same coordinates the cone search was built from.
	metadataMap["ra"] = m.Ra
	metadataMap["dec"] = m.Dec
	metadataMap["distance"] = m.Distance
	return json.Marshal(metadataMap)
}

type MetadataResult struct {
	Catalog        string             `json:"catalog"`
	Data           []MetadataExtended `json:"data"`
	Total          int                `json:"total"`
	TotalInCatalog int                `json:"total_in_catalog"`
}

func ResultFromSelectionMetadata(metadata selectionResult[repository.Metadata]) []MetadataResult {
	result := make([]MetadataResult, 0)
	grouped := make(map[string][]MetadataExtended)
	for i, m := range metadata.Data {
		grouped[m.Catalog] = append(grouped[m.Catalog], MetadataExtended{
			Metadata: m,
			Distance: metadata.Distance[i],
		})
	}
	for catalog, data := range grouped {
		result = append(result, MetadataResult{
			Catalog:        catalog,
			Data:           data,
			Total:          metadata.Total,
			TotalInCatalog: metadata.CatalogCounts[catalog],
		})
	}
	return result
}
