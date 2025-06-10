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

	return &Web{getenv, config, conesearchService, metadataService, templateCache}, nil
}
