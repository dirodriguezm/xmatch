package di

import (
	"xmatch/service/internal/core"
	httpservice "xmatch/service/internal/http_service"
	"xmatch/service/internal/repository"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func ContainerBuilder() {
	container.MustSingleton(container.Global, func() core.Repository {
		return &repository.SqliteRepository{}
	})
	container.MustSingleton(container.Global, func(r core.Repository) (*core.ConesearchService, error) {
		return core.NewConesearchService(
			core.WithNside(18),
			core.WithScheme(healpix.Nest),
			core.WithCatalog("vlass"),
			core.WithRepository(r),
		)
	})
	container.MustSingleton(container.Global, func(service *core.ConesearchService) httpservice.HttpServer {
		return httpservice.NewHttpServer(service)
	})
}
