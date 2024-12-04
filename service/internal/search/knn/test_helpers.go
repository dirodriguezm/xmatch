package knn

import (
	"strconv"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type ObjectBuilder struct {
	t *testing.T

	NumObjects int
}

func Objects(t *testing.T) *ObjectBuilder {
	return &ObjectBuilder{t: t}
}

func (builder *ObjectBuilder) WithNumObjects(n int) *ObjectBuilder {
	builder.t.Helper()

	// TODO: remove this limit when I figure out how to create a list of coordinates
	limit := 5
	if n > limit {
		builder.NumObjects = 5
		return builder
	}
	builder.NumObjects = n
	return builder
}

func (builder *ObjectBuilder) Build() []repository.Mastercat {
	builder.t.Helper()

	raList := []float64{179.593, 179.59312500000001, 179.59375, 179.59416666666667, 179.5958333333333}
	objects := []repository.Mastercat{}
	for i := 0; i < builder.NumObjects; i++ {
		newObject := repository.Mastercat{
			ID:   strconv.FormatInt(int64(i), 10),
			Ipix: 1, // ipix does not really matter here, but objects in the KNN search will most probably have the same Ipix
			Ra:   raList[i],
			Dec:  14.5297050,
			Cat:  "test",
		}
		objects = append(objects, newObject)
	}
	return objects
}
