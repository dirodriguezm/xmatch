package httpservice

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"xmatch/service/internal/core"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	conesearchService *core.ConesearchService
}

func NewHttpServer(conesearchService *core.ConesearchService) HttpServer {
	return HttpServer{conesearchService: conesearchService}
}

func (server HttpServer) SetupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/conesearch", conesearch)
	r.SetTrustedProxies([]string{"localhost"})
	return r
}

func (server HttpServer) InitServer() {
	r := server.SetupServer()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func conesearch(c *gin.Context) {
	ra := c.Query("ra")
	dec := c.Query("dec")
	radius := c.Query("radius")
	catalog := c.DefaultQuery("catalog", "all")
	nneighbor := c.DefaultQuery("nneighbor", "1")
	parsedRa, errMsg := validateRa(ra)
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
	}
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
	if radius <= 0 {
		return -999, "Radius can't be lower than 0\n"
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
	if parsedRa < 0 {
		return -999, "RA can't be lower than 0\n"
	}
	if parsedRa > 360 {
		return -999, "RA can't be greater than 360\n"
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
	if parsedDec < -90 {
		return -999, "Dec can't be lower than -90\n"
	}
	if parsedDec > 90 {
		return -999, "Dec can't be greater than 90\n"
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
	if parsedNneighbor <= 0 {
		return -999, fmt.Sprintf("Nneighbor must be a positive integer\n")
	}
	return parsedNneighbor, ""
}
