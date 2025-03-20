package httpservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/metadata"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/dirodriguezm/xmatch/service/docs"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	conesearchService *conesearch.ConesearchService
	metadataService   *metadata.MetadataService
}

func NewHttpServer(
	conesearchService *conesearch.ConesearchService,
	metadataService *metadata.MetadataService,
) (*HttpServer, error) {
	if conesearchService == nil {
		return nil, fmt.Errorf("ConesearchService was nil while creating HttpServer")
	}
	if metadataService == nil {
		return nil, fmt.Errorf("MetadataService was nil while creating HttpServer")
	}
	return &HttpServer{conesearchService: conesearchService, metadataService: metadataService}, nil
}

func (server *HttpServer) SetupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/conesearch", server.conesearchHandler)
		v1.POST("/bulk-conesearch", server.conesearchBulkHandler)
		v1.GET("/metadata", server.metadataHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.SetTrustedProxies([]string{"localhost"})
	return r
}

func (server *HttpServer) InitServer() {
	r := server.SetupServer()
	r.Run() // listen and serve on 0.0.0.0:8080
}

// Search for objects in a given region
//
//	@Summary		Search for objects in a given region
//	@Description	Search for objects in a given region using ra, dec and radius
//	@Tags			conesearch
//	@Accept			json
//	@Produce		json
//	@Param			ra			query		string	true	"Right ascension in degrees"
//	@Param			dec			query		string	true	"Declination in degrees"
//	@Param			radius		query		string	true	"Radius in degrees"
//	@Param			catalog		query		string	false	"Catalog to search in"
//	@Param			nneighbor	query		string	false	"Number of neighbors to return"
//	@Success		200			{array}		repository.Mastercat
//	@Success		204			{string}	string
//	@Failure		400			{object}	conesearch.ValidationError
//	@Failure		500			{string}	string
//	@Router			/conesearch [get]
func (server *HttpServer) conesearchHandler(c *gin.Context) {
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

	result, err := server.conesearchService.Conesearch(parsedRa, parsedDec, parsedRadius, parsedNneighbor, catalog)
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

type BulkRequest struct {
	Ra        []float64 `json:"ra"`
	Dec       []float64 `json:"dec"`
	Radius    float64   `json:"radius"`
	Catalog   string    `json:"catalog"`
	Nneighbor int       `json:"nneighbor"`
}

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
func (server *HttpServer) conesearchBulkHandler(c *gin.Context) {
	var bulkRequest BulkRequest
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

	result, err := server.conesearchService.BulkConesearch(bulkRequest.Ra, bulkRequest.Dec, bulkRequest.Radius, bulkRequest.Nneighbor, bulkRequest.Catalog)
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

// Find metadata by id
//
//	@Summary		Search for metadata by id
//	@Description	Search for metadata by id
//	@Tags			metadata
//	@Accept			json
//	@Produce		json
//	@Param			id		query		string	true	"ID to search for"
//	@Param			catalog	query		string	true	"Catalog to search in"
//	@Success		200		{object}	repository.AllwiseMetadata
//	@Success		204		{string}	string
//	@Failure		400		{object}	metadata.ValidationError
//	@Failure		500		{string}	string
//	@Router			/metadata [get]
func (server *HttpServer) metadataHandler(c *gin.Context) {
	id := c.Query("id")
	catalog := c.Query("catalog")

	result, err := server.metadataService.FindByID(c.Request.Context(), id, catalog)
	if err != nil {
		if errors.As(err, &metadata.ValidationError{}) {
			c.JSON(http.StatusBadRequest, err)
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
