// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlite_writer

import (
	"context"
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/dirodriguezm/xmatch/service/mocks"
	"github.com/stretchr/testify/mock"
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

	params := mastercat.ToInsertParams()
	repo := &mocks.Repository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertObject",
		mock.Anything,
		mock.Anything,
		[]repository.InsertObjectParams{params.(repository.InsertObjectParams)},
	).Return(nil)
	w := NewSqliteWriter(repo, make(chan writer.WriterInput[repository.ParquetMastercat]), make(chan bool), context.Background(), src)

	w.Start()
	w.BaseWriter.InboxChannel <- writer.WriterInput[repository.ParquetMastercat]{
		Rows: []repository.ParquetMastercat{mastercat},
	}
	close(w.BaseWriter.InboxChannel)
	<-w.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}

func TestStart_Allwise(t *testing.T) {
	source_id := "test"
	w1mpro := 1.0
	w1sigmpro := 1.0
	w2mpro := 2.0
	w2sigmpro := 2.0
	allwise := repository.AllwiseMetadata{
		Source_id: &source_id,
		W1mpro:    &w1mpro,
		W1sigmpro: &w1sigmpro,
		W2mpro:    &w2mpro,
		W2sigmpro: &w2sigmpro,
	}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params := allwise.ToInsertParams()
	repo := &mocks.Repository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertAllwise",
		mock.Anything,
		mock.Anything,
		[]repository.InsertAllwiseParams{params.(repository.InsertAllwiseParams)},
	).Return(nil)
	w := NewSqliteWriter(
		repo,
		make(chan writer.WriterInput[repository.AllwiseMetadata]),
		make(chan bool),
		context.Background(),
		src,
	)

	w.Start()
	w.BaseWriter.InboxChannel <- writer.WriterInput[repository.AllwiseMetadata]{
		Rows:  []repository.AllwiseMetadata{allwise},
		Error: nil,
	}
	close(w.BaseWriter.InboxChannel)
	<-w.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}
