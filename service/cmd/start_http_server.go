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
