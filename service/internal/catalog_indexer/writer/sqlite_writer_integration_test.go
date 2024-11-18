package writer_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/require"
)

var ctr container.Container

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "debug")

	depth := 5
	rootPath, err := findRootModulePath(depth)
	if err != nil {
		slog.Error("could not fin root module path", "depth", depth)
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
	tmpDir, err := os.MkdirTemp("", "sqlite_writer_integration_test_*")
	if err != nil {
		slog.Error("could not make temp dir")
		panic(err)
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
catalog_indexer:
  source:
    url: "buffer:"
    type: "csv"
  reader:
    batch_size: 500
    type: "csv"
`
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		slog.Error("could not write config file")
		panic(err)
	}
	os.Setenv("CONFIG_PATH", configPath)

	// build DI container
	ctr = di.BuildIndexerContainer()

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
	m.Run()
	os.Remove(configPath)
	os.Remove(dbFile)
	os.Remove(tmpDir)
}

func TestActor(t *testing.T) {
	ch := make(chan indexer.IndexerResult)
	var repo conesearch.Repository
	err := ctr.Resolve(&repo)
	require.NoError(t, err)
	ctx := context.Background()
	w := writer.NewSqliteWriter(repo, ch, ctx)

	w.Start()
	ch <- indexer.IndexerResult{Objects: []repository.Mastercat{
		{ID: "oid1", Ipix: 1, Ra: 1, Dec: 1, Cat: "vlass"},
		{ID: "oid2", Ipix: 2, Ra: 2, Dec: 2, Cat: "vlass"},
	}}
	close(ch)
	<-w.Done

	// check the database
	objects, err := repo.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Len(t, objects, 2)
}

func findRootModulePath(maxDepth int) (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dirs, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		if dir.Name() == "go.mod" {
			return currDir, nil
		}
	}

	os.Chdir("..")
	return findRootModulePath(maxDepth - 1)
}
