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

	"github.com/gin-gonic/gin"
)

func (w *Web) home(c *gin.Context) {
	ctx := c.Request.Context()
	data := newTemplateData(ctx)
	if err := w.render(c, http.StatusOK, "home.tmpl.html", data); err != nil {
		w.serverError(c, fmt.Errorf("Failed to render home template: %v", err))
	}
}

func (web *Web) testHTMX(c *gin.Context) {
	ctx := c.Request.Context()
	data := newTemplateData(ctx)
	if err := web.render(c, http.StatusOK, "htmxtest.tmpl.html", data); err != nil {
		web.serverError(c, fmt.Errorf("Failed to render htmxtest template: %v", err))
	}
}

func (web *Web) notFound(c *gin.Context) {
	ctx := c.Request.Context()
	data := newTemplateData(ctx)
	if err := web.render(c, http.StatusNotFound, "notfound.tmpl.html", data); err != nil {
		web.serverError(c, fmt.Errorf("Failed to render not found template: %v", err))
	}
}

func (web *Web) stars(c *gin.Context) {
	ctx := c.Request.Context()
	data := newTemplateData(ctx)
	if err := web.render(c, http.StatusOK, "stars.tmpl.html", data); err != nil {
		web.serverError(c, fmt.Errorf("Failed to render stars template: %v", err))
	}
}
