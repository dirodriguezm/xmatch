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
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (api *API) SetupRoutes(r *gin.Engine) {
	if r == nil {
		panic("api: gin engine cannot be nil")
	}
	r.Use(gin.Recovery())
	if os.Getenv("USE_LOGGER") != "" {
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
