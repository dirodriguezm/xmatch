package conesearch

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type ConesearchOption func(service *ConesearchService) error

func WithScheme(scheme healpix.OrderingScheme) ConesearchOption {
	return func(service *ConesearchService) error {
		service.Scheme = scheme
		return nil
	}
}

func WithResolution(res int) ConesearchOption {
	return func(service *ConesearchService) error {
		if service.Scheme == healpix.Nest && !isPowerOfTwo(res) {
			return errors.New("resolution must be a power of 2 when using Nest")
		}
		if res <= 0 {
			return errors.New("resolution must be a positive integer")
		}
		service.Resolution = res
		return nil
	}
}

func WithCatalogs(catalogs []repository.Catalog) ConesearchOption {
	return func(service *ConesearchService) error {
		allowed := []string{"vlass", "wise", "ztf"}
		for i := range catalogs {
			catName := strings.ToLower(catalogs[i].Name)
			if !slices.Contains(allowed, catName) {
				msg := fmt.Sprintf("specified catalog not available, please use one of %s", allowed)
				return errors.New(msg)
			}
		}
		service.Catalogs = catalogs
		return nil
	}
}

func WithRepository(repository Repository) ConesearchOption {
	return func(service *ConesearchService) error {
		service.repository = repository
		return nil
	}
}

func isPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	nfloat := float64(n)
	return math.Ceil(math.Log2(nfloat)) == math.Floor(math.Log2(nfloat))
}
