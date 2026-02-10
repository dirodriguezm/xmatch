package main

import (
	"context"
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/web"
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
	ctr := di.BuildServiceContainer(ctx, getenv, stdout)
	var api *api.API
	var web *web.Web
	ctr.Resolve(&api)
	ctr.Resolve(&web)

	r := gin.New()
	r.Use(gin.Recovery())
	if getenv("USE_LOGGER") != "" {
		r.Use(func(c *gin.Context) {
			slog.Info("request", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.Next()
		})
	}
	r.SetTrustedProxies([]string{"localhost"})

	api.SetupRoutes(r)
	web.SetupRoutes(r)

	err := r.Run()
	return err
}
