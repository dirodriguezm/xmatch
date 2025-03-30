package conesearch_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"github.com/golobby/container/v3"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

var ctr container.Container

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "debug")

	rootPath, err := utils.FindRootModulePath(5)
	if err != nil {
		panic(err)
	}

	// remove test database, ignore errors
	dbFile := filepath.Join(rootPath, "test.db")
	os.Remove(dbFile)

	// // set db connection environment variable
	// err = os.Setenv("DB_CONN", fmt.Sprintf("file://%s", dbFile))
	// if err != nil {
	// 	slog.Error("could not set environment variable")
	// 	panic(err)
	// }

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

	// build DI container
	ctr = di.BuildServiceContainer()

	// run tests
	m.Run()

	// cleanup
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(tmpDir)
}

func TestConesearch(t *testing.T) {
	var service *conesearch.ConesearchService
	err := ctr.Resolve(&service)
	if err != nil {
		t.Error(err)
	}

	objects := []repository.InsertObjectParams{
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "B", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	var repo conesearch.Repository
	err = ctr.Resolve(&repo)
	if err != nil {
		t.Error(err)
	}
	for _, obj := range objects {
		repo.InsertObject(context.Background(), obj)
	}

	result, err := service.Conesearch(0, 0, 1, 10, "all")
	if err != nil {
		t.Error(err)
	}
	require.Len(t, result, 1, "conesearch should get one object but got %d", len(result))

	CleanDB(t, repo)
}

func TestBulkConesearch(t *testing.T) {
	// initialize service
	var service *conesearch.ConesearchService
	err := ctr.Resolve(&service)
	if err != nil {
		t.Error(err)
	}

	// insert objects
	objects := []repository.InsertObjectParams{
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "B", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	var repo conesearch.Repository
	err = ctr.Resolve(&repo)
	if err != nil {
		t.Error(err)
	}
	for _, obj := range objects {
		repo.InsertObject(context.Background(), obj)
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

		require.Len(t, tc.expected, len(result), "testCase: %v | result: %v", tc, result)
		for i := range result {
			id := result[i].ID
			require.Contains(t, tc.expected, id, "testCase: %v | result: %v", tc, result)
		}
	}

	CleanDB(t, repo)
}

func CleanDB(t *testing.T, repo conesearch.Repository) {
	err := repo.RemoveAllObjects(context.Background())
	require.NoError(t, err)
}
