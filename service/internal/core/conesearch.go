package core

import (
	"math"
	"xmatch/service/internal/repository"
	"xmatch/service/pkg/assertions"

	"github.com/dirodriguezm/healpix"
)

type Repository interface {
	FindObjectIds(pixelList []int64) ([]string, error)
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
	return service, nil
}

func (c *ConesearchService) Conesearch(ra, dec, radius float64, nneighbor int) ([]string, error) {
	scheme := healpix.Nest
	nside := 14
	discResolution := 4
	mapper, err := healpix.NewHEALPixMapper(nside, scheme)
	if err != nil {
		return nil, err
	}
	radius = arcsecToRadians(radius)
	point := healpix.RADec(float64(ra), float64(dec))
	pixelRange := mapper.QueryDiscInclusive(point, radius, discResolution)
	pixelList := pixelRangeToList(pixelRange)
	oids, err := getObjectIds(pixelList, &repository.SqliteRepository{})
	// TODO: perform nearest neihbor search to filter oids
	return oids, err
}

func arcsecToRadians(arcsec float64) float64 {
	return (arcsec / 3600) * (math.Pi / 180)
}

func pixelRangeToList(pixelRange []healpix.PixelRange) []int64 {
	result := make([]int64, 0, len(pixelRange))
	for _, r := range pixelRange {
		for i := r.Start; i < r.Stop; i++ {
			result = append(result, i)
		}
	}
	return result
}

func getObjectIds(pixelList []int64, rep Repository) ([]string, error) {
	oids, err := rep.FindObjectIds(pixelList)
	if err != nil {
		return nil, err
	}
	return oids, nil
}
