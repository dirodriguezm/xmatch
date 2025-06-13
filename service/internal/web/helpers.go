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
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	// "runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func (web *Web) serverError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL
		// trace  = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "uri", uri)

	data := web.newTemplateData(c)
	web.render(c, http.StatusNotFound, "error.tmpl.html", data)
}

func (web *Web) notFound(c *gin.Context) {
	data := web.newTemplateData(c)
	web.render(c, http.StatusNotFound, "notfound.tmpl.html", data)
}

func (web *Web) render(c *gin.Context, status int, page string, data templateData) {
	ts, ok := web.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		web.serverError(c, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		web.serverError(c, err)
		return
	}

	c.Data(status, "text/html; charset=utf-8", buf.Bytes())
}

func (web *Web) newTemplateData(c *gin.Context) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		// Flash:           web.sessionManager.PopString(r.Context(), "flash"),
		// IsAuthenticated: web.isAuthenticated(r),
		// CSRFToken:       nosurf.Token(r),
	}
}
