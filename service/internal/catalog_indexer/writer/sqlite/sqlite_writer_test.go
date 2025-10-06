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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"

	"github.com/stretchr/testify/mock"
)

func TestReceive_Mastercat(t *testing.T) {
	mastercat := repository.Mastercat{
		ID:   "1",
		Ipix: int64(1),
		Ra:   1.0,
		Dec:  1.0,
		Cat:  "test",
	}

	repo := &conesearch.MockRepository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertObject",
		mock.Anything,
		mock.Anything,
		[]any{mastercat},
	).Return(nil)

	w := New(repo, context.Background(), "mastercat")
	w.Write(nil, actor.Message{Rows: []any{mastercat}, Error: nil})

	repo.AssertExpectations(t)
}

func TestReceive_Allwise(t *testing.T) {
	allwise := repository.Allwise{
		ID:        "test",
		W1mpro:    sql.NullFloat64{Float64: 1.0, Valid: true},
		W1sigmpro: sql.NullFloat64{Float64: 1.0, Valid: true},
		W2mpro:    sql.NullFloat64{Float64: 2.0, Valid: true},
		W2sigmpro: sql.NullFloat64{Float64: 2.0, Valid: true},
	}

	repo := &conesearch.MockRepository{}
	repo.On("GetDbInstance").Return(nil)
	repo.On(
		"BulkInsertAllwise",
		mock.Anything,
		mock.Anything,
		[]any{allwise},
	).Return(nil)

	w := New(repo, context.Background(), "allwise")
	w.Write(nil, actor.Message{Rows: []any{allwise}, Error: nil})

	repo.AssertExpectations(t)
}
