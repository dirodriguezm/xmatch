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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestExtractRoute(t *testing.T) {
	tests := []struct {
		name     string
		route    string
		expected string
	}{
		{name: "Home Route", route: "/", expected: "home"},
		{name: "Signup Route", route: "/user/signup", expected: "signup"},
		{name: "Login Route", route: "/user/login", expected: "login"},
		{name: "Unknown Route", route: "/unknown", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &strings.Builder{}
			router := SetupRouter(t, stdout)

			var ctx context.Context
			router.Use(extractRoute())
			router.GET(tt.route, func(c *gin.Context) {
				ctx = c.Request.Context()
				c.String(200, "ok")
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.route, nil)

			router.ServeHTTP(w, req)

			route, ok := ctx.Value(Route).(string)
			if !ok {
				t.Fatal("could not get route from context")
			}
			Equal(t, route, tt.expected)
		})
	}

}

func TestLocalize(t *testing.T) {
	stdout := &strings.Builder{}
	router := SetupRouter(t, stdout)

	var ctx context.Context

	router.Use(localize())

	router.GET("/context", func(c *gin.Context) {
		ctx = c.Request.Context()
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/context", nil)
	router.ServeHTTP(w, req)

	loc, ok := ctx.Value(Localizer).(*i18n.Localizer)
	if !ok {
		t.Fatal("could not get localizer from context")
	}

	NotNil(t, loc)
}
