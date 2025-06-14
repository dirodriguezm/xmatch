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
	"fmt"
	"net/http"
	"time"

	"github.com/dirodriguezm/xmatch/service/ui"
	"github.com/gin-gonic/gin"
)

func (web *Web) SetupRoutes(r *gin.Engine) {
	if r == nil {
		panic("api: gin engine cannot be nil")
	}

	r.GET("/static/*filepath", func(c *gin.Context) {
		fileServer := http.FileServer(http.FS(ui.Files))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.NoRoute(web.notFound)

	r.GET("/", web.home)
	r.GET("/htmx", web.testHTMX)

	r.GET("/htmx-test", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC1123)))
	})

	r.POST("/htmx-click", func(c *gin.Context) {
		c.String(http.StatusOK, `
            <div id="click-result" class="flash">
                Button clicked at: %s
            </div>
        `, time.Now().Format(time.Kitchen))
	})
}
