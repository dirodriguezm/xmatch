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

package test_helpers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RegisterCatalogsInDB(ctx context.Context, dbFile string) error {
	conn := fmt.Sprintf("file:%s", dbFile)
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return fmt.Errorf("could not create sqlite3 connection: %w", err)
	}
	_, err = db.Exec("select 'test conn'")
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	repo := repository.New(db)
	if _, err := repo.InsertCatalog(ctx, repository.InsertCatalogParams{Name: "vlass", Nside: 18}); err != nil {
		return fmt.Errorf("could not insert catalog vlass: %w", err)
	}
	if _, err := repo.InsertCatalog(ctx, repository.InsertCatalogParams{Name: "allwise", Nside: 18}); err != nil {
		return fmt.Errorf("could not insert catalog allwise: %w", err)
	}
	return nil
}

func WriteConfigFile(configPath, config string) error {
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		slog.Error("could not write config file")
		return err
	}
	return os.Setenv("CONFIG_PATH", configPath)
}

func Migrate(dbFile string, rootPath string) error {
	mig, err := migrate.New(fmt.Sprintf("file://%s/internal/db/migrations", rootPath), fmt.Sprintf("sqlite3://%s", dbFile))
	if err != nil {
		slog.Error("Could not create Migrate instance")
		return err
	}
	return mig.Up()
}

func InsertAllwiseMastercat(nobjects int, db *sql.DB) error {
	repo := repository.New(db)
	for i := 0; i < nobjects; i++ {
		ra := i
		dec := i
		// ra can't be greater than 360
		if ra > 360 {
			ra = ra % 360
		}
		// dec can't be greater than 90
		if dec > 90 {
			dec = dec % 90
		}
		// dec can't be less than -90
		if dec < -90 {
			dec = dec % 90
		}

		point := healpix.RADec(float64(ra), float64(dec))
		mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
		if err != nil {
			return fmt.Errorf("could not create healpix mapper: %w", err)
		}
		ipix := mapper.PixelAt(point)

		// insert object
		_, err = repo.InsertObject(context.Background(), repository.InsertObjectParams{
			ID:   fmt.Sprintf("allwise-%d", i),
			Ipix: ipix,
			Ra:   float64(ra),
			Dec:  float64(dec),
			Cat:  "allwise",
		})
		if err != nil {
			return fmt.Errorf("could not insert allwise mastercat: %w", err)
		}
	}

	return nil
}

func InsertAllwiseMetadata(nobjects int, db *sql.DB) error {
	repo := repository.New(db)
	for i := 0; i < nobjects; i++ {
		metadata := repository.InsertAllwiseParams{
			ID:         fmt.Sprintf("allwise-%d", i),
			W1mpro:     sql.NullFloat64{Float64: 1.0, Valid: true},
			W1sigmpro:  sql.NullFloat64{Float64: 0.1, Valid: true},
			W2mpro:     sql.NullFloat64{Float64: 2.0, Valid: true},
			W2sigmpro:  sql.NullFloat64{Float64: 0.2, Valid: true},
			W3mpro:     sql.NullFloat64{Float64: 3.0, Valid: true},
			W3sigmpro:  sql.NullFloat64{Float64: 0.3, Valid: true},
			W4mpro:     sql.NullFloat64{Float64: 4.0, Valid: true},
			W4sigmpro:  sql.NullFloat64{Float64: 0.4, Valid: true},
			JM2mass:    sql.NullFloat64{Float64: 5.0, Valid: true},
			JMsig2mass: sql.NullFloat64{Float64: 0.5, Valid: true},
			HM2mass:    sql.NullFloat64{Float64: 6.0, Valid: true},
			HMsig2mass: sql.NullFloat64{Float64: 0.6, Valid: true},
			KM2mass:    sql.NullFloat64{Float64: 7.0, Valid: true},
			KMsig2mass: sql.NullFloat64{Float64: 0.7, Valid: true},
		}
		err := repo.InsertAllwise(context.Background(), metadata)
		if err != nil {
			return fmt.Errorf("could not insert allwise metadata: %w", err)
		}
	}
	return nil
}
