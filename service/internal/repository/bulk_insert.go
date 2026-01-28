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
)

func (q *Queries) BulkInsertObject(ctx context.Context, db *sql.DB, arg []any) error {
	tx, err := db.Begin()
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

func (q *Queries) BulkInsertAllwise(ctx context.Context, db *sql.DB, arg []any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertAllwiseWithoutParams(ctx, arg[i].(Allwise))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (q *Queries) BulkInsertGaia(ctx context.Context, db *sql.DB, arg []any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertGaiaWithoutParams(ctx, arg[i].(Gaia))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (q *Queries) BulkInsertErosita(ctx context.Context, db *sql.DB, arg []any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	for i := range arg {
		err = qtx.InsertErositaWithoutParams(ctx, arg[i].(Erosita))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
