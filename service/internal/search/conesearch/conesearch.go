package conesearch

import (
	"context"
	"database/sql"
	"log/slog"
	"math"

	"github.com/dirodriguezm/xmatch/service/internal/assertions"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/knn"

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
