package web

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (web *Web) localize() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("aca")
		ctx := c.Request.Context()
		acceptLang := c.Request.Header.Get("Accept-Language")
		if lang := c.Request.URL.Query().Get("lang"); lang != "" {
			slog.Debug("localize","lang",lang)
			acceptLang = lang
		}

		localizer := i18n.NewLocalizer(translations, acceptLang)

		c.Request = c.Request.WithContext(context.WithValue(ctx, Localizer, localizer))
		slog.Debug("aca")
		c.Next()
	}
}

// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// traceId := app.service.ReadGCPTraceID(r)

// ctx := context.WithValue(r.Context(), request.TraceID, traceId)
// next.ServeHTTP(w, r.WithContext(ctx))
// })
// }
