package sqlite_writer

import (
	"context"
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/dirodriguezm/xmatch/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	objects := []indexer.Row{
		{"oid": "oid1", "ipix": 1, "ra": 1, "dec": 1, "cat": "vlass"},
	}
	mastercat := repository.Mastercat{ID: "oid1", Ipix: 1, Ra: 1, Dec: 1, Cat: "vlass"}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params, err := row2insertParams(objects[0], src)
	require.NoError(t, err)
	repo := &mocks.Repository{}
	repo.On("InsertObject", mock.Anything, params).Return(mastercat, nil)
	writer := NewSqliteWriter(repo, make(chan indexer.WriterInput), make(chan bool), context.Background(), src)

	writer.Start()
	writer.BaseWriter.Inbox <- indexer.WriterInput{Rows: objects}
	close(writer.BaseWriter.Inbox)
	<-writer.BaseWriter.Done

	repo.AssertExpectations(t)
}
