package httpservice_test

import (
	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	ctr := di.BuildServiceContainer()
	var server httpservice.HttpServer
	ctr.Resolve(&server)
	router = server.SetupServer()
	m.Run()
}

func TestPingRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestConesearchValidation(t *testing.T) {
	type Expected struct {
		Status       int
		ErrorMessage string
	}
	testCases := map[string]Expected{
		"/conesearch":                                               {400, "RA can't be empty\n"},
		"/conesearch?ra=1":                                          {400, "Dec can't be empty\n"},
		"/conesearch?ra=1&dec=1":                                    {400, "Radius can't be empty\n"},
		"/conesearch?ra=1&dec=1&radius=1":                           {200, ""},
		"/conesearch?ra=1&dec=1&radius=1&catalog=a":                 {400, "Catalog must be one of [all wise vlass lsdr10]\n"},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise":              {200, ""},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise&nneighbor=-1": {400, "Nneighbor must be a positive integer\n"},
	}

	for testPath, expected := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testPath, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, expected.Status, w.Code)
		assert.Contains(t, w.Body.String(), expected.ErrorMessage)
	}
}
