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

package sqlite_writer_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var ctr container.Container

func TestMain(m *testing.M) {
	depth := 5
	rootPath, err := utils.FindRootModulePath(depth)
	if err != nil {
		slog.Error("could not find root module path", "depth", depth)
		panic(err)
	}

	// remove test database if exist
	dbFile := filepath.Join(rootPath, "test.db")
	os.Remove(dbFile)

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

	// build DI container
	ctr = di.BuildIndexerContainer(ctx, getenv, stdout)

	// create tables
	mig, err := migrate.New(fmt.Sprintf("file://%s/internal/db/migrations", rootPath), fmt.Sprintf("sqlite://%s", dbFile))
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
	ch := make(chan writer.WriterInput[repository.Mastercat])
	var repo conesearch.Repository
	err := ctr.Resolve(&repo)
	require.NoError(t, err)
	ctx := context.Background()
	done := make(chan struct{})
	src := source.ASource(t).WithUrl(fmt.Sprintf("files:%s", t.TempDir())).Build()
	w := sqlite_writer.NewSqliteWriter(repo, ch, done, ctx, src)

	w.Start()
	ids := []string{"1", "2"}
	ras := []float64{1, 2}
	decs := []float64{1, 2}
	ipixs := []int64{1, 2}
	cats := []string{"test", "test"}
	ch <- writer.WriterInput[repository.Mastercat]{Rows: []repository.Mastercat{
		{ID: ids[0], Ra: ras[0], Dec: decs[0], Ipix: ipixs[0], Cat: cats[0]},
		{ID: ids[1], Ra: ras[1], Dec: decs[1], Ipix: ipixs[1], Cat: cats[1]},
	}}
	close(ch)
	<-done

	// check the database
	objects, err := repo.GetAllObjects(ctx)
	require.NoError(t, err)
	require.Len(t, objects, 2)
}
