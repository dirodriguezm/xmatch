package sqlite_writer_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/app"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/testutils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var queries *repository.Queries
var db *sql.DB

func TestMain(m *testing.M) {
	rootPath, err := testutils.FindRootModulePath(5)
	if err != nil {
		slog.Error("could not find root module path", "error", err)
		panic(err)
	}

	dbDir, err := os.MkdirTemp("", "sqlite_writer_test_db_*")
	if err != nil {
		panic(err)
	}
	dbFile := filepath.Join(dbDir, "test.db")

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
    catalog_name: "vlass"
  reader:
    batch_size: 500
    type: "csv"
  database:
    url: "file:%s?_journal_mode=WAL&_sync=NORMAL&_busy_timeout=5000"
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

	cfg, err := app.Config(getenv)
	if err != nil {
		panic(err)
	}
	queries, err = app.Repository(cfg)
	if err != nil {
		panic(err)
	}
	db = queries.GetDbInstance()

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
	os.Remove(dbDir)
	os.Remove(tmpDir)
}

func TestReceive(t *testing.T) {
	ctx := context.Background()
	w := sqlite_writer.New(db, ctx, queries.BulkInsertObject)

	ids := []string{"1", "2"}
	ras := []float64{1, 2}
	decs := []float64{1, 2}
	ipixs := []int64{1, 2}
	cats := []string{"test", "test"}
	w.Write(nil, actor.Message{
		Rows: []any{
			repository.Mastercat{ID: ids[0], Ra: ras[0], Dec: decs[0], Ipix: ipixs[0], Cat: cats[0]},
			repository.Mastercat{ID: ids[1], Ra: ras[1], Dec: decs[1], Ipix: ipixs[1], Cat: cats[1]},
		},
		Error: nil,
	})

	// check the database
	objects, err := queries.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Len(t, objects, 2)
}
