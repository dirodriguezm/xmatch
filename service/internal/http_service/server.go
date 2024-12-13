package httpservice

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

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

	parsedRa, err := parseRa(ra)
	if err != nil {
		c.Error(err)
	}
	parsedDec, err := parseDec(dec)
	if parsedDec == -999 {
		c.Error(err)
	}
	parsedRadius, err := parseRadius(radius)
	if err != nil {
		c.Error(err)
	}
	catalog, err = parseCatalog(catalog)
	if err != nil {
		c.Error(err)
	}
	parsedNneighbor, err := parseNneighbor(nneighbor)
	if err != nil {
		c.Error(err)
	}

	for _, err := range c.Errors {
		if errors.As(err, &ParseError{}) {
			c.JSON(http.StatusBadRequest, err)
		}
	}
	if len(c.Errors) > 0 {
		return
	}

	result, err := server.conesearchService.Conesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor)
	if err != nil {
		// TODO: diferentiate between SQL errors and other types of error
		slog.Error("Could not execute conesearch", "error", err)
		c.JSON(http.StatusInternalServerError, "Could not execute conesearch")
		return
	}

	c.JSON(http.StatusOK, result)
}
