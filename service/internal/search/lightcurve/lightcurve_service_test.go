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
package lightcurve

import (
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestDetection struct {
	ID             string
	ObjectId       string
	Magnitude      float64
	MagnitudeError float32
	MJD            float64
}

func (d TestDetection) GetId() string {
	return d.ID
}

func (d TestDetection) GetObjectId() string {
	return d.ObjectId
}

func (d TestDetection) GetBrightness() float64 {
	return d.Magnitude
}

func (d TestDetection) GetBrightnessError() float32 {
	return d.MagnitudeError
}

func (d TestDetection) GetMjd() float64 {
	return d.MJD
}

func MockLightcurveFilter(Lightcurve, []conesearch.MetadataResult) Lightcurve {
	return Lightcurve{}
}

func TestGetObjectIds_Empty(t *testing.T) {
	mockService := NewMockConesearchService(t)
	mockService.EXPECT().FindMetadataByConesearch(
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		1,
		"all",
	).Return([]conesearch.MetadataResult{}, nil)

	lightcurveService, err := New([]ExternalClient{NewMockExternalClient(t)}, []LightcurveFilter{MockLightcurveFilter}, mockService)
	require.NoError(t, err)

	objs, err := lightcurveService.getObjects(0, 0, 0, 1)
	require.NoError(t, err)

	require.Equal(t, []conesearch.MetadataResult{}, objs)
}

func TestGetObjectIds_NonEmpty(t *testing.T) {
	mockService := NewMockConesearchService(t)
	mockService.EXPECT().FindMetadataByConesearch(
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("float64"),
		1,
		"all",
	).Return([]conesearch.MetadataResult{
		{
			Catalog: "allwise",
			Data:    []conesearch.MetadataExtended{{Metadata: repository.Gaia{ID: "ALLWISE1"}, Distance: 0.5}},
		},
		{
			Catalog: "gaia",
			Data:    []conesearch.MetadataExtended{{Metadata: repository.Gaia{ID: "GAIA1"}, Distance: 0.5}},
		},
	}, nil)

	lightcurveService, err := New([]ExternalClient{NewMockExternalClient(t)}, []LightcurveFilter{MockLightcurveFilter}, mockService)
	require.NoError(t, err)

	objs, err := lightcurveService.getObjects(0, 0, 0, 1)
	require.NoError(t, err)

	ids := make([]string, 0)
	for i := range objs {
		for j := range objs[i].Data {
			ids = append(ids, objs[i].Data[j].GetId())
		}
	}

	require.Equal(t, []string{"ALLWISE1", "GAIA1"}, ids)
}

func TestMergeClientResults_NoError(t *testing.T) {
	clientResult1 := ClientResult{
		Lightcurve: Lightcurve{
			Detections: []LightcurveObject{TestDetection{"1", "ALLWISE1", 1, 1, 1}},
		},
		Error: nil,
	}
	clientResult2 := ClientResult{
		Lightcurve: Lightcurve{
			Detections: []LightcurveObject{TestDetection{"2", "GAIA1", 1, 1, 1}},
		},
		Error: nil,
	}

	results := make(chan ClientResult, 2)
	results <- clientResult1
	results <- clientResult2
	close(results)

	lightcurveService, err := New([]ExternalClient{NewMockExternalClient(t)}, []LightcurveFilter{MockLightcurveFilter}, NewMockConesearchService(t))
	require.NoError(t, err)

	lightcurve, err := lightcurveService.mergeClientResults(results)
	require.NoError(t, err)

	require.Equal(t, []LightcurveObject{
		TestDetection{"1", "ALLWISE1", 1, 1, 1},
		TestDetection{"2", "GAIA1", 1, 1, 1},
	}, lightcurve.Detections)
}

func TestMergeClientResults_WithError(t *testing.T) {
	clientResult1 := ClientResult{
		Lightcurve: Lightcurve{
			Detections: []LightcurveObject{TestDetection{"1", "ALLWISE1", 1, 1, 1}},
		},
		Error: nil,
	}
	clientResult2 := ClientResult{
		Lightcurve: Lightcurve{},
		Error:      fmt.Errorf("first error"),
	}
	clientResult3 := ClientResult{
		Lightcurve: Lightcurve{},
		Error:      fmt.Errorf("second error"),
	}

	results := make(chan ClientResult, 3)
	results <- clientResult1
	results <- clientResult2
	results <- clientResult3
	close(results)

	lightcurveService, err := New([]ExternalClient{NewMockExternalClient(t)}, []LightcurveFilter{MockLightcurveFilter}, NewMockConesearchService(t))
	require.NoError(t, err)

	lightcurve, err := lightcurveService.mergeClientResults(results)
	require.EqualError(t, err, "first error")

	require.Equal(t, []LightcurveObject{
		TestDetection{"1", "ALLWISE1", 1, 1, 1},
	}, lightcurve.Detections)
}

func TestMergeLightcurves(t *testing.T) {
	lightcurve1 := Lightcurve{
		Detections: []LightcurveObject{TestDetection{"1", "ALLWISE1", 1, 1, 1}},
	}
	lightcurve2 := Lightcurve{
		Detections: []LightcurveObject{TestDetection{"2", "GAIA1", 1, 1, 1}},
	}

	service, err := New([]ExternalClient{NewMockExternalClient(t)}, []LightcurveFilter{MockLightcurveFilter}, NewMockConesearchService(t))
	require.NoError(t, err)

	result := service.mergeLightcurves([]Lightcurve{lightcurve1, lightcurve2})

	require.Equal(t, []LightcurveObject{
		TestDetection{"1", "ALLWISE1", 1, 1, 1},
		TestDetection{"2", "GAIA1", 1, 1, 1},
	}, result.Detections)
}
