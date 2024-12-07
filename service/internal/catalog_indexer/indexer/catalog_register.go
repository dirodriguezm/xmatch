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
