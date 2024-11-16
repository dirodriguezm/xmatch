package writer_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

func TestMain(m *testing.M) {
	os.Setenv("LOG_LEVEL", "debug")

	rootPath, err := findRootModulePath(5)
	if err != nil {
		panic(err)
	}

	// remove test database, ignore errors
	dbFile := fmt.Sprintf("%s/test.db", rootPath)
	os.Remove(dbFile)

	// create test database
	err = os.Setenv("DB_CONN", fmt.Sprintf("file://%s", dbFile))
	if err != nil {
		panic(err)
	}

	// build DI container
	di.ContainerBuilder()

	// create tables
	mig, err := migrate.New(fmt.Sprintf("file://%s/internal/db/migrations", rootPath), fmt.Sprintf("sqlite3://%s", dbFile))
	if err != nil {
		panic(err)
	}
	err = mig.Up()
	if err != nil {
		slog.Error("Error during migrations", "error", err)
	}
	m.Run()
}

func TestActor(t *testing.T) {
	ch := make(chan indexer.IndexerResult)
	var repo conesearch.Repository
	err := container.Resolve(&repo)
	require.NoError(t, err)
	ctx := context.Background()
	w := writer.New(repo, ch, ctx)

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
