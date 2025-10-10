// Copyright 2024-2025 Diego Rodriguez Mancini
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

package api

import (
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
)

type API struct {
	conesearchService *conesearch.ConesearchService
	metadataService   *metadata.MetadataService
	lightcurveService *lightcurve.LightcurveService
	config            *config.ServiceConfig
	getenv            func(string) string
}

func New(
	conesearchService *conesearch.ConesearchService,
	metadataService *metadata.MetadataService,
	lightcurveService *lightcurve.LightcurveService,
	config *config.ServiceConfig,
	getenv func(string) string,
) (*API, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	if metadataService == nil {
		return nil, fmt.Errorf("MetadataService was nil while creating HttpServer")
	}
	if lightcurveService == nil {
		return nil, fmt.Errorf("LightcurveService was nil while creating HttpServer")
	}
	return &API{conesearchService, metadataService, lightcurveService, config, getenv}, nil
}
