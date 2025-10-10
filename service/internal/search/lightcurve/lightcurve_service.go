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
	"slices"
	"sync"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type ExternalClient interface {
	FetchLightcurve(float64, float64, float64, int) ClientResult
}

type ConesearchService interface {
	Conesearch(float64, float64, float64, int, string) ([]conesearch.MastercatResult, error)
}

type ClientResult struct {
	Lightcurve Lightcurve
	Error      error
}

type LightcurveService struct {
	externalClients   []ExternalClient
	conesearchService ConesearchService
}

func New(externalClients []ExternalClient, conesearchService ConesearchService) (*LightcurveService, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("conesearchService was nil while creating LightcurveService")
	}
	return &LightcurveService{externalClients, conesearchService}, nil
}

// GetLightcurve retrieves lightcurve data by querying multiple external clients concurrently.
// It takes celestial coordinates (ra, dec), search radius, and maximum number of objects as parameters,
// then merges the results from all external clients into a single Lightcurve.
//
// Sequence diagram (same row means concurrent, lines mean waiting):
//
//	ClientA - ClientB - ... - ClientN - Xwave Conesearch
//
//	----------------------------------------------
//
//	Xwave Lightcurve - Filter client results by Conesearch IDs
//
//	----------------------------------------------
//
//	Merge lightcurves from xwave and external clients
//
// # Parameters:
//   - ra: Right ascension coordinate in degrees
//   - dec: Declination coordinate in degrees
//   - radius: Search radius in degrees
//   - nobjects: Maximum number of objects to retrieve
//
// Returns:
//   - Lightcurve: Combined lightcurve data from all external clients
//   - error: Any error encountered during the fetch operation
func (service *LightcurveService) GetLightcurve(ra, dec, radius float64, nobjects int) (Lightcurve, error) {
	var wg sync.WaitGroup
	results := make(chan ClientResult, len(service.externalClients))

	for _, externalClient := range service.externalClients {
		wg.Add(1)
		go func(client ExternalClient) {
			defer wg.Done()
			results <- client.FetchLightcurve(ra, dec, radius, nobjects)
		}(externalClient)
	}

	wg.Add(1)
	var objectIds []string
	var conesearchError error
	go func() {
		defer wg.Done()
		objectIds, conesearchError = service.getObjectIds(ra, dec, radius, nobjects)
	}()

	wg.Wait()
	close(results)

	if conesearchError != nil {
		return Lightcurve{}, fmt.Errorf("could not execute conesearch: %w", conesearchError)
	}

	mergedResults, err := service.mergeClientResults(results)
	if err != nil {
		return mergedResults, err
	}

	var xWaveLightcurve Lightcurve
	var xWaveError error
	wg.Add(1)
	go func() {
		xWaveLightcurve, xWaveError = service.getXwaveLightcurve(objectIds)
	}()

	var filteredLightcurve Lightcurve
	wg.Add(1)
	go func() {
		filteredLightcurve = service.filterById(objectIds, mergedResults)
	}()

	wg.Wait()
	if xWaveError != nil {
		return Lightcurve{}, fmt.Errorf("could not get x-wave lightcurve: %w", xWaveError)
	}

	return service.mergeLightcurves([]Lightcurve{xWaveLightcurve, filteredLightcurve}), nil
}

// getObjectIds retrieves object IDs from the conesearch service for the given celestial coordinates.
// It queries the conesearch service with the specified parameters and extracts object IDs from the results.
//
// Parameters:
//   - ra: Right ascension coordinate in degrees
//   - dec: Declination coordinate in degrees
//   - radius: Search radius in degrees
//   - neighbors: Maximum number of objects to retrieve
//
// Returns:
//   - []string: Slice of object IDs found in the search area
//   - error: Any error encountered during the conesearch operation
func (service *LightcurveService) getObjectIds(ra, dec, radius float64, neighbors int) ([]string, error) {
	ids := make([]string, 0)
	objects, err := service.conesearchService.Conesearch(ra, dec, radius, neighbors, "all")
	if err != nil {
		return nil, fmt.Errorf("could not execute conesearch: %w", err)
	}
	for _, objectsByCatalog := range objects {
		for _, object := range objectsByCatalog.Data {
			ids = append(ids, object.ID)
		}
	}
	return ids, nil
}

// getXwaveLightcurve retrieves lightcurve data from the Xwave service for the given object IDs.
// Currently returns an empty lightcurve as a placeholder implementation.
//
// Parameters:
//   - ids: Slice of object IDs to fetch lightcurve data for
//
// Returns:
//   - Lightcurve: Lightcurve data from Xwave service
//   - error: Any error encountered during the fetch operation
func (service *LightcurveService) getXwaveLightcurve(_ []string) (Lightcurve, error) {
	return Lightcurve{}, nil
}

// filterById filters a lightcurve to include only detections, non-detections, and forced photometry
// that belong to the specified object IDs. This is used to ensure only objects indexed by xwave
// are returned from external clients is in the final result.
//
// Parameters:
//   - ids: Slice of object IDs to filter by
//   - lightcurve: Lightcurve data to filter
//
// Returns:
//   - Lightcurve: Filtered lightcurve containing only data for the specified object IDs
func (service *LightcurveService) filterById(ids []string, lightcurve Lightcurve) Lightcurve {
	newLightcurve := Lightcurve{}

	for _, detection := range lightcurve.Detections {
		if slices.Contains(ids, detection.GetObjectId()) {
			newLightcurve.Detections = append(newLightcurve.Detections, detection)
		}
	}

	for _, nonDetection := range lightcurve.NonDetections {
		if slices.Contains(ids, nonDetection.GetObjectId()) {
			newLightcurve.NonDetections = append(newLightcurve.NonDetections, nonDetection)
		}
	}

	for _, forcedPhotometry := range lightcurve.ForcedPhotometry {
		if slices.Contains(ids, forcedPhotometry.GetObjectId()) {
			newLightcurve.ForcedPhotometry = append(newLightcurve.ForcedPhotometry, forcedPhotometry)
		}
	}

	return newLightcurve
}

// mergeClientResults merges lightcurve data from multiple external client results received through a channel.
// It aggregates detections, non-detections, and forced photometry from all successful client responses.
// Returns early with an error if any client result contains an error.
//
// Parameters:
//   - results: Channel of ClientResult containing lightcurve data from external clients
//
// Returns:
//   - Lightcurve: Merged lightcurve data from all successful clients
//   - error: First error encountered from any client, if any
func (service *LightcurveService) mergeClientResults(results <-chan ClientResult) (Lightcurve, error) {
	lightcurve := Lightcurve{}

	for result := range results {
		if result.Error != nil {
			return lightcurve, result.Error
		}

		lightcurve.Detections = append(lightcurve.Detections, result.Lightcurve.Detections...)
		lightcurve.NonDetections = append(lightcurve.NonDetections, result.Lightcurve.NonDetections...)
		lightcurve.ForcedPhotometry = append(lightcurve.ForcedPhotometry, result.Lightcurve.ForcedPhotometry...)
	}

	return lightcurve, nil
}

// mergeLightcurves combines multiple lightcurves into a single lightcurve by concatenating
// all detections, non-detections, and forced photometry from each input lightcurve.
//
// Parameters:
//   - lightcurves: Slice of Lightcurve to merge
//
// Returns:
//   - Lightcurve: Combined lightcurve containing all data from input lightcurves
func (service *LightcurveService) mergeLightcurves(lightcurves []Lightcurve) Lightcurve {
	newLightcurve := Lightcurve{}

	for _, lightcurve := range lightcurves {
		for _, detection := range lightcurve.Detections {
			newLightcurve.Detections = append(newLightcurve.Detections, detection)
		}
		for _, nonDetection := range lightcurve.NonDetections {
			newLightcurve.NonDetections = append(newLightcurve.NonDetections, nonDetection)
		}
		for _, forcedPhotometry := range lightcurve.ForcedPhotometry {
			newLightcurve.ForcedPhotometry = append(newLightcurve.ForcedPhotometry, forcedPhotometry)
		}
	}

	return newLightcurve
}
