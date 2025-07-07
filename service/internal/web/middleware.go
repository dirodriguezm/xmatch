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

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func localize() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		acceptLang := c.Request.Header.Get("Accept-Language")
		if lang := c.Request.URL.Query().Get("lang"); lang != "" {
			acceptLang = lang
		}

		localizer := i18n.NewLocalizer(translations, acceptLang)

		c.Request = c.Request.WithContext(context.WithValue(ctx, Localizer, localizer))
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
