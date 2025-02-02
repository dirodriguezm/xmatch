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
	id := "1"
	ra := 1.0
	dec := 1.0
	ipix := int64(1)
	cat := "test"
	mastercat := repository.ParquetMastercat{
		ID:   &id,
		Ipix: &ipix,
		Ra:   &ra,
		Dec:  &dec,
		Cat:  &cat,
	}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params, err := row2insertParams(mastercat)
	require.NoError(t, err)
	repo := &mocks.Repository{}
	repo.On("InsertObject", mock.Anything, params).Return(
		repository.Mastercat{ID: id, Ra: ra, Dec: dec, Ipix: ipix, Cat: cat},
		nil,
	)
	writer := NewSqliteWriter(repo, make(chan indexer.WriterInput[repository.ParquetMastercat]), make(chan bool), context.Background(), src)

	writer.Start()
	writer.BaseWriter.InboxChannel <- indexer.WriterInput[repository.ParquetMastercat]{
		Rows: []repository.ParquetMastercat{mastercat},
	}
	close(writer.BaseWriter.InboxChannel)
	<-writer.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}
