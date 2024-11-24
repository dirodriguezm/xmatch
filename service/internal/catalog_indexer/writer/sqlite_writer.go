package writer

import (
	"context"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter struct {
	inbox      chan indexer.IndexerResult
	repository conesearch.Repository
	ctx        context.Context
	done       chan bool
}

func NewSqliteWriter(repository conesearch.Repository, ch chan indexer.IndexerResult, done chan bool, ctx context.Context) *SqliteWriter {
	slog.Debug("Creating new SqliteWriter")
	return &SqliteWriter{
		inbox:      ch,
		repository: repository,
		ctx:        ctx,
		done:       done,
	}
}

func (w *SqliteWriter) Start() {
	slog.Debug("Starting Writer")
	go func() {
		defer func() {
			slog.Debug("Writer Done")
			w.done <- true
			close(w.done)
		}()
		for msg := range w.inbox {
			w.receive(msg)
		}
	}()
}

func (w *SqliteWriter) receive(msg indexer.IndexerResult) {
	slog.Debug("Writer received message", "message", msg)
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(msg.Error)
	}
	for _, object := range msg.Objects {
		_, err := w.repository.InsertObject(w.ctx, masterCat2InsertParams(object))
		if err != nil {
			slog.Error("SqliteWriter could not write object to database")
			panic(err)
		}
	}
}

func (w *SqliteWriter) Done() {
	<-w.done
}

func masterCat2InsertParams(o1 repository.Mastercat) repository.InsertObjectParams {
	return repository.InsertObjectParams{
		ID:   o1.ID,
		Ipix: o1.Ipix,
		Ra:   o1.Ra,
		Dec:  o1.Dec,
		Cat:  o1.Cat,
	}
}
