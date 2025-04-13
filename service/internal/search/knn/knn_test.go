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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	expectedObjectIds []string
	objectList        []repository.Mastercat
	radius            float64
}

func TestKnn(t *testing.T) {
	objectList := Objects(t).WithNumObjects(5).Build()
	testCases := []testCase{
		{expectedObjectIds: []string{"0"}, objectList: objectList, radius: 1},
		{expectedObjectIds: []string{"0", "1"}, objectList: objectList, radius: 2},
		{expectedObjectIds: []string{"0", "1"}, objectList: objectList, radius: 3},
		{expectedObjectIds: []string{"0", "1", "2"}, objectList: objectList, radius: 4},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 5},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 6},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 7},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 8},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 9},
		{expectedObjectIds: []string{"0", "1", "2", "3"}, objectList: objectList, radius: 10},
		{expectedObjectIds: []string{"0", "1", "2", "3", "4"}, objectList: objectList, radius: 11},
	}

	for _, tc := range testCases {
		result := NearestNeighborSearch(tc.objectList, 179.5928264, 14.5297050, tc.radius, 5)
		require.Lenf(t, result, len(tc.expectedObjectIds), "Result objects are more than actual neighbors")
		for i, res := range result {
			require.Equal(t, res.ID, tc.expectedObjectIds[i])
		}
	}
}
