package web

import (
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
)

type Web struct {
	conesearchService *conesearch.ConesearchService
	metadataService   *metadata.MetadataService
	config            *config.ServiceConfig
	getenv            func(string) string
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
	return &Web{conesearchService, metadataService, config, getenv}, nil
}
