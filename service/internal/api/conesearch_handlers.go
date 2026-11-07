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
//	@Summary		Search for objects in a given region using multiple coordinates
//	@Description	Search for objects in a given region using lists of ra, dec and a single radius. Radius, catalog and nneighbor are query parameters; the body only carries coordinates, as a JSON object with ra/dec arrays, a raw CSV body with ra and dec columns, or a multipart/form-data upload in the file field. Radius is in arcseconds. Results are grouped by catalog.
//	@Tags			conesearch
//	@Accept			json,text/csv,multipart/form-data
//	@Produce		json
//
//	@Param			radius		query		string					true	"Radius in arcseconds"
//	@Param			catalog		query		string					false	"Catalog to search in (default: all)"
//	@Param			nneighbor	query		string					false	"Number of neighbors to return (default: 1)"
//	@Param			request		body		BulkConesearchRequest	false	"Bulk conesearch request with ra and dec arrays"
//
//	@Success		200			{array}		conesearch.MastercatResult
//	@Success		204			"No Content"
//	@Failure		400			{object}	conesearch.ValidationError
//	@Failure		500			{string}	string
//	@Router			/bulk-conesearch [post]
func (api *API) conesearchBulk(c *gin.Context) {
	parsedRadius, err := parseRadius(c.Query("radius"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	parsedNneighbor, err := parseNneighbor(c.DefaultQuery("nneighbor", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	catalog := c.DefaultQuery("catalog", "all")

	ra, dec, err := api.parseBulkConesearchInput(c)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := api.conesearchService.BulkConesearch(
		ra,
		dec,
		parsedRadius,
		parsedNneighbor,
		catalog,
		api.config.BulkChunkSize,
		api.config.MaxBulkConcurrency,
	)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	handleServiceSuccess(result, c)
}

func (api *API) parseBulkConesearchInput(c *gin.Context) ([]float64, []float64, error) {
	contentType := c.ContentType()
	switch {
	case isBulkConesearchMultipartContentType(contentType):
		return api.parseBulkConesearchCSVFile(c)
	case isBulkConesearchCSVContentType(contentType):
		return parseBulkConesearchCSV(c.Request.Body, api.config.BulkConesearchMaxRows)
	default:
		var bulkRequest BulkConesearchRequest
		if err := c.ShouldBindJSON(&bulkRequest); err != nil {
			return nil, nil, NewParseError("", "body", err.Error())
		}
		return bulkRequest.Ra, bulkRequest.Dec, nil
	}
}

// Search for objects in a given region
//
//	@Summary		Search for objects in a given region
//	@Description	Search for objects in a given region using ra, dec and radius. Radius is in arcseconds. Returns mastercat results grouped by catalog, or metadata results grouped by catalog when getMetadata is true.
//	@Tags			conesearch
//	@Accept			json
//	@Produce		json
//	@Param			ra			query		string	true	"Right ascension in degrees"
//	@Param			dec			query		string	true	"Declination in degrees"
//	@Param			radius		query		string	true	"Radius in arcseconds"
//	@Param			catalog		query		string	false	"Catalog to search in (default: all)"
//	@Param			nneighbor	query		string	false	"Number of neighbors to return (default: 1)"
//	@Param			getMetadata	query		string	false	"Return metadata results instead of mastercat results (default: false)"
//	@Success		200			{array}		conesearch.MastercatResult
//	@Success		204			"No Content"
//	@Failure		400			{object}	conesearch.ValidationError
//	@Failure		500			{string}	string
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
			return
		}
		handleServiceSuccess(result, c)
	}
}

func handleServiceError(serviceErr error, c *gin.Context) {
	if errors.As(serviceErr, &conesearch.ValidationError{}) {
		c.JSON(http.StatusBadRequest, serviceErr)
	} else {
		_ = c.Error(serviceErr)
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
