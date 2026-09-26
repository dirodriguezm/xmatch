package repository

import (
	"context"

	"github.com/dirodriguezm/healpix"
)

type MastercatReader interface {
	FindObjectsInPixelRanges(context.Context, []healpix.PixelRange) ([]Mastercat, error)
}

type MastercatWriter interface {
	BulkInsertObject(context.Context, []any) error
}
