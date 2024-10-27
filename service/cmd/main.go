package main

import (
	"xmatch/service/internal/core"
	httpservice "xmatch/service/internal/http_service"
	"xmatch/service/internal/repository"
	"xmatch/service/pkg/container"

	"github.com/dirodriguezm/healpix"
)

func AppContainer() container.Container {
	ctr := container.NewContainer()
	err := ctr.Register("sqliteRepository", func() core.Repository {
		return &repository.SqliteRepository{}
	})
	if err != nil {
		panic("could not register repository")
	}
	err = ctr.Register("conesearchService", func(r core.Repository) (*core.ConesearchService, error) {
		return core.NewConesearchService(
			core.WithNside(18),
			core.WithScheme(healpix.Nest),
			core.WithCatalog("vlass"),
			core.WithRepository(r),
		)
	})
	if err != nil {
		panic("could not register conesearch service")
	}
	err = ctr.Register("httpServer", func(service core.ConesearchService) (httpservice.HttpServer, error) {
		return httpservice.NewHttpServer(&service), nil
	})
	if err != nil {
		panic("could not register http server")
	}
	return ctr
}

func main() {
	appContainer := AppContainer()
	var server httpservice.HttpServer
	err := appContainer.ResolveWithBinds("httpServer", &server, []string{"conesearchService", "sqliteRepository"})
	if err != nil {
		panic(err)
	}
}
