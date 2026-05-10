package repository

import "context"

type MastercatReader interface {
	FindObjects(context.Context, []int64) ([]Mastercat, error)
}

type MastercatWriter interface {
	BulkInsertObject(context.Context, []any) error
}
