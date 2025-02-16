package httpservice_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/di"
	httpservice "github.com/dirodriguezm/xmatch/service/internal/http_service"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
)

var router *gin.Engine
var ctr container.Container

func beforeTest(t *testing.T) {
	// clear database
	var db *sql.DB
	ctr.Resolve(&db)

	_, err := db.Exec("DELETE FROM mastercat;")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM allwise;")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "debug")
	slog.Info("Setting up test environment")

	depth := 5
	rootPath, err := utils.FindRootModulePath(depth)
	if err != nil {
		panic(fmt.Errorf("could not find root module path: %w", err))
	}

	// remove test database if exist
	dbFile := filepath.Join(rootPath, "test.db")
	os.Remove(dbFile)

	// set db connection environment variable
	err = os.Setenv("DB_CONN", fmt.Sprintf("file://%s", dbFile))
	if err != nil {
		panic(fmt.Errorf("could not set environment variable: %w", err))
	}

	// create a config file
	tmpDir, err := os.MkdirTemp("", "server_test_*")
	if err != nil {
		panic(fmt.Errorf("could not make temp dir: %w", err))
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
service:
  database:
    url: "file:%s"
`
	config = fmt.Sprintf(config, dbFile)
	err = test_helpers.WriteConfigFile(configPath, config)
	if err != nil {
		panic(err)
	}

	// create tables
	err = test_helpers.Migrate(dbFile, rootPath)
	if err != nil {
		panic(err)
	}

	// register catalogs
	err = test_helpers.RegisterCatalogsInDB(context.Background(), dbFile)
	if err != nil {
		panic(err)
	}

	ctr = di.BuildServiceContainer()

	// initialize server
	var server *httpservice.HttpServer
	err = ctr.Resolve(&server)
	if err != nil {
		panic(fmt.Errorf("could not resolve server: %w", err))
	}
	router = server.SetupServer()

	// run tests
	m.Run()

	// cleanup
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(tmpDir)
}
