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
	Done       chan bool
}

func New(repository conesearch.Repository, ch chan indexer.IndexerResult, ctx context.Context) *SqliteWriter {
	return &SqliteWriter{
		inbox:      ch,
		repository: repository,
		ctx:        ctx,
		Done:       make(chan bool),
	}
}

func (w *SqliteWriter) Start() {
	go func() {
		defer func() {
			w.Done <- true
			close(w.Done)
		}()
		for msg := range w.inbox {
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
	}()
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
