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

	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/dirodriguezm/xmatch/service/internal/testutils"
	"github.com/gin-gonic/gin"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"
)

var router *gin.Engine
var configPath string

func beforeTest(t *testing.T) {
	// clear database
	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	stdout := &strings.Builder{}
	cfg, err := app.Config(getenv)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		t.Fatalf("creating database connection: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Try to delete from tables, but don't fail if they don't exist yet
	// The tables are created in TestMain via migrations
	_, _ = db.Exec("DELETE FROM mastercat;")
	_, _ = db.Exec("DELETE FROM allwise;")
}

func TestMain(m *testing.M) {
	slog.Info("Setting up test environment")

	rootPath, err := testutils.FindRootModulePath(5)
	if err != nil {
		panic(fmt.Errorf("could not find root module path: %w", err))
	}

	// remove test database if exist
	dbDir, err := os.MkdirTemp("", "api_test_db_*")
	if err != nil {
		panic(fmt.Errorf("could not make db temp dir: %w", err))
	}
	dbFile := filepath.Join(dbDir, "test.db")

	// create a config file
	tmpDir, err := os.MkdirTemp("", "server_test_*")
	if err != nil {
		panic(fmt.Errorf("could not make temp dir: %w", err))
	}
	configPath = filepath.Join(tmpDir, "config.yaml")
	config := `
service:
  database:
    url: "file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
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

	cfg, err := app.Config(getenv)
	if err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		panic(fmt.Errorf("creating database connection: %w", err))
	}

	queries := app.ServiceRepository(db)

	resolver := catalog.NewResolver(queries)

	conesearchService, err := app.ConesearchService(queries, resolver)
	if err != nil {
		_ = db.Close()
		panic(fmt.Errorf("creating conesearch service: %w", err))
	}

	metadataService, err := app.MetadataService(resolver)
	if err != nil {
		_ = db.Close()
		panic(fmt.Errorf("creating metadata service: %w", err))
	}

	lightcurveService, err := app.LightcurveService(cfg, conesearchService)
	if err != nil {
		_ = db.Close()
		panic(fmt.Errorf("creating lightcurve service: %w", err))
	}

	api, err := app.API(conesearchService, metadataService, lightcurveService, cfg.Service, getenv)
	if err != nil {
		_ = db.Close()
		panic(fmt.Errorf("creating API: %w", err))
	}

	router = gin.New()
	api.SetupRoutes(router)

	// run tests
	code := m.Run()

	// cleanup
	_ = db.Close()
	_ = os.Remove(configPath)
	_ = os.Remove(dbFile)
	_ = os.Remove(dbDir)
	_ = os.Remove(tmpDir)

	os.Exit(code)
}
