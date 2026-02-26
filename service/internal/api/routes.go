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

package api

import (
	"net/http"

	"github.com/dirodriguezm/xmatch/service/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (api *API) SetupRoutes(r *gin.Engine) {
	if r == nil {
		panic("api: gin engine cannot be nil")
	}
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://xwave-rho.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	if api.getenv("USE_LOGGER") != "" {
		r.Use(gin.Logger())
	}

	docs.SwaggerInfo.Host = api.config.Host
	docs.SwaggerInfo.BasePath = api.config.BasePath

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/conesearch", api.conesearch)
		v1.POST("/bulk-conesearch", api.conesearchBulk)
		v1.GET("/metadata", api.metadata)
		v1.POST("/bulk-metadata", api.metadataBulk)
		v1.GET("/lightcurve", api.Lightcurve)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.SetTrustedProxies([]string{"localhost"})
}
