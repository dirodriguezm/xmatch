// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type BulkInsertFunc[T any] func(context.Context, *Queries, T) error

func (q *Queries) beginBulkInsertTx(ctx context.Context) (*sql.Tx, error) {
	db, ok := q.db.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("bulk inserts require repository backed by *sql.DB")
	}
	return db.BeginTx(ctx, nil)
}

func BulkInsert[T any](ctx context.Context, q *Queries, rows []any, insert BulkInsertFunc[T]) error {
	tx, err := q.beginBulkInsertTx(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			slog.Error("rollback error", "err", rollbackErr)
		}
	}()
	qtx := q.WithTx(tx)
	var zero T
	for i := range rows {
		row, ok := rows[i].(T)
		if !ok {
			return fmt.Errorf("bulk insert row %d: expected %T, got %T", i, zero, rows[i])
		}
		err = insert(ctx, qtx, row)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	committed = true
	return nil
}

func (q *Queries) BulkInsertObject(ctx context.Context, rows []any) error {
	return BulkInsert(ctx, q, rows, func(ctx context.Context, qtx *Queries, row Mastercat) error {
		return qtx.InsertObject(ctx, InsertObjectParams(row))
	})
}
