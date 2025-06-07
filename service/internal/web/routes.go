// Copyright 2024-2025 Matías Medina Silva
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

package web

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (web *Web) SetupRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())
	if os.Getenv("USE_LOGGER") != "" {
		r.Use(gin.Logger())
	}

	// Simple route that responds with HTML
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Placeholder Page</title>
				<style>
					body { 
						font-family: Arial, sans-serif; 
						text-align: center; 
						padding: 50px; 
					}
					h1 { color: #333; }
				</style>
			</head>
			<body>
				<h1>Website Coming Soon</h1>
				<p>This is a placeholder page. The real website will be here soon!</p>
			</body>
			</html>
		`))
	})

	r.SetTrustedProxies([]string{"localhost"})
}
