package api

import (
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
)

type API struct {
	conesearchService *conesearch.ConesearchService
	metadataService   *metadata.MetadataService
	config            *config.ServiceConfig
	getEnv            func(string) string
}

func New(
	conesearchService *conesearch.ConesearchService,
	metadataService *metadata.MetadataService,
	config *config.ServiceConfig,
	getEnv func(string) string,
) (*API, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	if metadataService == nil {
		return nil, fmt.Errorf("MetadataService was nil while creating HttpServer")
	}
	return &API{conesearchService, metadataService, config, getEnv}, nil
}
