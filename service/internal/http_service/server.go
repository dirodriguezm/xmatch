package httpservice

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	conesearchService *conesearch.ConesearchService
}

func NewHttpServer(conesearchService *conesearch.ConesearchService) (*HttpServer, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	return &HttpServer{conesearchService: conesearchService}, nil
}

func (server *HttpServer) SetupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/conesearch", server.conesearchHandler)
	r.SetTrustedProxies([]string{"localhost"})
	return r
}

func (server *HttpServer) InitServer() {
	r := server.SetupServer()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func (server *HttpServer) conesearchHandler(c *gin.Context) {
	ra := c.Query("ra")
	dec := c.Query("dec")
	radius := c.Query("radius")
	catalog := c.DefaultQuery("catalog", "all")
	nneighbor := c.DefaultQuery("nneighbor", "1")
	parsedRa, errMsg := validateRa(ra)
	// TODO: move validation to service
	if parsedRa == -999 {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	parsedDec, errMsg := validateDec(dec)
	if parsedDec == -999 {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	parsedRadius, errMsg := validateRadius(radius)
	if parsedRadius == -999 {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	catalog, errMsg = validateCatalog(catalog)
	if catalog == "" {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	parsedNneighbor, errMsg := validateNneighbor(nneighbor)
	if parsedNneighbor == -999 {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	result, err := server.conesearchService.Conesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor)
	if err != nil {
		slog.Error("Could not execute conesearch", "error", err)
		c.String(http.StatusInternalServerError, "Could not execute conesearch")
		return
	}
	c.JSON(http.StatusOK, result)
}

func validateRadius(rad string) (float64, string) {
	radius, err := strconv.ParseFloat(rad, 64)
	if err != nil {
		msg := "Could not parse radius `%s`\n"
		if rad == "" {
			msg = "Radius can't be empty%s\n"
		}
		return -999, fmt.Sprintf(msg, rad)
	}
	if err := conesearch.ValidateRadius(radius); err != nil {
		return -999, fmt.Sprintf("Invalid radius: %s", err.Error())
	}
	return radius, ""
}

func validateRa(ra string) (float64, string) {
	parsedRa, err := strconv.ParseFloat(ra, 64)
	if err != nil {
		msg := "Could not parse RA `%s`\n"
		if ra == "" {
			msg = "RA can't be empty%s\n"
		}
		return -999, fmt.Sprintf(msg, ra)
	}
	if err := conesearch.ValidateRa(parsedRa); err != nil {
		return -999, fmt.Sprintf("Invalid ra: %s", err.Error())
	}
	return parsedRa, ""
}

func validateDec(dec string) (float64, string) {
	parsedDec, err := strconv.ParseFloat(dec, 64)
	if err != nil {
		msg := "Could not parse Dec `%s`\n"
		if dec == "" {
			msg = "Dec can't be empty%s\n"
		}
		return -999, fmt.Sprintf(msg, dec)
	}
	if err := conesearch.ValidateDec(parsedDec); err != nil {
		return -999, fmt.Sprintf("Invalid dec: %s", err.Error())
	}
	return parsedDec, ""
}

func validateCatalog(catalog string) (string, string) {
	allowed := []string{"all", "wise", "vlass", "lsdr10"}
	isAllowed := false
	catalog = strings.ToLower(catalog)
	for _, cat := range allowed {
		if catalog == cat {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return "", fmt.Sprintf("Catalog must be one of %s\n", allowed)
	}
	return catalog, ""
}

func validateNneighbor(nneighbor string) (int, string) {
	parsedNneighbor, err := strconv.Atoi(nneighbor)
	if err != nil {
		if nneighbor == "" {
			return -999, "Nneighbor can't be empty\n"
		}
		return -999, fmt.Sprintf("Could not parse nneighbor %v\n", nneighbor)
	}
	if err := conesearch.ValidateNneighbor(parsedNneighbor); err != nil {
		return -999, fmt.Sprintf("Invalid nneighbor: %s", err.Error())
	}
	return parsedNneighbor, ""
}
