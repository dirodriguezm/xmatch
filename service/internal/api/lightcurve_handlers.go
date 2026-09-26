package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get lightcurve data for coordinates
//
//	@Summary		Get lightcurve data
//	@Description	Retrieve time-series photometry data (lightcurve) for astronomical objects at the specified coordinates.
//	@Description	Returns detections, non-detections, and forced photometry measurements.
//	@ID				lightcurve
//	@Tags			lightcurve
//	@Accept			json
//	@Produce		json
//	@Param			ra			query		number	true	"Right Ascension (J2000) in degrees"		minimum(0) maximum(360) example(180.5)
//	@Param			dec			query		number	true	"Declination (J2000) in degrees"			minimum(-90) maximum(90) example(-45.0)
//	@Param			radius		query		number	true	"Search radius in arcseconds"				minimum(0) example(5.0)
//	@Param			catalog		query		string	false	"Catalog to query (all, ztf, neowise, allwise)"	default(all)
//	@Param			nneighbor	query		integer	false	"Number of neighbors to return"				default(1) minimum(1)
//	@Success		200			{object}	LightcurveResponse	"Lightcurve data"
//	@Failure		400			{string}	string				"Invalid parameters"
//	@Failure		500			{string}	string				"Internal server error"
//	@Router			/lightcurve [get]
func (api *API) Lightcurve(c *gin.Context) {
	ra := c.Query("ra")
	dec := c.Query("dec")
	radius := c.Query("radius")
	catalog := c.DefaultQuery("catalog", "all")
	nneighbor := c.DefaultQuery("nneighbor", "1")

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
	parsedCatalog, err := parseLightcurveCatalog(catalog)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	lightcurve, err := api.lightcurveService.GetLightcurve(parsedRa, parsedDec, parsedRadius, parsedNneighbor, parsedCatalog)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Could not fetch lightcurve")
		return
	}

	response, err := newLightcurveResponse(lightcurve)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Could not serialize lightcurve")
		return
	}

	c.JSON(http.StatusOK, response)
}
