package httpservice_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
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
	"github.com/stretchr/testify/require"
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

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "pong", w.Body.String())
}

func TestConesearchValidation(t *testing.T) {
	type Expected struct {
		Status int
		Error  map[string]string
	}
	testCases := map[string]Expected{
		"/conesearch": {400, map[string]string{
			"Field":    "RA",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/conesearch?ra=1": {400, map[string]string{
			"Field":    "Dec",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/conesearch?ra=1&dec=1": {400, map[string]string{
			"Field":    "radius",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/conesearch?ra=1&dec=1&radius=1": {200, nil},
		"/conesearch?ra=1&dec=1&radius=1&catalog=a": {400, map[string]string{
			"Field":    "catalog",
			"Reason":   "Catalog not available",
			"ErrValue": "a",
		}},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise": {200, nil},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise&nneighbor=-1": {400, map[string]string{
			"Field":    "nneighbor",
			"ErrValue": "-1",
			"Reason":   "Nneighbor must be a positive integer",
		}},
	}

	for testPath, expected := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testPath, nil)
		router.ServeHTTP(w, req)

		require.Equal(t, expected.Status, w.Code)
		if w.Code == 200 {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		require.Truef(t, maps.EqualFunc(expected.Error, result, func(a string, b interface{}) bool {
			return a == b.(string)
		}), "On %s: values are not equal\n Expected: %v\nReceived: %v", testPath, expected.Error, result)
	}
}
