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

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter struct {
	repository conesearch.Repository
	bulkInsert func(context.Context, *sql.DB, []any) error
	ctx        context.Context
	db         *sql.DB
}

func New(
	repo conesearch.Repository,
	ctx context.Context,
	bulkInsert func(context.Context, *sql.DB, []any) error,
) *SqliteWriter {
	slog.Debug("Creating new SqliteWriter")
	return &SqliteWriter{
		repository: repo,
		ctx:        ctx,
		db:         repo.GetDbInstance(),
		bulkInsert: bulkInsert,
	}
}

func (w *SqliteWriter) Write(a *actor.Actor, msg actor.Message) {
	defer func() {
		if r := recover(); r != nil {
			w.Stop(a)
			panic(r)
		}
	}()

	slog.Debug("SqliteWriter received message", "insert into", w.bulkInsert)
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(fmt.Errorf("SqliteWriter received error: %w", msg.Error))
	}

	err := w.bulkInsert(w.ctx, w.db, msg.Rows)
	if err != nil {
		panic(fmt.Errorf("SqliteWriter could not write objects to database: %w", err))
	}
}

func (w *SqliteWriter) Stop(a *actor.Actor) {
	slog.Debug("Stopping SqliteWriter", "insert into", w.bulkInsert)
	err := w.db.Close()
	if err != nil {
		panic(err)
	}
}
