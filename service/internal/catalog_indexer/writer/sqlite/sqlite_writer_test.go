package sqlite_writer

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func TestReceive_Mastercat(t *testing.T) {
	mastercat := repository.Mastercat{
		ID:   "1",
		Ipix: int64(1),
		Ra:   1.0,
		Dec:  1.0,
		Cat:  "test",
	}

	bulkInsertFn := func(ctx context.Context, db *sql.DB, rows []any) error {
		return nil
	}

	w := New(nil, context.Background(), bulkInsertFn)
	w.Write(nil, actor.Message{Rows: []any{mastercat}, Error: nil})
}

func TestReceive_Allwise(t *testing.T) {
	allwise := repository.Allwise{
		ID:        "test",
		W1mpro:    repository.NullFloat64{sql.NullFloat64{Float64: 1.0, Valid: true}},
		W1sigmpro: repository.NullFloat64{sql.NullFloat64{Float64: 1.0, Valid: true}},
		W2mpro:    repository.NullFloat64{sql.NullFloat64{Float64: 2.0, Valid: true}},
		W2sigmpro: repository.NullFloat64{sql.NullFloat64{Float64: 2.0, Valid: true}},
	}

	bulkInsertFn := func(ctx context.Context, db *sql.DB, rows []any) error {
		return nil
	}

	w := New(nil, context.Background(), bulkInsertFn)
	w.Write(nil, actor.Message{Rows: []any{allwise}, Error: nil})
}
