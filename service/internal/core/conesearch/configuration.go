package conesearch

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/dirodriguezm/healpix"
)

type ConesearchOption func(service *ConesearchService) error

func WithScheme(scheme healpix.OrderingScheme) ConesearchOption {
	return func(service *ConesearchService) error {
		service.Scheme = scheme
		return nil
	}
}

func WithNside(nside int) ConesearchOption {
	return func(service *ConesearchService) error {
		if nside < 1 || nside > 29 {
			return errors.New("nside must be between 1 and 29")
		}
		service.Nside = nside
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

func WithCatalog(catalog string) ConesearchOption {
	return func(service *ConesearchService) error {
		allowed := []string{"vlass"}
		catalog = strings.ToLower(catalog)
		for _, cat := range allowed {
			if catalog == cat {
				service.Catalog = catalog
				return nil
			}
		}
		msg := fmt.Sprintf("specified catalog not available, please use one of %s", allowed)
		return errors.New(msg)
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
