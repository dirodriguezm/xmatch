package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/app"

	"github.com/gin-gonic/gin"
)

// @title						CrossWave HTTP API
// @version					1.0
// @description.markdown		api
// @host						localhost:8080
// @BasePath		/v1
// @contact.name	Diego Rodriguez Mancini
// @contact.email	diegorodriguezmancini@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @tag.name			conesearch
// @tag.description	Search for astronomical objects by celestial coordinates (RA/Dec)
//
// @tag.name			metadata
// @tag.description	Retrieve detailed catalog information for specific objects
//
// @tag.name			lightcurve
// @tag.description	Get time-series photometry data for objects
//
// @externalDocs.description	ALeRCE Documentation
// @externalDocs.url			https://alerce.science
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

	repo := app.ServiceRepository(db)

	conesearchService, err := app.ConesearchService(repo)
	if err != nil {
		return fmt.Errorf("creating conesearch service: %w", err)
	}

	metadataService, err := app.MetadataService(repo)
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
