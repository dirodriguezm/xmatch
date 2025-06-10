package web

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func (web *Web) serverError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL
		trace  = string(debug.Stack())
	)

	slog.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

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
