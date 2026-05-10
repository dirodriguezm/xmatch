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

package conesearch

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"sync"

	"github.com/dirodriguezm/xmatch/service/internal/assertions"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/knn"

	"github.com/dirodriguezm/healpix"
)

type MastercatStore interface {
	FindObjects(context.Context, []int64) ([]repository.Mastercat, error)
	InsertMastercat(context.Context, repository.Mastercat) error
	GetAllObjects(context.Context) ([]repository.Mastercat, error)
	RemoveAllObjects(context.Context) error
	BulkInsertObject(context.Context, *sql.DB, []any) error
}

type CatalogRegistry interface {
	GetCatalogs(context.Context) ([]repository.Catalog, error)
	InsertCatalog(context.Context, repository.InsertCatalogParams) error
	GetDbInstance() *sql.DB
}

type indexedResult struct {
	index  int
	result knn.KnnResult[repository.Mastercat]
}

type ConesearchService struct {
	Scheme     healpix.OrderingScheme
	Resolution int
	Catalogs   []repository.Catalog
	store      MastercatStore
	resolver   *catalog.Resolver
	mappers    map[int64]*healpix.HEALPixMapper
	ctx        context.Context
}

func NewConesearchService(options ...ConesearchOption) (*ConesearchService, error) {
	ctx := context.Background()
	service := &ConesearchService{
		Scheme:     healpix.Nest,
		Resolution: 4,
		Catalogs:   []repository.Catalog{},
		store:      nil,
		resolver:   nil,
		mappers:    map[int64]*healpix.HEALPixMapper{},
		ctx:        ctx,
	}
	for _, opt := range options {
		err := opt(service)
		if err != nil {
			return nil, err
		}
	}
	assertions.NotNil(service.store)
	assertions.NotZero(service.Catalogs)
	assertions.NotZero(service.Scheme)

	var err error
	service.mappers, err = createServiceMappers(service.Catalogs, service.Scheme)
	if err != nil {
		return nil, err
	}

	slog.Debug("Created new ConesearchService", "scheme", service.Scheme, "catalogs", service.Catalogs, "resolution", service.Resolution)
	return service, nil
}

func createServiceMappers(catalogs []repository.Catalog, scheme healpix.OrderingScheme) (map[int64]*healpix.HEALPixMapper, error) {
	if len(catalogs) == 0 {
		return nil, fmt.Errorf("catalogs was empty while creating service mappers")
	}
	mappers := make(map[int64]*healpix.HEALPixMapper)
	for i := range catalogs {
		if _, ok := mappers[catalogs[i].Nside]; ok {
			continue
		}
		mapper, err := healpix.NewHEALPixMapper(int(catalogs[i].Nside), scheme)
		if err != nil {
			return nil, err
		}
		mappers[catalogs[i].Nside] = mapper
	}

	if len(mappers) == 0 {
		return nil, fmt.Errorf("mappers was empty while creating service mappers")
	}

	return mappers, nil
}

func (c *ConesearchService) Conesearch(ra, dec, radius float64, nneighbor int, catalog string) ([]MastercatResult, error) {
	if err := ValidateArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	radius_radians := arcsecToRadians(radius)
	point := healpix.RADec(float64(ra), float64(dec))
	objects := make([]repository.Mastercat, 0)
	for _, v := range c.mappers {
		pixelRanges := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
		pixelList := pixelRangeToList(pixelRanges)
		objs, err := c.getObjects(pixelList, catalog)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return ResultFromKnn(knn.NearestNeighborSearch(objects, ra, dec, radius, nneighbor), 0), nil
}

func (c *ConesearchService) FindMetadataByConesearch(
	ra, dec, radius float64,
	nneighbor int,
	catalog string,
) ([]MetadataResult, error) {
	if err := ValidateArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	objects, err := findMetadata(healpix.RADec(float64(ra), float64(dec)), arcsecToRadians(radius), c, catalog)
	if err != nil {
		return nil, fmt.Errorf("could not find allwise metadata: %w", err)
	}

	return ResultFromKnnMetadata(knn.NearestNeighborSearchForMetadata(objects, ra, dec, radius, nneighbor, catalog)), nil
}

func findMetadata(
	point healpix.Pointing,
	radius_radians float64,
	c *ConesearchService,
	catalog string,
) ([]repository.MetadataWithCoordinates, error) {
	objects := make([]repository.MetadataWithCoordinates, 0)
	for _, v := range c.mappers {
		pixelRanges := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
		pixelList := pixelRangeToList(pixelRanges)
		objs, err := c.getMetadata(pixelList, catalog)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}
	return objects, nil
}

func (c *ConesearchService) BulkConesearch(
	ra, dec []float64,
	radius float64,
	nneighbor int,
	catalog string,
	chunkSize int,
	maxBulkConcurrency int,
) ([]MastercatResult, error) {
	if err := ValidateBulkArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	radius_radians := arcsecToRadians(radius)
	numChunks := (len(ra) + chunkSize - 1) / chunkSize
	resultsChan := make(chan indexedResult, numChunks)
	errChan := make(chan error, numChunks)
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxBulkConcurrency)

	for _, v := range c.mappers {
		for i := 0; i < len(ra); i += chunkSize {
			wg.Add(1)

			end := min(i+chunkSize, len(ra))
			chunkRa := ra[i:end]
			chunkDec := dec[i:end]

			go func(chunkRa, chunkDec []float64, baseIndex int) {
				sem <- struct{}{}

				defer func() {
					<-sem
					wg.Done()
				}()

				for j := range chunkRa {
					point := healpix.RADec(chunkRa[j], chunkDec[j])
					pixelRange := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
					pixelList := pixelRangeToList(pixelRange)
					objs, err := c.getObjects(pixelList, catalog)
					if err != nil {
						errChan <- err
						break
					}

					resultsChan <- indexedResult{
						index:  baseIndex + j,
						result: knn.NearestNeighborSearch(objs, chunkRa[j], chunkDec[j], radius, nneighbor),
					}
				}

			}(chunkRa, chunkDec, i)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	resultsByIndex := make([][]MastercatResult, len(ra))
	for indexed := range resultsChan {
		resultsByIndex[indexed.index] = ResultFromKnn(indexed.result, indexed.index)
	}
	for err := range errChan {
		return nil, err
	}

	uniqueObjects := make([]MastercatResult, 0)
	seenIDs := make(map[int]map[string]bool)
	for i := range resultsByIndex {
		if resultsByIndex[i] == nil {
			resultsByIndex[i] = []MastercatResult{}
		}
		if seenIDs[i] == nil {
			seenIDs[i] = make(map[string]bool)
		}
		for _, mastercatResult := range resultsByIndex[i] {
			for j := range mastercatResult.Data {
				id := mastercatResult.Data[j].ID
				if !seenIDs[i][id] {
					seenIDs[i][id] = true
					uniqueObjects = append(uniqueObjects, MastercatResult{
						Catalog: mastercatResult.Catalog,
						Data:    []MastercatExtended{mastercatResult.Data[j]},
						Index:   i,
					})
				}
			}
		}
	}
	return uniqueObjects, nil
}

func arcsecToRadians(arcsec float64) float64 {
	return (arcsec / 3600) * (math.Pi / 180)
}

func pixelRangeToList(pixelRanges []healpix.PixelRange) []int64 {
	result := make([]int64, 0, len(pixelRanges))
	for _, r := range pixelRanges {
		for i := r.Start; i < r.Stop; i++ {
			result = append(result, i)
		}
	}
	return result
}

func (c *ConesearchService) getObjects(pixelList []int64, catalog string) ([]repository.Mastercat, error) {
	objects, err := c.store.FindObjects(c.ctx, pixelList)
	if err != nil {
		return nil, err
	}
	if catalog != "all" {
		return filterByCatalog(objects, catalog), nil
	}
	return objects, nil
}

func (c *ConesearchService) getMetadata(pixelList []int64, catalogName string) ([]repository.MetadataWithCoordinates, error) {
	objects := make([]repository.MetadataWithCoordinates, 0)

	catalogList := c.resolveCatalogList(catalogName)
	for _, name := range catalogList {
		adapter, err := c.resolver.Get(name)
		if err != nil {
			return nil, err
		}
		objs, err := adapter.GetFromPixels(c.ctx, pixelList)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}
	return objects, nil
}

func (c *ConesearchService) resolveCatalogList(catalogName string) []string {
	if strings.ToLower(catalogName) == "all" {
		names := make([]string, len(c.Catalogs))
		for i, cat := range c.Catalogs {
			names[i] = cat.Name
		}
		return names
	}
	return []string{catalogName}
}

func filterByCatalog(objects []repository.Mastercat, catalog string) []repository.Mastercat {
	result := make([]repository.Mastercat, 0)
	catalog = strings.ToLower(catalog)
	for _, obj := range objects {
		if strings.ToLower(obj.Cat) == catalog {
			result = append(result, obj)
		}
	}
	return result
}
