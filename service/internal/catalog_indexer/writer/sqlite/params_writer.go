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

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type ParamsWriter[T any] interface {
	BulkWrite([]T) error
}

type MastercatWriter struct {
	Repo conesearch.Repository
	Ctx  context.Context
	Db   *sql.DB
}

func (w MastercatWriter) BulkWrite(objs []repository.InsertObjectParams) error {
	return w.Repo.BulkInsertObject(w.Ctx, w.Db, objs)
}

type AllwiseWriter struct {
	Repo conesearch.Repository
	Ctx  context.Context
	Db   *sql.DB
}

func (w AllwiseWriter) BulkWrite(objs []repository.Metadata) error {
	return w.Repo.BulkInsertAllwise(w.Ctx, w.Db, objs)
}
