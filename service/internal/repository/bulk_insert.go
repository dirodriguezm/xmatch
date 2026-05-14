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
	"fmt"
)

func (q *Queries) beginBulkInsertTx(ctx context.Context) (*sql.Tx, error) {
	db, ok := q.db.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("bulk inserts require repository backed by *sql.DB")
	}
	return db.BeginTx(ctx, nil)
}

func (q *Queries) BulkInsertObject(ctx context.Context, arg []any) error {
	tx, err := q.beginBulkInsertTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertMastercat(ctx, arg[i].(Mastercat))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (q *Queries) BulkInsertAllwise(ctx context.Context, arg []any) error {
	tx, err := q.beginBulkInsertTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertAllwise(ctx, InsertAllwiseParams(arg[i].(Allwise)))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (q *Queries) BulkInsertGaia(ctx context.Context, arg []any) error {
	tx, err := q.beginBulkInsertTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertGaia(ctx, InsertGaiaParams(arg[i].(Gaia)))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (q *Queries) BulkInsertErosita(ctx context.Context, arg []any) error {
	tx, err := q.beginBulkInsertTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertErosita(ctx, InsertErositaParams(arg[i].(Erosita)))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
