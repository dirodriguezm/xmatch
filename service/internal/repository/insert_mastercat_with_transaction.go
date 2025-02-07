package repository

import (
	"context"
	"database/sql"
)

func (q *Queries) BulkInsertObject(ctx context.Context, db *sql.DB, arg []InsertObjectParams) error {
	tx, err := db.Begin()
	if err != nil {
		return nil
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		_, err = qtx.InsertObject(ctx, arg[i])
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
