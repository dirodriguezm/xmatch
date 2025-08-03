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
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dirodriguezm/xmatch/service/ui"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (w *Web) localize() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		acceptLang := c.Request.Header.Get("Accept-Language")
		if lang := c.Request.URL.Query().Get("lang"); lang != "" {
			acceptLang = lang
		}

		localizer := i18n.NewLocalizer(w.translations, acceptLang)

		c.Request = c.Request.WithContext(context.WithValue(ctx, Localizer, localizer))
		slog.Debug("localize", "lang", acceptLang)
		c.Next()
	}
}

func extractRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		route := c.Request.URL.Path

		routes := map[string]string{
			"/":            "home",
			"/user/signup": "signup",
			"/user/login":  "login",
		}

		selected := routes[route]

		c.Request = c.Request.WithContext(context.WithValue(ctx, Route, selected))
		c.Next()
	}
}

// maxSpeed is the maximum speed for the stars animation.
// it is used to calculate the delays for the stars animation.
// the value is in seconds, so 600 means 600 seconds (10 minutes).
const maxSpeed = 600 // this would be 600s in animation-delay style

// calculateStarsDelay calculates the delay for the stars animation based on the current time.
// it uses the current Unix timestamp to create a pseudo-random delay for each speed category.
// the delays are negative to ensure that the animation starts immediately.
func calculateStarsDelay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		seconds := time.Now().Unix()

		r := seconds % maxSpeed

		delays := delay{
			Fast:   -int16(maxSpeed + r),
			Medium: -int16(maxSpeed*2 + r),
			Slow:   -int16(maxSpeed*4 + r),
		}

		c.Request = c.Request.WithContext(context.WithValue(ctx, Delays, delays))
		c.Next()
	}
}

// cacheControl sets the Cache-Control and Expires headers for static assets.
// it allows the browser to cache the assets for 30 days (2592000 seconds).
func cacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		file := strings.Split(c.Request.URL.Path, "static")
		if len(file) < 2 {
			c.Next()
			return
		}

		filename := fmt.Sprintf("static%s", file[1])

		slog.Debug("cacheControl", "file", filename)
		data, err := ui.Files.ReadFile(filename)
		slog.Debug("cacheControl", "fileSize", len(data), "error", err)
		if err != nil {
			c.Next()
			return
		}

		etag := fmt.Sprintf("%x", sha256.Sum256(data))

		c.Header("ETag", etag)

		if match := c.Request.Header.Get("If-None-Match"); match == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}

		c.Next()
	}
}
