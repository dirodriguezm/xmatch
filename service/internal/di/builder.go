package di

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"xmatch/service/internal/core"
	httpservice "xmatch/service/internal/http_service"
	"xmatch/service/internal/repository"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func ContainerBuilder() {
	container.MustSingleton(container.Global, func() *sql.DB {
		defaultConn := ":memory:"
		conn, ok := os.LookupEnv("DB_CONN")
		if !ok {
			conn = defaultConn
		}
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			panic(err)
		}
		return db
	})
	container.MustSingleton(container.Global, func(db *sql.DB) core.Repository {
		return repository.New(db)
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
