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
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (w *Web) serverError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL
		trace  = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	ctx := c.Request.Context()
	data := w.newTemplateData(ctx)

	err = w.render(c, http.StatusInternalServerError, "error.tmpl.html", data)
	if err != nil {
		slog.Error("Failed to render error template", "error", err.Error())
	}
}

func (web *Web) render(c *gin.Context, status int, page string, data templateData) error {
	ts, ok := web.templateCache[page]
	if !ok {
		return fmt.Errorf("the template %s does not exist", page)
	}

	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	c.Data(status, "text/html; charset=utf-8", buf.Bytes())
	return nil
}
