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
	"database/sql"
	"errors"
	"net/http"

	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
	"github.com/gin-gonic/gin"
)

// Find metadata by id
//
//	@Summary		Get metadata by object ID
//	@Description	Retrieve detailed catalog metadata for a specific astronomical object by its identifier.
//	@ID				metadata
//	@Tags			metadata
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"Object identifier"		example(J120000.00-450000.0)
//	@Param			catalog	query		string	true	"Catalog to search in"	example(allwise)
//	@Success		200		{object}	repository.Allwise		"Object metadata"
//	@Success		204		{string}	string					"Object not found"
//	@Failure		400		{object}	metadata.ValidationError	"Invalid parameters"
//	@Failure		500		{string}	string					"Internal server error"
//	@Router			/metadata [get]
func (api *API) metadata(c *gin.Context) {
	id := c.Query("id")
	catalog := c.Query("catalog")

	result, err := api.metadataService.FindByID(c.Request.Context(), id, catalog)
	if err != nil {
		if errors.As(err, &metadata.ValidationError{}) {
			c.JSON(http.StatusBadRequest, err)
			// WARN: sql reference should be handled inside service, not in this layer
		} else if errors.Is(err, sql.ErrNoRows) {
			c.Writer.WriteHeader(http.StatusNoContent)
		} else if errors.As(err, &metadata.ArgumentError{}) {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Could not execute metadata query")
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// Find metadata by multiple ids
//
//	@Summary		Bulk get metadata by object IDs
//	@Description	Retrieve detailed catalog metadata for multiple astronomical objects by their identifiers.
//	@Description	All IDs are queried in parallel for optimal performance.
//	@ID				bulk-metadata
//	@Tags			metadata
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BulkMetadataRequest	true	"Bulk metadata request with list of object IDs"
//	@Success		200		{array}		repository.Allwise		"Objects metadata"
//	@Success		204		{string}	string					"No objects found"
//	@Failure		400		{object}	metadata.ValidationError	"Invalid parameters"
//	@Failure		500		{string}	string					"Internal server error"
//	@Router			/bulk-metadata [post]
func (api *API) metadataBulk(c *gin.Context) {
	var bulkRequest BulkMetadataRequest
	if err := c.ShouldBindJSON(&bulkRequest); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := api.metadataService.BulkFindByID(c.Request.Context(), bulkRequest.Ids, bulkRequest.Catalog)
	if err != nil {
		if errors.As(err, &metadata.ValidationError{}) {
			c.JSON(http.StatusBadRequest, err)
			// WARN: sql reference should be handled inside service, not in this layer
		} else if errors.Is(err, sql.ErrNoRows) {
			c.Writer.WriteHeader(http.StatusNoContent)
		} else if errors.As(err, &metadata.ArgumentError{}) {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Could not execute metadata query")
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
