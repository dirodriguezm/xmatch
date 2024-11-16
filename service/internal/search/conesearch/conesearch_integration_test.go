package conesearch_test

import (
	"context"
	"fmt"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"log/slog"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/require"
)

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

func TestConesearch(t *testing.T) {
	var service *conesearch.ConesearchService
	err := container.Resolve(&service)
	if err != nil {
		t.Error(err)
	}

	objects := []repository.InsertObjectParams{
		// ra and dec don't matter here
		{ID: "A", Ipix: 326417514496, Ra: 0, Dec: 0, Cat: "vlass"},
		{ID: "C", Ipix: 327879198247, Ra: 10, Dec: 10, Cat: "vlass"},
	}
	var repo conesearch.Repository
	err = container.Resolve(&repo)
	if err != nil {
		t.Error(err)
	}
	for _, obj := range objects {
		repo.InsertObject(context.Background(), obj)
	}

	result, err := service.Conesearch(0, 0, 1, 10) // TODO: Revisar con Lore algunos casos de prueba
	if err != nil {
		t.Error(err)
	}
	require.Len(t, result, 1)
}
