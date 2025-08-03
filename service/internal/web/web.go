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
	"html/template"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
	"github.com/go-playground/form/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Web struct {
	getenv           func(string) string
	config           *config.ServiceConfig
	conesearch       *conesearch.ConesearchService
	metadata         *metadata.MetadataService
	templateCache    map[string]*template.Template
	translations     *i18n.Bundle
	defaultLocalizer *i18n.Localizer
	formDecoder      *form.Decoder
}

func New(
	conesearch *conesearch.ConesearchService,
	metadata *metadata.MetadataService,
	config *config.ServiceConfig,
	getenv func(string string) string,
) (*Web, error) {
	if conesearch == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	if metadata == nil {
		return nil, fmt.Errorf("MetadataService was nil while creating HttpServer")
	}
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("err creating template cache: %v", err)
	}

	w := &Web{
		getenv:        getenv,
		config:        config,
		conesearch:    conesearch,
		metadata:      metadata,
		templateCache: templateCache,
	}

	if err := w.loadTranslations(); err != nil {
		return nil, fmt.Errorf("Failed to load translations: %v", err)
	}

	return w, nil
}

// func TestWeb(t *testing.T) (Web, *strings.Builder) {
// t.Helper()
// stdout := &strings.Builder{}

// getenv := func(key string) string {
// switch key {
// case "SERVICE_PORT"
// return "8080"

// c

// w := &Web{

// }
