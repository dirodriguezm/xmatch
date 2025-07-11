// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
