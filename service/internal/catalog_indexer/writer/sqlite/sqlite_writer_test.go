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

func TestStart_Mastercat(t *testing.T) {
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
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertObject",
		mock.Anything,
		mock.Anything,
		[]repository.InsertObjectParams{params.(repository.InsertObjectParams)},
	).Return(nil)
	writer := NewSqliteWriter(repo, make(chan indexer.WriterInput[repository.ParquetMastercat]), make(chan bool), context.Background(), src)

	writer.Start()
	writer.BaseWriter.InboxChannel <- indexer.WriterInput[repository.ParquetMastercat]{
		Rows: []repository.ParquetMastercat{mastercat},
	}
	close(writer.BaseWriter.InboxChannel)
	<-writer.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}

func TestStart_Allwise(t *testing.T) {
	designation := "test"
	w1mpro := 1.0
	w1sigmpro := 1.0
	w2mpro := 2.0
	w2sigmpro := 2.0
	allwise := repository.AllwiseMetadata{
		Designation: &designation,
		W1mpro:      &w1mpro,
		W1sigmpro:   &w1sigmpro,
		W2mpro:      &w2mpro,
		W2sigmpro:   &w2sigmpro,
	}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params, err := row2insertParams(allwise)
	require.NoError(t, err)
	repo := &mocks.Repository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertAllwise",
		mock.Anything,
		mock.Anything,
		[]repository.InsertAllwiseParams{params.(repository.InsertAllwiseParams)},
	).Return(nil)
	writer := NewSqliteWriter(
		repo,
		make(chan indexer.WriterInput[repository.AllwiseMetadata]),
		make(chan bool),
		context.Background(),
		src,
	)

	writer.Start()
	writer.BaseWriter.InboxChannel <- indexer.WriterInput[repository.AllwiseMetadata]{
		Rows:  []repository.AllwiseMetadata{allwise},
		Error: nil,
	}
	close(writer.BaseWriter.InboxChannel)
	<-writer.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}
