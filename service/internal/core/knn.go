package core

import (
	"fmt"
	"math"
	"xmatch/service/pkg/repository"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type KNNObject struct {
	Obj repository.Mastercat
}

func (knn KNNObject) Dimensions() int {
	return 2
}

func (knn KNNObject) Dimension(i int) float64 {
	dimensions := []float64{knn.Obj.Ra, knn.Obj.Dec}
	return dimensions[i]
}

func (c *ConesearchService) nearestNeighborSearch(objects []repository.Mastercat, ra, dec, radius float64, maxNeighbors int) []repository.Mastercat {
	pts := []kdtree.Point{}
	for _, obj := range objects {
		pts = append(pts, KNNObject{Obj: obj})
	}
	tree := kdtree.New(pts)

	nearObjs := tree.KNN(&points.Point2D{X: ra, Y: dec}, maxNeighbors)

	// now we need to check that distance between nearest objects and center is actually lower than radius
	result := []repository.Mastercat{}
	for _, obj := range nearObjs {
		if euclideanDistance(obj, &points.Point2D{X: ra, Y: dec}) > radius {
			continue
		}
		result = append(result, obj.(KNNObject).Obj)
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
