package sqlite_writer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
)

type SqliteWriter struct {
	bulkInsert func(context.Context, []any) error
	ctx        context.Context
}

func New(
	ctx context.Context,
	bulkInsert func(context.Context, []any) error,
) *SqliteWriter {
	slog.Debug("Creating new SqliteWriter")
	return &SqliteWriter{
		ctx:        ctx,
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

	err := w.bulkInsert(w.ctx, msg.Rows)
	if err != nil {
		panic(fmt.Errorf("SqliteWriter could not write objects to database: %w", err))
	}
}

func (w *SqliteWriter) Stop(a *actor.Actor) {
	slog.Debug("Stopping SqliteWriter", "insert into", w.bulkInsert)
}
