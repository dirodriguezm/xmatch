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

package sqlite_writer

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter struct {
	repository conesearch.Repository
	ctx        context.Context
	db         *sql.DB
	table      string
}

func New(repo conesearch.Repository, ctx context.Context, table string) *SqliteWriter {
	slog.Debug("Creating new SqliteWriter")
	w := &SqliteWriter{repository: repo, ctx: ctx, db: repo.GetDbInstance(), table: table}
	return w
}

func (w *SqliteWriter) Write(a *actor.Actor, msg actor.Message) {
	defer func() {
		if r := recover(); r != nil {
			w.Stop(a)
			panic(r)
		}
	}()

	slog.Debug("SqliteWriter received message", "table", w.table)
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(fmt.Errorf("SqliteWriter received error: %w", msg.Error))
	}

	var err error
	switch strings.ToLower(w.table) {
	case "mastercat":
		slog.Debug("SqliteWriter Writing Mastercat", "len", len(msg.Rows))
		err = w.repository.BulkInsertObject(w.ctx, w.db, msg.Rows)
	case "allwise":
		slog.Debug("SqliteWriter Writing Allwise", "len", len(msg.Rows))
		err = w.repository.BulkInsertAllwise(w.ctx, w.db, msg.Rows)
	case "gaia":
		slog.Debug("SqliteWriter Writing Gaia", "len", len(msg.Rows))
		err = w.repository.BulkInsertGaia(w.ctx, w.db, msg.Rows)
	default:
		err = fmt.Errorf("Table %s not supported", w.table)
	}

	if err != nil {
		slog.Error("SqliteWriter could not write objects to database")
		panic(err)
	}
}

func (w *SqliteWriter) Stop(a *actor.Actor) {
	slog.Debug("Stopping SqliteWriter", "table", w.table)
	err := w.db.Close()
	if err != nil {
		panic(err)
	}
}
