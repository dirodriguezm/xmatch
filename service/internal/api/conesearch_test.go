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

package api_test

import (
	"bytes"
	"context"
	"database/sql"

	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"

	"github.com/stretchr/testify/require"
)

func TestConesearch_Validation(t *testing.T) {
	type Expected struct {
		Status int
		Error  map[string]string
	}
	testCases := map[string]Expected{
		"/v1/conesearch": {400, map[string]string{
			"Field":    "RA",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/v1/conesearch?ra=1": {400, map[string]string{
			"Field":    "Dec",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/v1/conesearch?ra=1&dec=1": {400, map[string]string{
			"Field":    "radius",
			"Reason":   "Could not parse float.",
			"ErrValue": "",
		}},
		"/v1/conesearch?ra=1&dec=1&radius=1":                 {204, nil},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=a":       {204, nil},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=allwise": {204, nil},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=erosita": {204, nil},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=allwise&nneighbor=-1": {400, map[string]string{
			"Field":    "nneighbor",
			"ErrValue": "-1",
			"Reason":   "Nneighbor must be a positive integer",
		}},
	}

	for testPath, expected := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testPath, nil)
		router.ServeHTTP(w, req)

		require.Equalf(t, expected.Status, w.Code, "On %s", testPath)
		if w.Code == 200 || w.Code == 204 {
			continue
		}

		var result map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		require.Truef(t, maps.EqualFunc(expected.Error, result, func(a string, b any) bool {
			return a == b.(string)
		}), "On %s: values are not equal\n Expected: %v\nReceived: %v", testPath, expected.Error, result)
	}
}

func TestConesearch(t *testing.T) {
	beforeTest(t)

	// insert allwise mastercat
	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	stdout := &strings.Builder{}

	cfg, err := app.Config(getenv)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		t.Fatalf("creating database connection: %v", err)
	}

	err = test_helpers.InsertAllwiseMastercat(100, db)
	if err != nil {
		t.Fatal(err)
	}

	for i := range 10 {
		w := httptest.NewRecorder()
		ra := i
		dec := i
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/conesearch?ra=%d&dec=%d&radius=1&catalog=allwise", ra, dec), nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var result []conesearch.MastercatResult
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("could not unmarshal response: %v\n%v\nOn ra: %v, dec: %v", err, w.Body.String(), ra, dec)
		}
		require.Len(t, result, 1, "On ra=%d, dec=%d", ra, dec)
		require.Len(t, result[0].Data, 1, "On ra=%d, dec=%d", ra, dec)
		require.Equal(t, 1, result[0].Total, "On ra=%d, dec=%d", ra, dec)
		require.Equal(t, 1, result[0].TotalInCatalog, "On ra=%d, dec=%d", ra, dec)
	}
}

func TestConesearch_NoResult(t *testing.T) {
	beforeTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/conesearch?ra=1&dec=1&radius=1&catalog=allwise", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestConesearch_NNeighbor(t *testing.T) {
	beforeTest(t)

	// insert allwise mastercat
	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	stdout := &strings.Builder{}

	cfg, err := app.Config(getenv)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		t.Fatalf("creating database connection: %v", err)
	}

	err = test_helpers.InsertAllwiseMastercat(1, db)
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.New(db)
	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	if err != nil {
		t.Fatal(fmt.Errorf("could not create healpix mapper: %w", err))
	}
	ctx := context.Background()

	point := healpix.RADec(0.0000001, 0)
	ipix := mapper.PixelAt(point)
	err = repo.InsertObject(ctx, repository.InsertObjectParams{
		ID:   "allwise-1",
		Ra:   0.0000001,
		Dec:  0,
		Ipix: ipix,
		Cat:  "allwise",
	})
	point = healpix.RADec(0.0000002, 0)
	ipix = mapper.PixelAt(point)
	err = repo.InsertObject(ctx, repository.InsertObjectParams{
		ID:   "allwise-2",
		Ipix: ipix,
		Ra:   0.0000002,
		Dec:  0,
		Cat:  "allwise",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/conesearch?ra=0&dec=0&radius=1&catalog=allwise&nneighbor=5", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result []conesearch.MastercatResult
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	require.Len(t, result[0].Data, 3)
	for i := range 3 {
		require.Equal(t, fmt.Sprintf("allwise-%d", i), result[0].Data[i].ID)
	}
}

func TestBulkConesearch(t *testing.T) {
	beforeTest(t)

	// insert allwise mastercat
	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	stdout := &strings.Builder{}

	cfg, err := app.Config(getenv)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	logger := app.ServiceLogger(getenv, stdout)
	slog.SetDefault(logger)

	db, err := app.ServiceDatabase(cfg)
	if err != nil {
		t.Fatalf("creating database connection: %v", err)
	}

	err = test_helpers.InsertAllwiseMastercat(100, db)
	if err != nil {
		t.Fatal(err)
	}

	for i := range 10 {
		w := httptest.NewRecorder()
		ra := make([]float64, 10)
		dec := make([]float64, 10)

		for j := range 10 {
			// set ra and dec from 0,0 to 99,90
			ra[j] = float64(i*10 + j)
			dec[j] = float64((i*10 + j) % 90)
		}

		jsonBody := map[string]any{
			"ra":        ra,
			"dec":       dec,
			"radius":    1,
			"catalog":   "allwise",
			"nneighbor": 100,
		}

		bbody, err := json.Marshal(jsonBody)
		require.NoError(t, err)

		body := bytes.NewReader(bbody)

		// create the request using the body
		req, err := http.NewRequest("POST", "/v1/bulk-conesearch", body)
		require.NoError(t, err)

		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code, "Request: %v | Response: %v", jsonBody, w.Body.String())

		var result []repository.Mastercat
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("could not unmarshal response: %v\n%v\nOn ra: %v, dec: %v", err, w.Body.String(), ra, dec)
		}
		require.GreaterOrEqualf(t, len(result), 1, "On ra=%d, dec=%d", ra, dec)
	}
}

func TestConesearch_TruncatedMatchCounts(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	// The allwise objects are the nearest matches; the vlass objects are
	// within the radius but cut by nneighbor, so the vlass catalog produces
	// no group while still contributing to total.
	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: "allwise-1", Ra: 0.0001, Dec: 0, Cat: "allwise"},
		{ID: "allwise-2", Ra: 0.0002, Dec: 0, Cat: "allwise"},
		{ID: "vlass-1", Ra: 0.05, Dec: 0, Cat: "vlass"},
		{ID: "vlass-2", Ra: 0.06, Dec: 0, Cat: "vlass"},
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/conesearch?ra=0&dec=0&radius=3600&catalog=all&nneighbor=2", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result []conesearch.MastercatResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))

	require.Len(t, result, 1)
	require.Equal(t, "allwise", result[0].Catalog)
	require.Equal(t, 4, result[0].Total)
	require.Equal(t, 2, result[0].TotalInCatalog)
	require.Len(t, result[0].Data, 2)
}

func TestConesearch_MetadataMatchCounts(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	// Objects are separated by fractions of an arcsecond so the metadata
	// path retrieves them from a small pixel set.
	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: "allwise-0", Ra: 0, Dec: 0, Cat: "allwise"},
		{ID: "allwise-1", Ra: 0.0002, Dec: 0, Cat: "allwise"},
		{ID: "allwise-2", Ra: 0.0004, Dec: 0, Cat: "allwise"},
	})
	require.NoError(t, test_helpers.InsertAllwiseMetadata(3, db))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/conesearch?ra=0&dec=0&radius=2&catalog=all&nneighbor=2&getMetadata=true", nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var result []conesearch.MetadataResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))

	// Three metadata objects lie within two arcseconds; nneighbor keeps two.
	require.Len(t, result, 1)
	require.Equal(t, "AllWISE", result[0].Catalog)
	require.Equal(t, 3, result[0].Total)
	require.Equal(t, 3, result[0].TotalInCatalog)
	require.Len(t, result[0].Data, 2)
}

func TestBulkConesearch_MatchCounts(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: "allwise-1", Ra: 0.0001, Dec: 0, Cat: "allwise"},
		{ID: "allwise-2", Ra: 0.0002, Dec: 0, Cat: "allwise"},
		{ID: "vlass-1", Ra: 0.05, Dec: 0, Cat: "vlass"},
		{ID: "vlass-2", Ra: 0.06, Dec: 0, Cat: "vlass"},
	})

	// Position 0 has four in-radius matches, truncated to the two nearest
	// allwise objects (the fully cut vlass catalog still counts in total).
	// Position 1 has no matches, so it contributes no rows.
	bbody, err := json.Marshal(map[string]any{
		"ra":        []float64{0, 50},
		"dec":       []float64{0, 50},
		"radius":    3600,
		"catalog":   "all",
		"nneighbor": 2,
	})
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/v1/bulk-conesearch", bytes.NewReader(bbody))
	require.NoError(t, err)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var result []conesearch.MastercatResult
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))

	require.Len(t, result, 2)
	ids := make([]string, 0, len(result))
	for _, row := range result {
		require.Equal(t, 0, row.Index)
		require.Equal(t, "allwise", row.Catalog)
		require.Equal(t, 4, row.Total)
		require.Equal(t, 2, row.TotalInCatalog)
		require.Len(t, row.Data, 1)
		ids = append(ids, row.Data[0].ID)
	}
	require.ElementsMatch(t, []string{"allwise-1", "allwise-2"}, ids)
}

func TestBulkConesearch_InvalidBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/bulk-conesearch", strings.NewReader(`{"ra": "not-an-array"}`))
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, "body", result["Field"])
	require.NotEmpty(t, result["Reason"])
	require.Contains(t, result, "ErrValue")
}

// newConesearchTestDB opens the shared test database and registers its cleanup.
func newConesearchTestDB(t *testing.T) *sql.DB {
	t.Helper()

	getenv := func(key string) string {
		switch key {
		case "LOG_LEVEL":
			return "debug"
		case "CONFIG_PATH":
			return configPath
		default:
			return ""
		}
	}
	stdout := &strings.Builder{}

	cfg, err := app.Config(getenv)
	require.NoError(t, err)

	slog.SetDefault(app.ServiceLogger(getenv, stdout))

	db, err := app.ServiceDatabase(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	return db
}

// insertMastercatObjects inserts objects with their healpix pixel computed at
// the mapper resolution used by the test service.
func insertMastercatObjects(t *testing.T, db *sql.DB, objects []repository.Mastercat) {
	t.Helper()

	mapper, err := healpix.NewHEALPixMapper(18, healpix.Nest)
	require.NoError(t, err)

	repo := repository.New(db)
	for _, obj := range objects {
		err := repo.InsertObject(context.Background(), repository.InsertObjectParams{
			ID:   obj.ID,
			Ra:   obj.Ra,
			Dec:  obj.Dec,
			Cat:  obj.Cat,
			Ipix: mapper.PixelAt(healpix.RADec(obj.Ra, obj.Dec)),
		})
		require.NoError(t, err)
	}
}
