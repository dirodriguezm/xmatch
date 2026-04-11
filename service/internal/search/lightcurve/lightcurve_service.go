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
	"strings"
	"sync"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type ExternalClient interface {
	FetchLightcurve(float64, float64, float64, int) ClientResult
}

type ConesearchService interface {
	FindMetadataByConesearch(float64, float64, float64, int, string) ([]conesearch.MetadataResult, error)
}

type LightcurveFilter func(Lightcurve, []conesearch.MetadataResult) Lightcurve

func DummyLightcurveFilter(l Lightcurve, _ []conesearch.MetadataResult) Lightcurve {
	return l
}

type ClientResult struct {
	Lightcurve Lightcurve
	Catalog    string
	Filter     LightcurveFilter
	Error      error
}

type Source struct {
	Catalog string
	Client  ExternalClient
	Filter  LightcurveFilter
}

type LightcurveService struct {
	sources           []Source
	conesearchService ConesearchService
}

func New(
	sources []Source,
	conesearchService ConesearchService,
) (*LightcurveService, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("conesearchService was nil while creating LightcurveService")
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("sources was empty while creating LightcurveService")
	}
	for i := range sources {
		if sources[i].Client == nil {
			return nil, fmt.Errorf("source client was nil while creating LightcurveService")
		}
		if sources[i].Filter == nil {
			sources[i].Filter = DummyLightcurveFilter
		}
	}
	return &LightcurveService{sources, conesearchService}, nil
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
func (service *LightcurveService) GetLightcurve(ra, dec, radius float64, nobjects int, catalog string) (Lightcurve, error) {
	selectedSources, err := service.selectSources(catalog)
	if err != nil {
		return Lightcurve{}, err
	}

	// Step 1: Fetch external clients and conesearch data concurrently
	clientData := make(chan ClientResult, len(selectedSources))
	service.fetchClientData(clientData, selectedSources, ra, dec, radius, nobjects)
	metadataResult := make(chan []conesearch.MetadataResult, 1)
	errors := make(chan error, 1)
	service.fetchConesearchData(metadataResult, errors, ra, dec, radius, nobjects, metadataCatalog(catalog))

	// Wait for all data to be fetched
	clientResults, err := service.collectClientResults(clientData)
	mergedClientResult := service.mergeCollectedClientResults(clientResults)
	if err != nil {
		return mergedClientResult, err
	}
	objectIds, objects, err := service.extractObjectIds(metadataResult, errors)
	if err != nil {
		return Lightcurve{}, err
	}
	if len(objectIds) == 0 {
		return mergedClientResult, nil
	}

	// Step 2: Fetch Xwave data and apply filters concurrently
	xwaveLightcurve := make(chan Lightcurve, 1)
	xwaveError := make(chan error, 1)
	service.getXwaveLightcurve(objectIds, xwaveLightcurve, xwaveError)
	filteredOutput := make(chan Lightcurve, 1)
	service.filterLightcurve(filteredOutput, objects, clientResults)

	// Wait for data to be fetched and filtered. Then merge the results.
	select {
	case xWaveLightcurve := <-xwaveLightcurve:
		return service.mergeLightcurves([]Lightcurve{xWaveLightcurve, <-filteredOutput}), nil
	case err := <-xwaveError:
		return Lightcurve{}, fmt.Errorf("could not get x-wave lightcurve: %w", err)
	}
}

// fetchClientData concurrently fetches lightcurve data from all external clients.
// It spawns a goroutine for each external client and sends the results through the output channel.
// The channel is closed once all goroutines have completed.
//
// Parameters:
//   - output: Channel to send ClientResult data to
//   - ra: Right ascension coordinate in degrees
//   - dec: Declination coordinate in degrees
//   - radius: Search radius in degrees
//   - nobjects: Maximum number of objects to retrieve
func (service *LightcurveService) fetchClientData(output chan<- ClientResult, sources []Source, ra, dec, radius float64, nobjects int) {
	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)
		go func(source Source) {
			defer wg.Done()
			result := source.Client.FetchLightcurve(ra, dec, radius, nobjects)
			result.Catalog = source.Catalog
			result.Filter = source.Filter
			output <- result
		}(source)
	}

	go func() {
		wg.Wait()
		close(output)
	}()
}

// fetchConesearchData concurrently fetches metadata from the conesearch service.
// It executes the conesearch query in a goroutine and sends results through the output channel.
// If an error occurs, it sends the error through the errors channel.
//
// Parameters:
//   - output: Channel to send metadata results to
//   - errors: Channel to send errors to
//   - ra: Right ascension coordinate in degrees
//   - dec: Declination coordinate in degrees
//   - radius: Search radius in degrees
//   - nobjects: Maximum number of objects to retrieve
func (service *LightcurveService) fetchConesearchData(
	output chan<- []conesearch.MetadataResult,
	errors chan<- error,
	ra, dec, radius float64,
	nobjects int,
	catalog string,
) {
	go func() {
		defer close(output)

		result, err := service.getObjects(ra, dec, radius, nobjects, catalog)
		if err != nil {
			errors <- err
			return
		}
		output <- result
	}()
}

// extractObjectIds extracts object IDs and metadata from conesearch results.
// It waits for metadata results from the conesearch service and extracts all object IDs
// and their corresponding metadata objects from the results.
//
// Parameters:
//   - metadataResult: Channel receiving conesearch metadata results
//   - errors: Channel receiving any errors from conesearch
//
// Returns:
//   - []string: Slice of extracted object IDs
//   - []MetadataExtended: Slice of metadata objects
//   - error: Any error encountered during extraction
func (service *LightcurveService) extractObjectIds(
	metadataResult <-chan []conesearch.MetadataResult,
	errors chan error,
) ([]string, []conesearch.MetadataResult, error) {
	objectIds := make([]string, 0)
	objects := make([]conesearch.MetadataResult, 0)
	select {
	case res := <-metadataResult:
		for i := range res {
			for j := range res[i].Data {
				objectIds = append(objectIds, res[i].Data[j].GetId())
			}
			objects = append(objects, conesearch.MetadataResult{Catalog: res[i].Catalog, Data: res[i].Data})
		}
		return objectIds, objects, nil
	case err := <-errors:
		return nil, nil, err
	}
}

// getObjects retrieves objects from the conesearch service for the given celestial coordinates.
// It queries the conesearch service with the specified parameters and extracts object IDs from the results.
//
// Parameters:
//   - ra: Right ascension coordinate in degrees
//   - dec: Declination coordinate in degrees
//   - radius: Search radius in degrees
//   - neighbors: Maximum number of objects to retrieve
//
// Returns:
//   - []MetadataResult: Slice of objects indexed by catalog found in the search area
//   - error: Any error encountered during the conesearch operation
func (service *LightcurveService) getObjects(ra, dec, radius float64, neighbors int, catalog string) ([]conesearch.MetadataResult, error) {
	objects, err := service.conesearchService.FindMetadataByConesearch(ra, dec, radius, neighbors, catalog)
	if err != nil {
		return nil, fmt.Errorf("could not execute conesearch: %w", err)
	}
	return objects, nil
}

// getXwaveLightcurve retrieves lightcurve data from the Xwave service for the given object IDs.
// Currently returns an empty lightcurve as a placeholder implementation.
//
// Parameters:
//   - ids: Slice of object IDs to fetch lightcurve data for
//   - result: Channel to send the result to
//   - err: Channel to send errors to
func (service *LightcurveService) getXwaveLightcurve(_ []string, result chan<- Lightcurve, _ chan error) {
	go func() {
		defer close(result)
		result <- Lightcurve{}
	}()
}

// filterLightcurve applies all configured lightcurve filters to the input lightcurve.
// Each filter is applied concurrently and the results are merged into a single lightcurve.
// The filtered result is sent through the output channel.
//
// Parameters:
//   - output: Channel to send the filtered lightcurve to
//   - objects: Metadata objects used for filtering
//   - lightcurve: Input lightcurve to be filtered
func (service *LightcurveService) filterLightcurve(output chan<- Lightcurve, objects []conesearch.MetadataResult, clientResults []ClientResult) {
	go func() {
		defer close(output)

		filteredLightcurves := make([]Lightcurve, len(clientResults))

		for i := range clientResults {
			filter := clientResults[i].Filter
			if filter == nil {
				filter = DummyLightcurveFilter
			}
			filteredLightcurves[i] = filter(clientResults[i].Lightcurve, objects)
		}

		output <- service.mergeLightcurves(filteredLightcurves)
	}()
}

func (service *LightcurveService) collectClientResults(results <-chan ClientResult) ([]ClientResult, error) {
	clientResults := make([]ClientResult, 0)

	for result := range results {
		if result.Error != nil {
			return clientResults, result.Error
		}

		clientResults = append(clientResults, result)
	}

	return clientResults, nil
}

func (service *LightcurveService) selectSources(catalog string) ([]Source, error) {
	normalizedCatalog := normalizeCatalog(catalog)
	if normalizedCatalog == "all" {
		return service.sources, nil
	}

	selectedSources := make([]Source, 0)
	for i := range service.sources {
		if service.sources[i].Catalog == normalizedCatalog {
			selectedSources = append(selectedSources, service.sources[i])
		}
	}
	if len(selectedSources) == 0 {
		return nil, fmt.Errorf("unsupported lightcurve catalog %q", catalog)
	}

	return selectedSources, nil
}

func metadataCatalog(catalog string) string {
	normalizedCatalog := normalizeCatalog(catalog)
	if normalizedCatalog == "neowise" {
		return "allwise"
	}
	return normalizedCatalog
}

func normalizeCatalog(catalog string) string {
	normalizedCatalog := strings.ToLower(catalog)
	if normalizedCatalog == "allwise" {
		return "neowise"
	}
	return normalizedCatalog
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
	clientResults, err := service.collectClientResults(results)
	if err != nil {
		return service.mergeCollectedClientResults(clientResults), err
	}

	return service.mergeCollectedClientResults(clientResults), nil
}

func (service *LightcurveService) mergeCollectedClientResults(results []ClientResult) Lightcurve {
	lightcurves := make([]Lightcurve, len(results))
	for i := range results {
		lightcurves[i] = results[i].Lightcurve
	}

	return service.mergeLightcurves(lightcurves)
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
