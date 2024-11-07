package di

import (
	"database/sql"
	"github.com/dirodriguezm/xmatch/service/internal/core/conesearch"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/pkg/repository"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func ContainerBuilder() {
	container.MustSingleton(container.Global, func() *slog.LevelVar {
		levels := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"info":  slog.LevelInfo,
			"error": slog.LevelError,
			"warn":  slog.LevelWarn,
			"":      slog.LevelInfo,
		}
		var programLevel = new(slog.LevelVar)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
		slog.SetDefault(logger)
		programLevel.Set(levels[os.Getenv("LOG_LEVEL")])
		return programLevel
	})

	container.MustSingleton(container.Global, func() *sql.DB {
		conn := os.Getenv("DB_CONN")
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite3 connection", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})

	container.MustSingleton(container.Global, func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})

	container.MustSingleton(container.Global, func(r conesearch.Repository) (*conesearch.ConesearchService, error) {
		return conesearch.NewConesearchService(
			conesearch.WithNside(18),
			conesearch.WithScheme(healpix.Nest),
			conesearch.WithCatalog("vlass"),
			conesearch.WithRepository(r),
		)
	})

	container.MustSingleton(container.Global, func(service *conesearch.ConesearchService) httpservice.HttpServer {
		return httpservice.NewHttpServer(service)
	})
}
