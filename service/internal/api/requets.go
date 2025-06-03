package api

type BulkConesearchRequest struct {
	Ra        []float64 `json:"ra"`
	Dec       []float64 `json:"dec"`
	Radius    float64   `json:"radius"`
	Catalog   string    `json:"catalog"`
	Nneighbor int       `json:"nneighbor"`
}

type BulkMetadataRequest struct {
	Ids     []string `json:"ids"`
	Catalog string   `json:"catalog"`
}
