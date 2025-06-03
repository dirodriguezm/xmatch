package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (api *API) SetupRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())
	if api.getEnv("USE_LOGGER") != "" {
		r.Use(gin.Logger())
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/conesearch", api.conesearch)
		v1.POST("/bulk-conesearch", api.conesearchBulk)
		v1.GET("/metadata", api.metadata)
		v1.POST("/bulk-metadata", api.metadataBulk)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.SetTrustedProxies([]string{"localhost"})
}
