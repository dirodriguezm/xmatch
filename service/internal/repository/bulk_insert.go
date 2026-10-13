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

func (q *Queries) beginBulkInsertTx(ctx context.Context) (*sql.Tx, error) {
	db, ok := q.db.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("bulk inserts require repository backed by *sql.DB")
	}
	return db.BeginTx(ctx, nil)
}

func BulkInsert[T Insertable](ctx context.Context, q *Queries, rows []T) error {
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
	for i := range rows {
		err = rows[i].Insert(ctx, qtx)
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
