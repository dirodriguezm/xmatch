package sqlite_writer_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
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
  database:
    url: "file:%s"
  indexer:
    ordering_scheme: "nested"
  indexer_writer:
    type: "sqlite"
`
	config = fmt.Sprintf(config, dbFile)
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
	ch := make(chan indexer.WriterInput[repository.ParquetMastercat])
	var repo conesearch.Repository
	err := ctr.Resolve(&repo)
	require.NoError(t, err)
	ctx := context.Background()
	done := make(chan bool)
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()
	w := sqlite_writer.NewSqliteWriter(repo, ch, done, ctx, src)

	w.Start()
	ids := []string{"1", "2"}
	ras := []float64{1, 2}
	decs := []float64{1, 2}
	ipixs := []int64{1, 2}
	cats := []string{"test", "test"}
	ch <- indexer.WriterInput[repository.ParquetMastercat]{Rows: []repository.ParquetMastercat{
		{ID: &ids[0], Ra: &ras[0], Dec: &decs[0], Ipix: &ipixs[0], Cat: &cats[0]},
		{ID: &ids[1], Ra: &ras[1], Dec: &decs[1], Ipix: &ipixs[1], Cat: &cats[1]},
	}}
	close(ch)
	<-done

	// check the database
	objects, err := repo.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Len(t, objects, 2)
}
