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
	Obj        repository.Mastercat
	AllwiseObj repository.GetAllwiseFromPixelsRow
	catalog    string
}

func (knn knnObject) Dimensions() int {
	return 2
}

func (knn knnObject) Dimension(i int) float64 {
	var dimensions []float64
	switch strings.ToLower(knn.catalog) {
	case "allwise":
		dimensions = []float64{knn.AllwiseObj.Ra, knn.AllwiseObj.Dec}
	default:
		dimensions = []float64{knn.Obj.Ra, knn.Obj.Dec}
	}
	return dimensions[i]
}

func NearestNeighborSearch(objects []repository.Mastercat, ra, dec, radius float64, maxNeighbors int) []repository.Mastercat {
	pts := []kdtree.Point{}
	for _, obj := range objects {
		pts = append(pts, knnObject{Obj: obj})
	}
	tree := kdtree.New(pts)

	nearObjs := tree.KNN(&points.Point2D{X: ra, Y: dec}, maxNeighbors)

	// now we need to check that distance between nearest objects and center is actually lower than radius
	result := []repository.Mastercat{}
	for _, obj := range nearObjs {
		dist := haversineDistance(obj, &points.Point2D{X: ra, Y: dec})
		if dist > radius {
			continue
		}
		result = append(result, obj.(knnObject).Obj)
	}
	return result
}

func NearestNeighborSearchForAllwiseMetadata(objects []repository.GetAllwiseFromPixelsRow, ra, dec, radius float64, maxNeighbors int) []repository.AllwiseMetadata {
	pts := []kdtree.Point{}
	for _, obj := range objects {
		pts = append(pts, knnObject{AllwiseObj: obj, catalog: "allwise"})
	}
	tree := kdtree.New(pts)

	nearObjs := tree.KNN(&points.Point2D{X: ra, Y: dec}, maxNeighbors)

	// now we need to check that distance between nearest objects and center is actually lower than radius
	result := []repository.AllwiseMetadata{}
	for _, obj := range nearObjs {
		dist := haversineDistance(obj, &points.Point2D{X: ra, Y: dec})
		if dist > radius {
			continue
		}
		result = append(result, obj.(knnObject).AllwiseObj.ToAllwiseMetadata())
	}
	return result
}

func euclideanDistance(p1, p2 kdtree.Point) float64 {
	if p1.Dimensions() != 2 || p2.Dimensions() != 2 {
		err := fmt.Errorf("Can't calculate distance between points of dimension %v and %v", p1.Dimensions(), p2.Dimensions())
		panic(err)
	}
	dsquared := math.Pow(p2.Dimension(0)-p1.Dimension(0), 2) + math.Pow(p2.Dimension(1)-p1.Dimension(1), 2)
	return math.Sqrt(dsquared)
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
