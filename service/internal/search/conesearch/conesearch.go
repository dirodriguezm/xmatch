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
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/knn"
	"github.com/dirodriguezm/xmatch/service/internal/utils"

	"github.com/dirodriguezm/healpix"
)

type Repository interface {
	FindObjects(context.Context, []int64) ([]repository.Mastercat, error)
	InsertObject(context.Context, repository.InsertObjectParams) (repository.Mastercat, error)
	BulkInsertObject(context.Context, *sql.DB, []repository.InsertObjectParams) error
	GetAllObjects(context.Context) ([]repository.Mastercat, error)
	GetCatalogs(context.Context) ([]repository.Catalog, error)
	InsertCatalog(context.Context, repository.InsertCatalogParams) (repository.Catalog, error)
	GetDbInstance() *sql.DB
	InsertAllwise(context.Context, repository.InsertAllwiseParams) error
	GetAllwise(context.Context, string) (repository.Allwise, error)
	BulkInsertAllwise(context.Context, *sql.DB, []repository.InsertAllwiseParams) error
	RemoveAllObjects(context.Context) error
	BulkGetAllwise(context.Context, []string) ([]repository.Allwise, error)
	GetAllwiseFromPixels(context.Context, []int64) ([]repository.GetAllwiseFromPixelsRow, error)
}

type ConesearchService struct {
	Scheme     healpix.OrderingScheme
	Resolution int
	Catalogs   []repository.Catalog
	repository Repository
	mappers    map[int64]*healpix.HEALPixMapper
	ctx        context.Context
}

func NewConesearchService(options ...ConesearchOption) (*ConesearchService, error) {
	ctx := context.Background()
	service := &ConesearchService{
		Scheme:     healpix.Nest,
		Resolution: 4,
		Catalogs:   []repository.Catalog{},
		repository: nil,
		mappers:    map[int64]*healpix.HEALPixMapper{},
		ctx:        ctx,
	}
	for _, opt := range options {
		err := opt(service)
		if err != nil {
			return nil, err
		}
	}
	assertions.NotNil(service.repository)
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
	return mappers, nil
}

func (c *ConesearchService) Conesearch(ra, dec, radius float64, nneighbor int, catalog string) ([]repository.Mastercat, error) {
	if err := ValidateArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	radius_radians := arcsecToRadians(radius)
	point := healpix.RADec(float64(ra), float64(dec))
	objects := make([]repository.Mastercat, 0)
	for _, v := range c.mappers {
		pixelRanges := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
		pixelList := pixelRangeToList(pixelRanges)
		objs, err := c.getObjects(pixelList)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	objects = knn.NearestNeighborSearch(objects, ra, dec, radius, nneighbor)
	return objects, nil
}

func (c *ConesearchService) FindMetadataByConesearch(ra, dec, radius float64, nneighbor int, catalog string) (any, error) {
	if err := ValidateArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	radius_radians := arcsecToRadians(radius)
	point := healpix.RADec(float64(ra), float64(dec))

	performQuery := func(v *healpix.HEALPixMapper) (any, error) {
		pixelRanges := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
		pixelList := pixelRangeToList(pixelRanges)
		return c.getMetadata(pixelList)
	}

	switch strings.ToLower(catalog) {
	case "allwise":
		objects := make([]repository.GetAllwiseFromPixelsRow, 0)
		for _, v := range c.mappers {
			objs, err := performQuery(v)
			if err != nil {
				return nil, err
			}
			objects = append(objects, objs.([]repository.GetAllwiseFromPixelsRow)...)
		}
		final := knn.NearestNeighborSearchForAllwiseMetadata(objects, ra, dec, radius, nneighbor)
		return final, nil
	case "all":
		return nil, fmt.Errorf("using all is not supported for metadata search")
	default:
		return nil, fmt.Errorf("catalog %s not found", catalog)
	}
}

func (c *ConesearchService) BulkConesearch(
	ra, dec []float64, radius float64, nneighbor int, catalog string, chunkSize int, maxBulkConcurrency int,
) ([]repository.Mastercat, error) {
	if err := ValidateBulkArguments(ra, dec, radius, nneighbor, catalog); err != nil {
		return nil, err
	}

	radius_radians := arcsecToRadians(radius)
	numChunks := (len(ra) + chunkSize - 1) / chunkSize
	resultsChan := make(chan []repository.Mastercat, numChunks)
	errChan := make(chan error, numChunks)
	var wg sync.WaitGroup

	sem := make(chan struct{}, maxBulkConcurrency)

	for _, v := range c.mappers {
		for i := 0; i < len(ra); i += chunkSize {
			wg.Add(1)

			end := i + chunkSize
			if end > len(ra) {
				end = len(ra)
			}
			chunkRa := ra[i:end]
			chunkDec := dec[i:end]

			go func(chunkRa, chunkDec []float64) {
				sem <- struct{}{}

				defer func() {
					<-sem
					wg.Done()
				}()

				for j := 0; j < len(chunkRa); j++ {
					point := healpix.RADec(chunkRa[j], chunkDec[j])
					pixelRange := v.QueryDiscInclusive(point, radius_radians, c.Resolution)
					pixelList := pixelRangeToList(pixelRange)
					objs, err := c.getObjects(pixelList)
					if err != nil {
						errChan <- err
					}

					objs = knn.NearestNeighborSearch(objs, chunkRa[j], chunkDec[j], radius, nneighbor)
					resultsChan <- objs
				}

			}(chunkRa, chunkDec)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	allObjects := make([]repository.Mastercat, 0)
	for result := range resultsChan {
		allObjects = append(allObjects, result...)
	}
	for err := range errChan {
		return nil, err
	}

	uniqueObjects := make([]repository.Mastercat, 0)
	ids := utils.Set{}
	for i := 0; i < len(allObjects); i++ {
		if !ids.Contains(allObjects[i].ID) {
			uniqueObjects = append(uniqueObjects, allObjects[i])
			ids.Add(allObjects[i].ID)
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

func (c *ConesearchService) getObjects(pixelList []int64) ([]repository.Mastercat, error) {
	// TODO: include catalog name in search
	objects, err := c.repository.FindObjects(c.ctx, pixelList)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func (c *ConesearchService) getMetadata(pixelList []int64) (any, error) {
	objects, err := c.repository.GetAllwiseFromPixels(c.ctx, pixelList)
	if err != nil {
		return nil, err
	}
	return objects, nil
}
