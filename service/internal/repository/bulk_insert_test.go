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
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkInsertSuccessfulCommitDoesNotLogRollbackError(t *testing.T) {
	var logs bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	q := New(db)
	err = BulkInsert[int](context.Background(), q, []any{1}, func(ctx context.Context, qtx *Queries, row int) error {
		_, err := qtx.db.ExecContext(ctx, "SELECT 1")
		return err
	})

	require.NoError(t, err)
	assert.NotContains(t, logs.String(), "rollback error")
	assert.NotContains(t, logs.String(), "transaction has already been committed or rolled back")
}
