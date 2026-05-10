package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"

	"github.com/gin-gonic/gin"
)

// @title			CrossWave HTTP API
// @version		1.0
// @description	API for the CrossWave Xmatch service. This service allows to search for objects in a given region and to retrieve metadata from the catalogs.
// @host			localhost:8080
// @BasePath		/v1
// @contact.name	Diego Rodriguez Mancini
// @contact.email	diegorodriguezmancini@gmail.com
func StartHttpServer(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	cfg, err := app.Config(getenv)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		return fmt.Errorf("creating database connection: %w", err)
	}
	defer db.Close()

	queries := app.ServiceRepository(db)

	resolver := catalog.NewResolver()
	resolver.RegisterStore("allwise", queries)
	resolver.RegisterStore("gaia", queries)
	resolver.RegisterStore("erosita", queries)

	conesearchService, err := app.ConesearchService(queries, resolver)
	if err != nil {
		return fmt.Errorf("creating conesearch service: %w", err)
	}

	metadataService, err := app.MetadataService(resolver)
	if err != nil {
		return fmt.Errorf("creating metadata service: %w", err)
	}

	lightcurveService, err := app.LightcurveService(cfg, conesearchService)
	if err != nil {
		return fmt.Errorf("creating lightcurve service: %w", err)
	}

	api, err := app.API(conesearchService, metadataService, lightcurveService, cfg.Service, getenv)
	if err != nil {
		return fmt.Errorf("creating API: %w", err)
	}

	r := gin.New()
	if getenv("USE_LOGGER") != "" {
		r.Use(func(c *gin.Context) {
			slog.Info("request", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.Next()
		})
	}

	api.SetupRoutes(r)

	err = r.Run()
	return err
}
