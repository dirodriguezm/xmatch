package test_helpers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func RegisterCatalogsInDB(ctx context.Context, dbFile string) {
	conn := fmt.Sprintf("file:%s", dbFile)
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		slog.Error("Could not create sqlite3 connection", "conn", conn)
		panic(err)
	}
	_, err = db.Exec("select 'test conn'")
	if err != nil {
		slog.Error("Could not connect to database", "conn", conn)
		panic(err)
	}

	repo := repository.New(db)
	if _, err := repo.InsertCatalog(ctx, repository.InsertCatalogParams{Name: "vlass", Nside: 18}); err != nil {
		slog.Error("Could not insert catalog", "conn", conn)
		panic(err)
	}
}
