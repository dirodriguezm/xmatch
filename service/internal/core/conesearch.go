package core

import (
	"math"
	"xmatch/service/internal/repository"

	"github.com/dirodriguezm/healpix"
)

func Conesearch(ra, dec, radius float64, catalog string, nneighbor int) ([]string, error) {
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
	oids, err := getObjectIds(pixelList)
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

func getObjectIds(pixelList []int64) ([]string, error) {
	rep := repository.NewSqliteRepository()
	oids, err := rep.FindObjectIds(pixelList)
	if err != nil {
		return nil, err
	}
	return oids, nil
}
