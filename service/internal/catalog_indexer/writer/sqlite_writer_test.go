package writer

import (
	"context"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStart(t *testing.T) {
	objects := []repository.Mastercat{
		{ID: "oid1", Ipix: 1, Ra: 1, Dec: 1, Cat: "vlass"},
	}
	repo := &mocks.Repository{}
	repo.On("InsertObject", mock.Anything, masterCat2InsertParams(objects[0])).Return(objects[0], nil)
	writer := NewSqliteWriter(repo, make(chan indexer.IndexerResult), make(chan bool), context.Background())

	writer.Start()
	writer.inbox <- indexer.IndexerResult{Objects: objects}
	close(writer.inbox)
	<-writer.done

	repo.AssertExpectations(t)
}
