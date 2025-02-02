package sqlite_writer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter[T any] struct {
	*writer.BaseWriter[T]
	repository conesearch.Repository
	ctx        context.Context
	src        *source.Source
}

func NewSqliteWriter[T any](
	repository conesearch.Repository,
	ch chan indexer.WriterInput[T],
	done chan bool,
	ctx context.Context,
	src *source.Source,
) *SqliteWriter[T] {
	slog.Debug("Creating new SqliteWriter")
	w := &SqliteWriter[T]{
		BaseWriter: &writer.BaseWriter[T]{
			DoneChannel:  done,
			InboxChannel: ch,
		},
		repository: repository,
		ctx:        ctx,
		src:        src,
	}
	w.Writer = w
	return w
}

func (w *SqliteWriter[T]) Receive(msg indexer.WriterInput[T]) {
	slog.Debug("Writer received message")
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(msg.Error)
	}
	for _, object := range msg.Rows {
		// convert the received row to insert params needed by the repository
		params, err := row2insertParams(object)
		if err != nil {
			slog.Error("SqliteWriter could not convert received object to insert params", "object", object)
			panic(err)
		}

		// insert converted rows
		err = insertData(w.repository, w.ctx, params)
		if err != nil {
			slog.Error("SqliteWriter could not write object to database", "object", object)
			panic(err)
		}
	}
}

func (w *SqliteWriter[T]) Stop() {
	w.DoneChannel <- true
}

func row2insertParams[T any](obj T) (any, error) {
	switch v := any(obj).(type) {
	case repository.ParquetMastercat:
		return v.ToInsertObjectParams(), nil
	default:
		return nil, fmt.Errorf("Parameter type not known: %T", v)
	}
}

func insertData[T any](repo conesearch.Repository, ctx context.Context, row T) error {
	switch v := any(row).(type) {
	case repository.InsertObjectParams:
		_, err := repo.InsertObject(ctx, v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Parameter type not known: %T", v)
	}
	return nil
}
