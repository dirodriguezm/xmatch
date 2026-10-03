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

// BulkConesearchRequest represents a bulk cone search query with multiple coordinates
type BulkConesearchRequest struct {
	// Right ascension values in degrees (0-360)
	Ra []float64 `json:"ra" example:"180.5,181.2"`
	// Declination values in degrees (-90 to 90)
	Dec []float64 `json:"dec" example:"-45.0,-45.5"`
	// Search radius in degrees
	Radius float64 `json:"radius" example:"0.01"`
	// Catalog name to search in (default: all)
	Catalog string `json:"catalog" example:"allwise"`
	// Number of neighbors to return per coordinate (default: 1)
	Nneighbor int `json:"nneighbor" example:"1"`
}

// BulkMetadataRequest represents a bulk metadata query with multiple object IDs
type BulkMetadataRequest struct {
	// List of object identifiers to search for
	Ids []string `json:"ids" example:"id1,id2"`
	// Catalog to search in
	Catalog string `json:"catalog" example:"allwise"`
}
