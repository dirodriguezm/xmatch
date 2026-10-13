package repository

import (
	"context"

	"github.com/dirodriguezm/healpix"
)

type MastercatReader interface {
	FindObjectsInPixelRanges(context.Context, []healpix.PixelRange) ([]Mastercat, error)
}
