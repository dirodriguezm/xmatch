package httpservice_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "debug")

	depth := 5
	rootPath, err := utils.FindRootModulePath(depth)
	if err != nil {
		slog.Error("could not find root module path", "depth", depth)
		panic(err)
	}

	// remove test database if exist
	dbFile := filepath.Join(rootPath, "test.db")
	os.Remove(dbFile)

	// set db connection environment variable
	err = os.Setenv("DB_CONN", fmt.Sprintf("file://%s", dbFile))
	if err != nil {
		slog.Error("could not set environment variable")
		panic(err)
	}

	// create a config file
	tmpDir, err := os.MkdirTemp("", "server_test_*")
	if err != nil {
		slog.Error("could not make temp dir")
		panic(err)
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
service:
  database:
    url: "file:%s"
`
	config = fmt.Sprintf(config, dbFile)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		slog.Error("could not write config file")
		panic(err)
	}
	os.Setenv("CONFIG_PATH", configPath)

	// create tables
	mig, err := migrate.New(fmt.Sprintf("file://%s/internal/db/migrations", rootPath), fmt.Sprintf("sqlite3://%s", dbFile))
	if err != nil {
		slog.Error("Could not create Migrate instance")
		panic(err)
	}
	err = mig.Up()
	if err != nil {
		slog.Error("Error during migrations", "error", err)
		panic(err)
	}

	test_helpers.RegisterCatalogsInDB(context.Background(), dbFile)

	ctr := di.BuildServiceContainer()

	// initialize server
	var server *httpservice.HttpServer
	ctr.Resolve(&server)
	router = server.SetupServer()

	// run tests
	m.Run()

	// cleanup
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(tmpDir)
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
