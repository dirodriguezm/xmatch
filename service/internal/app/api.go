package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/neowise"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"

	_ "github.com/mattn/go-sqlite3"
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

func ServiceRepository(db *sql.DB) *repository.Queries {
	return repository.New(db)
}

func ConesearchService(queries *repository.Queries, resolver *catalog.Resolver) (*conesearch.ConesearchService, error) {
	ctx := context.Background()
	catalogs, err := queries.GetCatalogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not find catalogs in DB when creating conesearch service: %w", err)
	}

	con, err := conesearch.NewConesearchService(
		conesearch.WithScheme(healpix.Nest),
		conesearch.WithMastercatStore(queries),
		conesearch.WithResolver(resolver),
		conesearch.WithCatalogs(catalogs),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create ConesearchService: %w", err)
	}
	return con, nil
}

func MetadataService(resolver *catalog.Resolver) (*metadata.MetadataService, error) {
	service, err := metadata.NewMetadataService(resolver)
	if err != nil {
		return nil, fmt.Errorf("could not create MetadataService: %w", err)
	}
	return service, nil
}

func LightcurveService(cfg config.Config, conesearchService *conesearch.ConesearchService) (*lightcurve.LightcurveService, error) {
	neowiseFilter := lightcurve.DummyLightcurveFilter
	if cfg.Service.LightcurveServiceConfig.NeowiseConfig.UseIdFilter || cfg.Service.LightcurveServiceConfig.NeowiseConfig.UseCntrFilter {
		neowiseFilter = neowise.Filter
	}
	ztfFilter := lightcurve.DummyLightcurveFilter
	if cfg.Service.LightcurveServiceConfig.ZtfDrConfig.UseIdFilter {
		ztfFilter = ztfdr.Filter
	}
	sources := []lightcurve.Source{
		{Catalog: "neowise", Client: neowise.NewNeowiseClient(), Filter: neowiseFilter},
		{Catalog: "ztf", Client: ztfdr.NewZtfDrClient(), Filter: ztfFilter},
	}

	service, err := lightcurve.New(
		sources,
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
