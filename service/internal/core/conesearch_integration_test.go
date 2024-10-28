package core_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	dbdb "xmatch/service/internal/db"
	"xmatch/service/internal/di"

	"github.com/golobby/container/v3"
)

func TestMain(m *testing.M) {
	// remove test database, ignore errors
	os.Remove("../../test.db")
	// create test database
	os.Setenv("DB_CONN", "file:../../test.db")

	// build DI container
	di.ContainerBuilder()

	// get db connection from the container
	var db *sql.DB
	container.Resolve(&db)

	// create tables
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, dbdb.SCHEMA); err != nil {
		panic(err)
	}
}
