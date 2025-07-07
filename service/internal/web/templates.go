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
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/dirodriguezm/xmatch/service/ui"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type templateData struct {
	CurrentYear int
	Form        any //this is for posts requests
	Local       *i18n.Localizer
	Route       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"t":         translate,
	"tc":        translateCount,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("could not load globs of html templates: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func newTemplateData(ctx context.Context) templateData {
	localizer, err := localizerFrom(ctx)
	if err != nil {
		localizer = defaultLocalizer
	}

	route, err := routeFrom(ctx)
	if err != nil {
		route = ""
	}

	return templateData{
		CurrentYear: time.Now().Year(),
		Local:       localizer,
		Route:       route,
	}
}
