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

package di

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"

	"github.com/charmbracelet/log"
	"github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/web"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dirodriguezm/healpix"
	"github.com/golobby/container/v3"
)

func BuildServiceContainer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) container.Container {
	ctr := container.New()

	// read config
	ctr.Singleton(func() *config.Config {
		cfg, err := config.Load(getenv)
		if err != nil {
			panic(err)
		}
		return cfg
	})

	ctr.Singleton(func() *slog.Logger {
		levels := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"info":  slog.LevelInfo,
			"error": slog.LevelError,
			"warn":  slog.LevelWarn,
			"":      slog.LevelInfo,
		}
		lvl := levels[getenv("LOG_LEVEL")]
		var logger *slog.Logger
		if getenv("ENVIRONMENT") == "local" {
			handler := log.NewWithOptions(stdout, log.Options{
				Level:           log.Level(lvl.Level()),
				ReportTimestamp: true,
			})
			logger = slog.New(handler)
		} else {
			logger = slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{Level: lvl}))
		}

		slog.SetDefault(logger)
		return logger
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

	// the metadata.Repository is a subset of conesearch.Repository
	ctr.Singleton(func(repo conesearch.Repository) metadata.Repository {
		return repo
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

	ctr.Singleton(func(r metadata.Repository) *metadata.MetadataService {
		service, err := metadata.NewMetadataService(r)
		if err != nil {
			panic(fmt.Errorf("Could not create MetadataService: %w", err))
		}
		return service
	})

	ctr.Singleton(func(
		conesearchService *conesearch.ConesearchService,
		metadataService *metadata.MetadataService,
		config *config.Config,
	) *api.API {
		api, err := api.New(conesearchService, metadataService, config.Service, getenv)
		if err != nil {
			panic(fmt.Errorf("Could not register API: %w", err))
		}
		if api == nil {
			panic("api nil while registering API")
		}
		return api
	})

	ctr.Singleton(func(
		conesearchService *conesearch.ConesearchService,
		metadataService *metadata.MetadataService,
		config *config.Config,
	) *web.Web {
		web, err := web.New(conesearchService, metadataService, config.Service, getenv)
		if err != nil {
			panic(fmt.Errorf("Could not register API: %w", err))
		}
		if web == nil {
			panic("web nil while registering Web")
		}
		return web
	})

	return ctr
}
