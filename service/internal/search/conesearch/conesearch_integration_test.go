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

package conesearch_test

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
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/dirodriguezm/xmatch/service/internal/testutils"

	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/allwise"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/erosita"
	_ "github.com/dirodriguezm/xmatch/service/internal/catalog/gaia"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var configPath string
var dbFile string

func TestMain(m *testing.M) {
	rootPath, err := testutils.FindRootModulePath(5)
	if err != nil {
		panic(err)
	}

	// remove test database, ignore errors
	dbDir, err := os.MkdirTemp("", "conesearch_test_db_*")
	if err != nil {
		panic(err)
	}
	dbFile = filepath.Join(dbDir, "test.db")

	// create temporary directory for config
	tmpDir, err := os.MkdirTemp("", "xmatch-test-*")
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// create config file
	configPath = filepath.Join(tmpDir, "config.yaml")
	config := fmt.Sprintf(`
service:
  database:
    url: "file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
  host: "localhost:8080"
  base_path: "/v1"
  bulk_chunk_size: 500
  max_bulk_concurrency: 1
  lightcurve_service:
    neowise:
      use_cntr_filter: true
`, dbFile)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		panic(fmt.Errorf("could not write config file: %w", err))
	}

	// create tables
	mig, err := migrate.New(fmt.Sprintf("file://%s/internal/db/migrations", rootPath), fmt.Sprintf("sqlite3://%s", dbFile))
	if err != nil {
		panic(fmt.Errorf("Could not create Migrate instance: %w", err))
	}
	err = mig.Up()
	if err != nil {
		panic(fmt.Errorf("Error during migrations: %w", err))
	}

	err = test_helpers.RegisterCatalogsInDB(context.Background(), dbFile)
	if err != nil {
		panic(fmt.Errorf("registering catalogs: %w", err))
	}

	// run tests
	code := m.Run()

	// cleanup
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(dbDir)
	os.Remove(tmpDir)

	os.Exit(code)
}

func TestConesearch(t *testing.T) {
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

	queries := app.ServiceRepository(db)
	defer CleanDB(t, queries)

	resolver := catalog.NewResolver(queries)
	service, err := app.ConesearchService(queries, resolver)
	if err != nil {
		t.Fatalf("creating conesearch service: %v", err)
	}

	objects := []repository.Mastercat{
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "B", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	for _, obj := range objects {
		err = queries.InsertMastercat(context.Background(), obj)
		if err != nil {
			t.Fatalf("inserting object: %v", err)
		}
	}

	result, err := service.Conesearch(0, 0, 1, 10, "all")
	if err != nil {
		t.Error(err)
	}
	require.Len(t, result, 1, "conesearch should get one object but got %d", len(result))
}

func TestConesearch_WithMetadata(t *testing.T) {
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

	queries := app.ServiceRepository(db)
	defer CleanDB(t, queries)

	resolver := catalog.NewResolver(queries)
	service, err := app.ConesearchService(queries, resolver)
	if err != nil {
		t.Fatalf("creating conesearch service: %v", err)
	}

	objects := []repository.Mastercat{
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "B", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	for _, obj := range objects {
		ctx := context.Background()
		err = queries.InsertMastercat(ctx, obj)
		if err != nil {
			t.Fatalf("inserting mastercat: %v", err)
		}
		err = queries.InsertAllwise(ctx, repository.InsertAllwiseParams{ID: obj.ID})
		if err != nil {
			t.Fatalf("inserting allwise: %v", err)
		}
	}

	result, err := service.FindMetadataByConesearch(0, 0, 1, 10, "allwise")
	if err != nil {
		t.Error(err)
	}
	require.Len(t, result, 1, "conesearch should get one object but got %d", len(result))
}

func TestBulkConesearch(t *testing.T) {
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

	queries := app.ServiceRepository(db)
	defer CleanDB(t, queries)

	resolver := catalog.NewResolver(queries)
	service, err := app.ConesearchService(queries, resolver)
	if err != nil {
		t.Fatalf("creating conesearch service: %v", err)
	}

	// insert objects
	objects := []repository.Mastercat{
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "B", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	for _, obj := range objects {
		err = queries.InsertMastercat(context.Background(), obj)
		if err != nil {
			t.Fatalf("inserting object: %v", err)
		}
	}

	// set up test cases
	type testCase struct {
		ra        []float64
		dec       []float64
		radius    float64
		nneighbor int
		expected  []string
	}
	testCases := []testCase{
		{ra: []float64{0}, dec: []float64{0}, radius: 1, nneighbor: 10, expected: []string{"A"}},
		{ra: []float64{0, 1, 2}, dec: []float64{0, 1, 2}, radius: 1, nneighbor: 10, expected: []string{"A"}},
		{ra: []float64{10}, dec: []float64{10}, radius: 1, nneighbor: 10, expected: []string{"B"}},
		{ra: []float64{0, 10}, dec: []float64{0, 10}, radius: 1, nneighbor: 10, expected: []string{"A", "B"}},
	}

	// test bulk conesearch
	for _, tc := range testCases {
		result, err := service.BulkConesearch(tc.ra, tc.dec, tc.radius, tc.nneighbor, "all", 1, 1)
		if err != nil {
			t.Error(err)
		}

		foundIDs := make(map[string]bool)
		for i := range result {
			for j := range result[i].Data {
				id := result[i].Data[j].ID
				require.Contains(t, tc.expected, id, "testCase: %v | result: %v", tc, result)
				foundIDs[id] = true
			}
		}
		for _, expectedID := range tc.expected {
			require.True(t, foundIDs[expectedID], "expected ID %s not found in result for testCase: %v", expectedID, tc)
		}
	}

	// test that coordinates with no matches are properly handled
	t.Run("coordinates with no matches return empty results", func(t *testing.T) {
		noMatchResult, err := service.BulkConesearch(
			[]float64{100, 200}, // coordinates with no objects nearby
			[]float64{50, 60},
			1,
			10,
			"all",
			1,
			1,
		)
		require.NoError(t, err)
		require.Len(t, noMatchResult, 0, "expected no matches for coordinates with no nearby objects, got: %v", noMatchResult)
	})

	// test that Index field is correctly set
	t.Run("index field is correctly set for matched coordinates", func(t *testing.T) {
		multiMatchResult, err := service.BulkConesearch(
			[]float64{0, 10}, // first matches A, second matches B
			[]float64{0, 10},
			1,
			10,
			"all",
			1,
			1,
		)
		require.NoError(t, err)
		require.Len(t, multiMatchResult, 2, "expected 2 results, got: %v", multiMatchResult)

		indexMap := make(map[int][]string)
		for _, r := range multiMatchResult {
			for _, d := range r.Data {
				indexMap[r.Index] = append(indexMap[r.Index], d.ID)
			}
		}
		require.Contains(t, indexMap[0], "A", "index 0 should contain ID A")
		require.Contains(t, indexMap[1], "B", "index 1 should contain ID B")
	})

	// test that non-matching coordinates in the middle are handled correctly
	t.Run("non-matching coordinates in the middle preserve index correctness", func(t *testing.T) {
		middleNoMatchResult, err := service.BulkConesearch(
			[]float64{0, 50, 10, 60},
			[]float64{0, 50, 10, 60},
			1,
			10,
			"all",
			1,
			1,
		)
		require.NoError(t, err)
		require.Len(t, middleNoMatchResult, 2, "expected 2 results, got: %v", middleNoMatchResult)

		indexMap := make(map[int][]string)
		for _, r := range middleNoMatchResult {
			for _, d := range r.Data {
				indexMap[r.Index] = append(indexMap[r.Index], d.ID)
			}
		}
		require.Contains(t, indexMap[0], "A", "index 0 should contain ID A")
		require.Contains(t, indexMap[2], "B", "index 2 should contain ID B")
		require.NotContains(t, indexMap, 1, "index 1 should have no matches")
		require.NotContains(t, indexMap, 3, "index 3 should have no matches")
	})

}

func CleanDB(t *testing.T, repo interface {
	RemoveAllObjects(context.Context) error
}) {
	err := repo.RemoveAllObjects(context.Background())
	require.NoError(t, err)
}
