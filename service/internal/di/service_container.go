package di

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func BuildServiceContainer() container.Container {
	ctr := container.New()

	// read config
	ctr.Singleton(func() *config.Config {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		return cfg
	})

	ctr.Singleton(func() *slog.LevelVar {
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

	// Register DB
	ctr.Singleton(func(cfg *config.Config) *sql.DB {
		conn := cfg.Service.Database.Url
		db, err := sql.Open("sqlite3", conn)
		if err != nil {
			slog.Error("Could not create sqlite3 connection", "conn", conn)
			panic(err)
		}
		_, err = db.Exec("select 'test conn'")
		if err != nil {
			slog.Error("Could not connect to database", "conn", conn)
			panic(err)
		}
		slog.Debug("Created database", "conn", conn)
		return db
	})

	ctr.Singleton(func(db *sql.DB) conesearch.Repository {
		return repository.New(db)
	})

	ctr.Singleton(func(r conesearch.Repository, cfg *config.Config) *conesearch.ConesearchService {
		ctx := context.TODO()
		catalogs, err := r.GetCatalogs(ctx)
		if err != nil {
			slog.Error("Could not find catalogs in DB when creating conesearch service", "error", err)
			panic(err)
		}

		con, err := conesearch.NewConesearchService(
			conesearch.WithScheme(healpix.Nest),
			conesearch.WithRepository(r),
			conesearch.WithCatalogs(catalogs),
		)
		if err != nil {
			slog.Error("Could not register ConesearchService")
			panic(err)
		}
		return con
	})

	ctr.Singleton(func(service *conesearch.ConesearchService) *httpservice.HttpServer {
		server, err := httpservice.NewHttpServer(service)
		if err != nil {
			slog.Error("Could not register HttpServer")
			panic(err)
		}
		if server == nil {
			panic("Server nil while registering HttpServer")
		}
		return server
	})
	return ctr
}