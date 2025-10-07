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
	"errors"
	"net/http"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/gin-gonic/gin"
)

// Search for objects in a given region using multiple coordinates
//
//	@Summary		Search for objects in a given region using multiple coordinates
//	@Description	Search for objects in a given region using list of ra, dec and a single radius
//	@Tags			conesearch
//	@Accept			json
//	@Produce		json
//
//	@Param			ra			body		[]float64	true	"Right ascension in degrees"
//	@Param			dec			body		[]float64	true	"Declination in degrees"
//	@Param			radius		body		float64		true	"Radius in degrees"
//	@Param			catalog		body		string		false	"Catalog to search in"
//	@Param			nneighbor	body		int			false	"Number of neighbors to return"
//
//	@Success		200			{array}		repository.Mastercat
//	@Success		204			{string}	string
//	@Failure		400			{object}	conesearch.ValidationError
//	@Failure		500			{string}	string
//	@Router			/bulk-conesearch [post]
func (api *API) conesearchBulk(c *gin.Context) {
	var bulkRequest BulkConesearchRequest
	if err := c.ShouldBindJSON(&bulkRequest); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if bulkRequest.Nneighbor == 0 {
		bulkRequest.Nneighbor = 1
	}
	if bulkRequest.Catalog == "" {
		bulkRequest.Catalog = "all"
	}

	result, err := api.conesearchService.BulkConesearch(
		bulkRequest.Ra,
		bulkRequest.Dec,
		bulkRequest.Radius,
		bulkRequest.Nneighbor,
		bulkRequest.Catalog,
		api.config.BulkChunkSize,
		api.config.MaxBulkConcurrency,
	)
	if err != nil {
		if errors.As(err, &conesearch.ValidationError{}) {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Could not execute conesearch")
		}
		return
	}
	if len(result) == 0 {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Search for objects in a given region
//
//		@Summary		Search for objects in a given region
//		@Description	Search for objects in a given region using ra, dec and radius
//		@Tags			conesearch
//		@Accept			json
//		@Produce		json
//		@Param			ra			query		string	true	"Right ascension in degrees"
//		@Param			dec			query		string	true	"Declination in degrees"
//		@Param			radius		query		string	true	"Radius in degrees"
//		@Param			catalog		query		string	false	"Catalog to search in"
//		@Param			nneighbor	query		string	false	"Number of neighbors to return"
//	 @Param			getMetadata	query		string	false	"Return metadata results"
//		@Success		200			{array}		repository.Mastercat
//		@Success		204			{string}	string
//		@Failure		400			{object}	conesearch.ValidationError
//		@Failure		500			{string}	string
//		@Router			/conesearch [get]
func (api *API) conesearch(c *gin.Context) {
	ra := c.Query("ra")
	dec := c.Query("dec")
	radius := c.Query("radius")
	catalog := c.DefaultQuery("catalog", "all")
	nneighbor := c.DefaultQuery("nneighbor", "1")
	getMetadata := c.DefaultQuery("getMetadata", "false")

	parsedRa, err := parseRa(ra)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	parsedDec, err := parseDec(dec)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	parsedRadius, err := parseRadius(radius)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	parsedNneighbor, err := parseNneighbor(nneighbor)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var result any
	var serviceErr error
	if getMetadata == "true" {
		result, serviceErr = api.conesearchService.FindMetadataByConesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor, catalog)
	} else {
		result, serviceErr = api.conesearchService.Conesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor, catalog)
	}

	if serviceErr != nil {
		if errors.As(serviceErr, &conesearch.ValidationError{}) {
			c.JSON(http.StatusBadRequest, serviceErr)
		} else {
			c.Error(serviceErr)
			c.JSON(http.StatusInternalServerError, "Could not execute conesearch")
		}
		return
	}
	if isEmptyResult(result) {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, result)
}

func isEmptyResult(result any) bool {
	// WARN: in this layer there should not be any references to repository, that should be done in service
	switch result := result.(type) {
	case []repository.Mastercat:
		return len(result) == 0
	case []repository.Allwise:
		return len(result) == 0
	default:
		return false
	}
}
