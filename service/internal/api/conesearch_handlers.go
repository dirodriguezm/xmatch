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

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/gin-gonic/gin"
)

// Search for objects in a given region using multiple coordinates
//
//	@Summary		Bulk cone search for objects
//	@Description	Search for objects in a given region using multiple RA/Dec coordinates and a single radius.
//	@Description	All coordinate pairs are searched in parallel for optimal performance.
//	@ID				bulk-conesearch
//	@Tags			conesearch
//	@Accept			json
//	@Produce		json
//	@Param			request		body		BulkConesearchRequest	true	"Bulk conesearch request with arrays of coordinates"
//	@Success		200			{array}		repository.Mastercat	"Found objects"
//	@Success		204			{string}	string					"No objects found"
//	@Failure		400			{object}	conesearch.ValidationError	"Invalid parameters"
//	@Failure		500			{string}	string					"Internal server error"
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
//	@Summary		Cone search for objects
//	@Description	Search for astronomical objects within a specified radius of given celestial coordinates (RA/Dec).
//	@Description	Returns matching objects from the specified catalog.
//	@ID				conesearch
//	@Tags			conesearch
//	@Accept			json
//	@Produce		json
//	@Param			ra			query		number	true	"Right Ascension (J2000) in degrees"			minimum(0) maximum(360) example(180.5)
//	@Param			dec			query		number	true	"Declination (J2000) in degrees"				minimum(-90) maximum(90) example(-45.0)
//	@Param			radius		query		number	true	"Search radius in degrees"						minimum(0) example(0.01)
//	@Param			catalog		query		string	false	"Catalog to search in"							default(all)
//	@Param			nneighbor	query		integer	false	"Maximum number of neighbors to return"			default(1) minimum(1)
//	@Param			getMetadata	query		boolean	false	"Include full metadata in response"				default(false)
//	@Success		200			{array}		repository.Mastercat	"Found objects"
//	@Success		204			{string}	string					"No objects found"
//	@Failure		400			{object}	conesearch.ValidationError	"Invalid parameters"
//	@Failure		500			{string}	string					"Internal server error"
//	@Router			/conesearch [get]
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

	if getMetadata == "true" {
		result, err := api.conesearchService.FindMetadataByConesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor, catalog)
		if err != nil {
			handleServiceError(err, c)
			return
		}
		handleServiceSuccess(result, c)
	} else {
		result, err := api.conesearchService.Conesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor, catalog)
		if err != nil {
			handleServiceError(err, c)
		}
		handleServiceSuccess(result, c)
	}
}

func handleServiceError(serviceErr error, c *gin.Context) {
	if errors.As(serviceErr, &conesearch.ValidationError{}) {
		c.JSON(http.StatusBadRequest, serviceErr)
	} else {
		c.Error(serviceErr)
		c.JSON(http.StatusInternalServerError, "Could not execute conesearch")
	}
}

func handleServiceSuccess[T any](result []T, c *gin.Context) {
	if len(result) == 0 {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, result)
}
