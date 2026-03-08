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

package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/neowise"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dirodriguezm/healpix"
)

func ServiceLogger(getenv func(string) string, stdout io.Writer) *slog.Logger {
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
}

func ServiceDatabase(cfg config.Config) (*sql.DB, error) {
	conn := cfg.Service.Database.Url
	if !strings.Contains(conn, "?") {
		conn += "?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
	}
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, fmt.Errorf("could not create sqlite connection: %w", err)
	}
	_, err = db.Exec("select 'test conn'")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	return db, nil
}

func ServiceRepository(db *sql.DB) conesearch.Repository {
	return repository.New(db)
}

func ConesearchService(repo conesearch.Repository) (*conesearch.ConesearchService, error) {
	ctx := context.Background()
	catalogs, err := repo.GetCatalogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not find catalogs in DB when creating conesearch service: %w", err)
	}

	con, err := conesearch.NewConesearchService(
		conesearch.WithScheme(healpix.Nest),
		conesearch.WithRepository(repo),
		conesearch.WithCatalogs(catalogs),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create ConesearchService: %w", err)
	}
	return con, nil
}

func MetadataService(repo conesearch.Repository) (*metadata.MetadataService, error) {
	service, err := metadata.NewMetadataService(repo)
	if err != nil {
		return nil, fmt.Errorf("could not create MetadataService: %w", err)
	}
	return service, nil
}

func LightcurveService(cfg config.Config, conesearchService *conesearch.ConesearchService) (*lightcurve.LightcurveService, error) {
	lightcurveFilterSlice := []lightcurve.LightcurveFilter{lightcurve.DummyLightcurveFilter}
	if cfg.Service.LightcurveServiceConfig.NeowiseConfig.UseCntrFilter {
		lightcurveFilterSlice = append(lightcurveFilterSlice, neowise.Filter)
	}

	service, err := lightcurve.New(
		[]lightcurve.ExternalClient{neowise.NewNeowiseClient()},
		lightcurveFilterSlice,
		conesearchService,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create LightcurveService: %w", err)
	}
	return service, nil
}

func API(conesearchService *conesearch.ConesearchService, metadataService *metadata.MetadataService, lightcurveService *lightcurve.LightcurveService, cfg config.ServiceConfig, getenv func(string) string) (*api.API, error) {
	return api.New(conesearchService, metadataService, lightcurveService, cfg, getenv)
}
