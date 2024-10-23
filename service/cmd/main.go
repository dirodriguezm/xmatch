package main

import (
	"xmatch/service/internal/core"
	httpservice "xmatch/service/internal/http_service"
	"xmatch/service/internal/repository"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func main() {
	err := container.Singleton(func() core.Repository {
		return &repository.SqliteRepository{}
	})
	if err != nil {
		panic("could not register repository")
	}
	err = container.Singleton(func(r core.Repository) (*core.ConesearchService, error) {
		return core.NewConesearchService(
			core.WithNside(18),
			core.WithScheme(healpix.Nest),
			core.WithCatalog("vlass"),
			core.WithRepository(r),
		)
	})
	httpservice.InitServer()
}
