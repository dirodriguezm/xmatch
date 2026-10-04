// Copyright 2024-2026 Diego Rodriguez Mancini
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

package conesearch

import (
	"sort"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// selectionResult is the output of selectNearest: the selected objects and
// their great-circle distances in arcseconds, ordered nearest-first. Data and
// Distance are index-aligned.
//
// Total and CatalogCounts describe the match counts over the radius-filtered
// candidate set before the nneighbor truncation: Total is the number of
// matching objects across all catalogs, and CatalogCounts maps each catalog to
// the number of its own matching objects.
type selectionResult[T any] struct {
	Data          []T
	Distance      []float64
	Total         int
	CatalogCounts map[string]int
}

// coordinateAccessor returns an object's sky position in degrees.
type coordinateAccessor[T any] func(T) (ra, dec float64)

// catalogAccessor returns the catalog an object belongs to.
type catalogAccessor[T any] func(T) string

// selectNearest returns the objects whose great-circle distance from
// (ra, dec) is at most radius, ordered nearest-first and truncated to at most
// nneighbor objects. Selection is independent of the input order. The result
// also carries the total number of in-radius matches and the per-catalog
// tally, both counted before the truncation.
func selectNearest[T any](
	objects []T,
	ra, dec, radius float64,
	nneighbor int,
	coordinates coordinateAccessor[T],
	catalogs catalogAccessor[T],
) selectionResult[T] {
	candidates := make([]neighbor[T], 0, len(objects))
	catalogCounts := make(map[string]int)
	for _, obj := range objects {
		objRa, objDec := coordinates(obj)
		distance := haversineDistance(ra, dec, objRa, objDec)
		if distance > radius {
			continue
		}
		candidates = append(candidates, neighbor[T]{object: obj, distance: distance})
		catalogCounts[catalogs(obj)]++
	}

	total := len(candidates)

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].distance < candidates[j].distance
	})

	if len(candidates) > nneighbor {
		candidates = candidates[:nneighbor]
	}

	result := selectionResult[T]{
		Data:          make([]T, 0, len(candidates)),
		Distance:      make([]float64, 0, len(candidates)),
		Total:         total,
		CatalogCounts: catalogCounts,
	}
	for _, candidate := range candidates {
		result.Data = append(result.Data, candidate.object)
		result.Distance = append(result.Distance, candidate.distance)
	}
	return result
}

type neighbor[T any] struct {
	object   T
	distance float64
}

func mastercatCoordinates(obj repository.Mastercat) (float64, float64) {
	return obj.Ra, obj.Dec
}

func metadataCoordinates(obj repository.Metadata) (float64, float64) {
	return obj.Ra, obj.Dec
}

func mastercatCatalog(obj repository.Mastercat) string {
	return obj.Cat
}

func metadataCatalog(obj repository.Metadata) string {
	return obj.Catalog
}
