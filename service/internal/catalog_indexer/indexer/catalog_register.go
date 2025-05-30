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

package indexer

import (
	"context"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type Catalog struct {
	Name  string
	Nside int64
}

type CatalogRegister struct {
	repo conesearch.Repository
	ctx  context.Context
	cfg  config.SourceConfig
}

func NewCatalogRegister(
	ctx context.Context,
	repo conesearch.Repository,
	cfg config.SourceConfig,
) *CatalogRegister {
	return &CatalogRegister{
		repo: repo,
		ctx:  ctx,
		cfg:  cfg,
	}
}

func (r *CatalogRegister) RegisterCatalog() {
	r.repo.InsertCatalog(r.ctx, repository.InsertCatalogParams{
		Name:  r.cfg.CatalogName,
		Nside: int64(r.cfg.Nside),
	})
}
