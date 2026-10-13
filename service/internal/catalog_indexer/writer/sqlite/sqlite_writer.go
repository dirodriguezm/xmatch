package sqlite_writer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type SqliteWriter[T repository.Insertable] struct {
	ctx   context.Context
	store *repository.Queries
}

// New builds a writer that bulk inserts rows of T through the store. T
// implements repository.Insertable, so the writer needs no per-catalog insert
// function.
func New[T repository.Insertable](
	ctx context.Context,
	store *repository.Queries,
) *SqliteWriter[T] {
	slog.Debug("Creating new SqliteWriter")
	return &SqliteWriter[T]{
		ctx:        ctx,
		store: store,
	}
}

func (w *SqliteWriter[T]) Write(a *actor.Actor[T], msg actor.Message[T]) {
	defer func() {
		if r := recover(); r != nil {
			w.Stop(a)
			panic(r)
		}
	}()

	slog.Debug("SqliteWriter received message", "rows", len(msg.Rows))
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(fmt.Errorf("SqliteWriter received error: %w", msg.Error))
	}

	err := repository.BulkInsert(w.ctx, w.store, msg.Rows)
	if err != nil {
		panic(fmt.Errorf("SqliteWriter could not write objects to database: %w", err))
	}
}

func (w *SqliteWriter[T]) Stop(a *actor.Actor[T]) {
	slog.Debug("Stopping SqliteWriter")
}
