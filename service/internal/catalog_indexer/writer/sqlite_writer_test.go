package writer

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

func TestStart(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "oid1", Ipix: 1, Ra: 1, Dec: 1, Cat: "vlass"},
	}
	repo := &conesearch.MockRepository{Objects: objects, Error: nil}
	repo.On("InsertObject", masterCat2InsertParams(objects[0])).Return(objects[0], nil)
	writer := NewSqliteWriter(repo, make(chan indexer.IndexerResult), context.Background())

	writer.Start()
	writer.inbox <- indexer.IndexerResult{Objects: objects}
	close(writer.inbox)
	<-writer.Done

	repo.AssertExpectations(t)
}
