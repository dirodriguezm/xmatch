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
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
	"html/template"
)

type Web struct {
	getenv            func(string) string
	config            *config.ServiceConfig
	conesearchService *conesearch.ConesearchService
	metadataService   *metadata.MetadataService
	templateCache     map[string]*template.Template
}

func New(
	conesearchService *conesearch.ConesearchService,
	metadataService *metadata.MetadataService,
	config *config.ServiceConfig,
	getenv func(string string) string,
) (*Web, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	if metadataService == nil {
		return nil, fmt.Errorf("MetadataService was nil while creating HttpServer")
	}
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("err creating template cache: %v", err)
	}
	if err := loadTranslations(); err != nil {
		return nil, fmt.Errorf("Failed to load translations: %v", err)
	}

	return &Web{getenv, config, conesearchService, metadataService, templateCache}, nil
}
