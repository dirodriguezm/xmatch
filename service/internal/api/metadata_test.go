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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/app"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestMetadata_FindByID(t *testing.T) {
	beforeTest(t)

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

	test_helpers.InsertAllwiseMastercat(10, db)
	test_helpers.InsertAllwiseMetadata(10, db)

	for i := range 10 {
		recorder := httptest.NewRecorder()
		endpoint := fmt.Sprintf("/v1/metadata?id=allwise-%v&catalog=allwise", i)
		req, _ := http.NewRequest("GET", endpoint, nil)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var result repository.GetAllwiseRow
		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatalf("could not unmarshal response: %v\n%v\nOn id: %v", err, recorder.Body.String(), i)
		}
		require.Equal(t, fmt.Sprintf("allwise-%v", i), result.ID)
	}
}

func TestMetadata_NoResult(t *testing.T) {
	beforeTest(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/metadata?id=allwise-1&catalog=allwise", nil)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestMetadata_Validation(t *testing.T) {
	beforeTest(t)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/metadata?id=allwise-1&catalog=invalid", nil)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestMetadata_BulkFindByID(t *testing.T) {
	beforeTest(t)

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

	test_helpers.InsertAllwiseMastercat(10, db)
	test_helpers.InsertAllwiseMetadata(10, db)

	ids := make([]string, 10)
	for i := range 10 {
		ids[i] = fmt.Sprintf("allwise-%v", i)
	}
	request := map[string]any{
		"ids":     ids,
		"catalog": "allwise",
	}

	bbody, err := json.Marshal(request)
	require.NoError(t, err)

	body := bytes.NewReader(bbody)

	req, err := http.NewRequest("POST", "/v1/bulk-metadata", body)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code, "Request: %v | Response: %v", body, w.Body.String())

	var result []repository.BulkGetAllwiseRow

	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("could not unmarshal response: %v\n%v", err, w.Body.String())
	}
	require.Len(t, result, 10)
}

// TestMetadata_ResponseKeys_Erosita pins the erosita response keys to the
// catalog table column names (snake_case) and the canonical ra/dec to the
// mastercat pair, not the catalog table's own nullable coordinates.
func TestMetadata_ResponseKeys_Erosita(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	insertErositaFixture(t, db, "erosita-1", repository.Erosita{
		Detuid:  repository.NullString{NullString: sql.NullString{String: "det-1", Valid: true}},
		Skytile: repository.NullInt64{NullInt64: sql.NullInt64{Int64: 12345, Valid: true}},
		Ra:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 100, Valid: true}},
		Dec:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 50, Valid: true}},
		MlFlux1: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 1.5e-13, Valid: true}},
	})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/metadata?id=erosita-1&catalog=erosita", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	requireSnakeCaseKeys(t, body)

	for _, key := range []string{"detuid", "skytile", "ml_flux_1", "ra", "dec"} {
		require.Containsf(t, body, key, "erosita response is missing key %q: %s", key, recorder.Body.String())
	}
	for _, key := range []string{"Detuid", "Skytile", "MlFlux1", "ra_2", "dec_2"} {
		require.NotContainsf(t, body, key, "erosita response should not contain the key %q", key)
	}

	var row repository.GetErositaRow
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &row))
	require.Equal(t, "det-1", row.Detuid.String)
	require.Equal(t, int64(12345), row.Skytile.Int64)
	require.InDelta(t, 1.5e-13, row.MlFlux1.Float64, 1e-20)
	// The fixture's mastercat pair is (10.5, -20.25); the erosita table's own
	// pair is (100, 50) and must not leak into the response.
	require.InDelta(t, 10.5, row.Ra, 1e-9)
	require.InDelta(t, -20.25, row.Dec, 1e-9)
}

// TestMetadata_ResponseKeys_Allwise pins the lowercase 2MASS mass columns and
// the mastercat coordinates carried by an allwise by-id response.
func TestMetadata_ResponseKeys_Allwise(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: "allwise-0", Ra: 12.5, Dec: -3.25, Cat: "allwise"},
	})
	require.NoError(t, test_helpers.InsertAllwiseMetadata(1, db))

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/metadata?id=allwise-0&catalog=allwise", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	requireSnakeCaseKeys(t, body)

	for _, key := range []string{"j_m_2mass", "j_msig_2mass", "h_m_2mass", "h_msig_2mass", "k_m_2mass", "k_msig_2mass", "ra", "dec"} {
		require.Containsf(t, body, key, "allwise response is missing key %q: %s", key, recorder.Body.String())
	}
	for _, key := range []string{"J_m_2mass", "K_msig_2mass"} {
		require.NotContainsf(t, body, key, "allwise response should not contain the PascalCase key %q", key)
	}

	var row repository.GetAllwiseRow
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &row))
	require.InDelta(t, 5.0, row.JM2mass.Float64, 1e-9)
	require.InDelta(t, 7.0, row.KM2mass.Float64, 1e-9)
	require.InDelta(t, 0.7, row.KMsig2mass.Float64, 1e-9)
	require.InDelta(t, 12.5, row.Ra, 1e-9)
	require.InDelta(t, -3.25, row.Dec, 1e-9)
}

// TestBulkMetadata_ResponseKeys_Gaia pins the gaia bulk response keys as
// unchanged: snake_case column names matching the gaia table.
func TestBulkMetadata_ResponseKeys_Gaia(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	insertGaiaFixture(t, db, "gaia-1")

	bbody, err := json.Marshal(map[string]any{
		"ids":     []string{"gaia-1"},
		"catalog": "gaia",
	})
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/v1/bulk-metadata", bytes.NewReader(bbody))
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

	var body []map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	require.Len(t, body, 1)
	requireSnakeCaseKeys(t, body[0])

	for _, key := range []string{"id", "source_id", "phot_g_mean_flux", "phot_g_mean_mag", "ra", "dec"} {
		require.Containsf(t, body[0], key, "gaia response is missing key %q: %s", key, recorder.Body.String())
	}
	require.NotContains(t, body[0], "PhotGMeanMag")

	var rows []repository.BulkGetGaiaRow
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &rows))
	require.Len(t, rows, 1)
	require.Equal(t, int64(424242), rows[0].SourceID.Int64)
	require.InDelta(t, 18.75, rows[0].PhotGMeanMag.Float64, 1e-9)
	require.InDelta(t, 60.0, rows[0].Ra, 1e-9)
}

// TestMetadata_ErositaNullableColumns checks that NULL catalog columns
// round-trip through the by-id response as explicit nulls instead of failing
// the scan or surfacing -9999 sentinels.
func TestMetadata_ErositaNullableColumns(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)

	// Only detuid and ml_flux_1 are set; mjd and every other nullable
	// column stay NULL in the erosita table.
	insertErositaFixture(t, db, "erosita-null", repository.Erosita{
		Detuid:  repository.NullString{NullString: sql.NullString{String: "det-null", Valid: true}},
		MlFlux1: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 1.5e-13, Valid: true}},
	})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/metadata?id=erosita-null&catalog=erosita", nil)
	require.NoError(t, err)
	router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))

	raw, ok := body["mjd"]
	require.True(t, ok, "mjd key is missing from the response: %s", recorder.Body.String())
	require.JSONEq(t, "null", string(raw))
	require.NotContains(t, recorder.Body.String(), "-9999")

	var row repository.GetErositaRow
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &row))
	require.False(t, row.Mjd.Valid, "null mjd should scan as an invalid NullFloat64")
	require.True(t, row.MlFlux1.Valid)
	require.True(t, row.Detuid.Valid)
}

// requireSnakeCaseKeys asserts that every JSON object key uses snake_case,
// i.e. no catalog leaks a PascalCase field name.
func requireSnakeCaseKeys(t *testing.T, body map[string]json.RawMessage) {
	t.Helper()
	for key := range body {
		require.Regexpf(t, `^[a-z][a-z0-9_]*$`, key, "response key %q is not snake_case", key)
	}
}

// insertGaiaFixture inserts a mastercat object and its gaia metadata row.
func insertGaiaFixture(t *testing.T, db *sql.DB, id string) {
	t.Helper()

	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: id, Ra: 60.0, Dec: -12.5, Cat: "gaia"},
	})

	repo := repository.New(db)
	err := repo.InsertGaia(context.Background(), repository.InsertGaiaParams(repository.Gaia{
		ID:            id,
		SourceID:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: 424242, Valid: true}},
		PhotGMeanFlux: repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 1234.5, Valid: true}},
		PhotGMeanMag:  repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: 18.75, Valid: true}},
	}))
	require.NoError(t, err)
}

// insertErositaFixture inserts a mastercat object and its erosita metadata
// row. Fields left at their zero value stay NULL in the erosita table.
func insertErositaFixture(t *testing.T, db *sql.DB, id string, metadata repository.Erosita) {
	t.Helper()

	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: id, Ra: 10.5, Dec: -20.25, Cat: "erosita"},
	})

	metadata.ID = id
	repo := repository.New(db)
	err := repo.InsertErosita(context.Background(), repository.InsertErositaParams(metadata))
	require.NoError(t, err)
}

func TestBulkMetadata_InvalidBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/bulk-metadata", strings.NewReader(`{"ids": "not-an-array"}`))
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, "body", result["Field"])
	require.NotEmpty(t, result["Reason"])
	require.Contains(t, result, "Value")
}
