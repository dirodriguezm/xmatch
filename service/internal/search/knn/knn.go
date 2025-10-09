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

package knn

import (
	"fmt"
	"math"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type knnObject struct {
	Obj         repository.Mastercat
	MetadataObj repository.MetadataWithCoordinates
	catalog     string
}

func (knn knnObject) Dimensions() int {
	return 2
}

func (knn knnObject) Dimension(i int) float64 {
	var dimensions []float64
	switch strings.ToLower(knn.catalog) {
	case "":
		dimensions = []float64{knn.Obj.Ra, knn.Obj.Dec}
	default:
		ra, dec := knn.MetadataObj.GetCoordinates()
		dimensions = []float64{ra, dec}
	}
	return dimensions[i]
}

type KnnResult[T any] struct {
	Data     []T
	Distance []float64
}

func NearestNeighborSearch(
	objects []repository.Mastercat,
	ra, dec, radius float64,
	maxNeighbors int,
) KnnResult[repository.Mastercat] {
	pts := []kdtree.Point{}
	for _, obj := range objects {
		pts = append(pts, knnObject{Obj: obj})
	}
	tree := kdtree.New(pts)

	nearObjs := tree.KNN(&points.Point2D{X: ra, Y: dec}, maxNeighbors)

	// now we need to check that distance between nearest objects and center is actually lower than radius
	result := KnnResult[repository.Mastercat]{}
	for _, obj := range nearObjs {
		dist := haversineDistance(obj, &points.Point2D{X: ra, Y: dec})
		if dist > radius {
			continue
		}
		result.Data = append(result.Data, obj.(knnObject).Obj)
		result.Distance = append(result.Distance, dist)
	}
	return result
}

func NearestNeighborSearchForMetadata(
	objects []repository.MetadataWithCoordinates,
	ra, dec, radius float64,
	maxNeighbors int,
	catalog string,
) KnnResult[repository.Metadata] {
	pts := []kdtree.Point{}
	for _, obj := range objects {
		pts = append(pts, knnObject{MetadataObj: obj, catalog: catalog})
	}
	tree := kdtree.New(pts)

	nearObjs := tree.KNN(&points.Point2D{X: ra, Y: dec}, maxNeighbors)

	// now we need to check that distance between nearest objects and center is actually lower than radius
	result := KnnResult[repository.Metadata]{}
	for _, obj := range nearObjs {
		dist := haversineDistance(obj, &points.Point2D{X: ra, Y: dec})
		if dist > radius {
			continue
		}
		result.Distance = append(result.Distance, dist)
		switch strings.ToLower(obj.(knnObject).catalog) {
		case "allwise":
			result.Data = append(result.Data, convertToAllwise(obj.(knnObject).MetadataObj.(repository.GetAllwiseFromPixelsRow)))
		case "gaia":
			result.Data = append(result.Data, convertToGaia(obj.(knnObject).MetadataObj.(repository.GetGaiaFromPixelsRow)))
		case "all":
			if metadataObj, ok := obj.(knnObject).MetadataObj.(repository.GetAllwiseFromPixelsRow); ok {
				result.Data = append(result.Data, convertToAllwise(metadataObj))
			} else if metadataObj, ok := obj.(knnObject).MetadataObj.(repository.GetGaiaFromPixelsRow); ok {
				result.Data = append(result.Data, convertToGaia(metadataObj))
			}
		default:
			panic("Unknown catalog to KNN Search for Metadata")
		}
	}
	return result
}

func convertToAllwise(obj repository.GetAllwiseFromPixelsRow) repository.Allwise {
	return repository.Allwise{
		ID:         obj.ID,
		W1mpro:     obj.W1mpro,
		W1sigmpro:  obj.W1sigmpro,
		W2mpro:     obj.W2mpro,
		W2sigmpro:  obj.W2sigmpro,
		W3mpro:     obj.W3mpro,
		W3sigmpro:  obj.W3sigmpro,
		W4mpro:     obj.W4mpro,
		W4sigmpro:  obj.W4sigmpro,
		JM2mass:    obj.JM2mass,
		JMsig2mass: obj.JMsig2mass,
		HM2mass:    obj.HM2mass,
		HMsig2mass: obj.HMsig2mass,
		KM2mass:    obj.KM2mass,
		KMsig2mass: obj.KMsig2mass,
	}
}

func convertToGaia(obj repository.GetGaiaFromPixelsRow) repository.Gaia {
	return repository.Gaia{
		ID:                  obj.ID,
		PhotGMeanFlux:       obj.PhotGMeanFlux,
		PhotGMeanFluxError:  obj.PhotGMeanFluxError,
		PhotGMeanMag:        obj.PhotGMeanMag,
		PhotBpMeanFlux:      obj.PhotBpMeanFlux,
		PhotBpMeanFluxError: obj.PhotBpMeanFluxError,
		PhotBpMeanMag:       obj.PhotBpMeanMag,
		PhotRpMeanFlux:      obj.PhotRpMeanFlux,
		PhotRpMeanFluxError: obj.PhotRpMeanFluxError,
		PhotRpMeanMag:       obj.PhotRpMeanMag,
	}
}

// Return the distance in arcsec, between two points in a sphere, using the Haversine Formula
func haversineDistance(p1, p2 kdtree.Point) float64 {
	if p1.Dimensions() != 2 || p2.Dimensions() != 2 {
		err := fmt.Errorf("Can't calculate distance between points of dimension %v and %v", p1.Dimensions(), p2.Dimensions())
		panic(err)
	}

	ra1 := p1.Dimension(0)
	dec1 := p1.Dimension(1)
	ra2 := p2.Dimension(0)
	dec2 := p2.Dimension(1)

	ra1Rad := ra1 * math.Pi / 180.0
	dec1Rad := dec1 * math.Pi / 180.0
	ra2Rad := ra2 * math.Pi / 180.0
	dec2Rad := dec2 * math.Pi / 180.0

	deltaRA := ra2Rad - ra1Rad
	deltaDec := dec2Rad - dec1Rad

	a := math.Sin(deltaDec/2)*math.Sin(deltaDec/2) +
		math.Cos(dec1Rad)*math.Cos(dec2Rad)*
			math.Sin(deltaRA/2)*math.Sin(deltaRA/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * 180.0 / math.Pi * 3600.0
}
