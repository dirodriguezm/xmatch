package conesearch

import (
	"context"
	"log/slog"
	"math"

	"github.com/dirodriguezm/xmatch/service/internal/core/knn"
	"github.com/dirodriguezm/xmatch/service/pkg/assertions"
	"github.com/dirodriguezm/xmatch/service/pkg/repository"

	"github.com/dirodriguezm/healpix"
)

type Repository interface {
	FindObjects(ctx context.Context, pixelList []int64) ([]repository.Mastercat, error)
	InsertObject(ctx context.Context, object repository.InsertObjectParams) (repository.Mastercat, error)
}

type ConesearchService struct {
	Scheme     healpix.OrderingScheme
	Nside      int
	Resolution int
	Catalog    string
	repository Repository
}

func NewConesearchService(options ...ConesearchOption) (*ConesearchService, error) {
	service := &ConesearchService{}
	for _, opt := range options {
		err := opt(service)
		if err != nil {
			return nil, err
		}
	}
	assertions.NotNil(service.repository)
	assertions.NotZero(service.Nside)
	assertions.NotZero(service.Catalog)
	assertions.NotZero(service.Scheme)
	if service.Resolution == 0 {
		service.Resolution = 4
	}
	slog.Debug("Created new ConesearchService", "repository", service.repository, "nside", service.Nside,
		"scheme", service.Scheme, "catalog", service.Catalog, "resolution", service.Resolution)
	return service, nil
}

func (c *ConesearchService) Conesearch(ra, dec, radius float64, nneighbor int) ([]repository.Mastercat, error) {
	if err := ValidateRa(ra); err != nil {
		return nil, err
	}
	if err := ValidateDec(dec); err != nil {
		return nil, err
	}
	if err := ValidateRadius(radius); err != nil {
		return nil, err
	}
	if err := ValidateNneighbor(nneighbor); err != nil {
		return nil, err
	}

	mapper, err := healpix.NewHEALPixMapper(c.Nside, c.Scheme)
	if err != nil {
		return nil, err
	}
	radius = arcsecToRadians(radius)
	point := healpix.RADec(float64(ra), float64(dec))
	pixelRanges := mapper.QueryDiscInclusive(point, radius, c.Resolution)
	pixelList := pixelRangeToList(pixelRanges)
	objs, err := c.getObjects(pixelList)
	if err != nil {
		return nil, err
	}
	objs = knn.NearestNeighborSearch(objs, ra, dec, radius, nneighbor)
	return objs, err
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
	ctx := context.Background()
	objects, err := c.repository.FindObjects(ctx, pixelList)
	if err != nil {
		return nil, err
	}
	return objects, nil
}