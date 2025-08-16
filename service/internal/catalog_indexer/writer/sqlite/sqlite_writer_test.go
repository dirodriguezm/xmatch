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
	"database/sql"
	"fmt"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"

	"github.com/dirodriguezm/xmatch/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStart_Mastercat(t *testing.T) {
	mastercat := repository.Mastercat{
		ID:   "1",
		Ipix: int64(1),
		Ra:   1.0,
		Dec:  1.0,
		Cat:  "test",
	}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params := repository.InsertObjectParams{
		ID:   mastercat.ID,
		Ipix: mastercat.Ipix,
		Ra:   mastercat.Ra,
		Dec:  mastercat.Dec,
		Cat:  mastercat.Cat,
	}
	repo := &mocks.Repository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertObject",
		mock.Anything,
		mock.Anything,
		[]repository.InsertObjectParams{params},
	).Return(nil)
	w := NewSqliteWriter(repo, make(chan writer.WriterInput[repository.Mastercat]), make(chan struct{}), context.Background(), src)

	w.Start()
	w.BaseWriter.InboxChannel <- writer.WriterInput[repository.Mastercat]{
		Rows: []repository.Mastercat{mastercat},
	}
	close(w.BaseWriter.InboxChannel)
	<-w.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}

func TestStart_Allwise(t *testing.T) {
	allwise := repository.Allwise{
		ID:        "test",
		W1mpro:    sql.NullFloat64{Float64: 1.0, Valid: true},
		W1sigmpro: sql.NullFloat64{Float64: 1.0, Valid: true},
		W2mpro:    sql.NullFloat64{Float64: 2.0, Valid: true},
		W2sigmpro: sql.NullFloat64{Float64: 2.0, Valid: true},
	}
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()

	params := repository.InsertAllwiseParams{
		ID:         allwise.ID,
		W1mpro:     allwise.W1mpro,
		W1sigmpro:  allwise.W1sigmpro,
		W2mpro:     allwise.W2mpro,
		W2sigmpro:  allwise.W2sigmpro,
		W3mpro:     allwise.W3mpro,
		W3sigmpro:  allwise.W3sigmpro,
		W4mpro:     allwise.W4mpro,
		W4sigmpro:  allwise.W4sigmpro,
		JM2mass:    allwise.JM2mass,
		JMsig2mass: allwise.JMsig2mass,
		HM2mass:    allwise.HM2mass,
		HMsig2mass: allwise.HMsig2mass,
		KM2mass:    allwise.KM2mass,
		KMsig2mass: allwise.KMsig2mass,
	}
	repo := &mocks.Repository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertAllwise",
		mock.Anything,
		mock.Anything,
		[]repository.InsertAllwiseParams{params},
	).Return(nil)
	w := NewSqliteWriter(
		repo,
		make(chan writer.WriterInput[repository.Allwise]),
		make(chan struct{}),
		context.Background(),
		src,
	)

	w.Start()
	w.BaseWriter.InboxChannel <- writer.WriterInput[repository.Allwise]{
		Rows:  []repository.Allwise{allwise},
		Error: nil,
	}
	close(w.BaseWriter.InboxChannel)
	<-w.BaseWriter.DoneChannel

	repo.AssertExpectations(t)
}
