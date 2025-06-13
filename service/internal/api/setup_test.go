// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	api "github.com/dirodriguezm/xmatch/service/internal/api"
	"github.com/dirodriguezm/xmatch/service/internal/di"
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
	slog.Info("Setting up test environment")

	depth := 5
	rootPath, err := utils.FindRootModulePath(depth)
	if err != nil {
		panic(fmt.Errorf("could not find root module path: %w", err))
	}

	// remove test database if exist
	dbFile := filepath.Join(rootPath, "test.db")
	os.Remove(dbFile)

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
  bulk_chunk_size: 1
  max_bulk_concurrency: 1
`
	config = fmt.Sprintf(config, dbFile)
	err = test_helpers.WriteConfigFile(configPath, config)
	if err != nil {
		panic(err)
	}

	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "DB_CONN":
			return fmt.Sprintf("file://%s", dbFile)
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	ctx := context.Background()
	stdout := &strings.Builder{}

	// create tables
	err = test_helpers.Migrate(dbFile, rootPath)
	if err != nil {
		panic(err)
	}

	// register catalogs
	err = test_helpers.RegisterCatalogsInDB(ctx, dbFile)
	if err != nil {
		panic(err)
	}

	ctr = di.BuildServiceContainer(ctx, getenv, stdout)

	// initialize server
	var api *api.API
	err = ctr.Resolve(&api)
	if err != nil {
		panic(fmt.Errorf("could not resolve server: %w", err))
	}
	router = gin.New()
	api.SetupRoutes(router)

	// run tests
	m.Run()

	// cleanup
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(tmpDir)
}
