package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
