package sqlite_writer

import (
	"context"
	"database/sql"
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
	params := make([]any, len(msg.Rows))
	for i, object := range msg.Rows {
		// convert the received row to insert params needed by the repository
		p, err := row2insertParams(object)
		if err != nil {
			slog.Error("SqliteWriter could not convert received object to insert params", "object", object)
			panic(err)
		}
		params[i] = p
	}
	// insert converted rows
	err := insertData(w.repository, w.ctx, w.repository.GetDbInstance(), params)
	if err != nil {
		slog.Error("SqliteWriter could not write objects to database")
		panic(err)
	}
}

func (w *SqliteWriter[T]) Stop() {
	w.DoneChannel <- true
}

func row2insertParams[T any](obj T) (any, error) {
	switch v := any(obj).(type) {
	case repository.ParquetMastercat:
		return v.ToInsertObjectParams(), nil
	case repository.AllwiseMetadata:
		return v.ToInsertParams(), nil
	default:
		return nil, fmt.Errorf("Parameter type not known: %T", v)
	}
}

func insertData[T any](repo conesearch.Repository, ctx context.Context, db *sql.DB, rows []T) error {
	insertObjectParams := make([]repository.InsertObjectParams, 0, len(rows))
	insertMetadataParams := make([]repository.InsertAllwiseParams, 0, len(rows))

	for i := range rows {
		if p, ok := any(rows[i]).(repository.InsertObjectParams); ok {
			insertObjectParams = append(insertObjectParams, p)
		} else if p, ok := any(rows[i]).(repository.InsertAllwiseParams); ok {
			insertMetadataParams = append(insertMetadataParams, p)
		} else {
			return fmt.Errorf("Parameter type not known: %T", rows[i])
		}
	}

	// now check which type of data we have and call the appropriate function
	if len(insertObjectParams) > 0 {
		return repo.BulkInsertObject(ctx, db, insertObjectParams)
	}
	if len(insertMetadataParams) > 0 {
		return repo.BulkInsertAllwise(ctx, db, insertMetadataParams)
	}
	return nil
}
