package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get lightcurve data for coordinates
//
//	@Summary		Get lightcurve data for coordinates
//	@Description	Get lightcurve data for specified coordinates with search radius and neighbor count
//	@Tags			lightcurve
//	@Accept			json
//	@Produce		json
//	@Param			ra			query		string	true	"Right Ascension coordinate"
//	@Param			dec			query		string	true	"Declination coordinate"
//	@Param			radius		query		string	true	"Search radius in arcseconds"
//	@Param			nneighbor	query		string	false	"Number of neighbors to return (default: 1)"
//	@Success		200			{object}	lightcurve.Lightcurve
//	@Failure		400			{string}	string
//	@Failure		500			{string}	string
//	@Router			/lightcurve [get]
func (api *API) Lightcurve(c *gin.Context) {
	ra := c.Query("ra")
	dec := c.Query("dec")
	radius := c.Query("radius")
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
	lightcurve, err := api.lightcurveService.GetLightcurve(parsedRa, parsedDec, parsedRadius, parsedNneighbor)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Could not fetch lightcurve")
	}

	c.JSON(http.StatusOK, lightcurve)
}
