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

package testutils

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// CatalogDB returns a repository backed by a fresh SQLite database
// with every migration applied. Migrations are resolved from this file, so
// callers can live in any package of the module.
func CatalogDB(t *testing.T) *repository.Queries {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "db", "migrations")

	dbFile := filepath.Join(t.TempDir(), "catalog-roundtrip.db")
	mig, err := migrate.New(fmt.Sprintf("file://%s", migrationsDir), fmt.Sprintf("sqlite3://%s", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = mig.Close()
	})
	require.NoError(t, mig.Up())

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000", dbFile))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	require.NoError(t, db.Ping())

	return repository.New(db)
}
