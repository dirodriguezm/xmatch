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
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
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
		"/v1/conesearch?ra=1&dec=1&radius=1": {204, nil},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=a": {400, map[string]string{
			"Field":    "catalog",
			"Reason":   "Catalog not available",
			"ErrValue": "a",
		}},
		"/v1/conesearch?ra=1&dec=1&radius=1&catalog=allwise": {204, nil},
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
	var db *sql.DB
	ctr.Resolve(&db)
	err := test_helpers.InsertAllwiseMastercat(100, db)
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

		var result []repository.Mastercat
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("could not unmarshal response: %v\n%v\nOn ra: %v, dec: %v", err, w.Body.String(), ra, dec)
		}
		require.GreaterOrEqualf(t, len(result), 1, "On ra=%d, dec=%d", ra, dec)
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
	var db *sql.DB
	ctr.Resolve(&db)
	err := test_helpers.InsertAllwiseMastercat(1, db)
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
	_, err = repo.InsertObject(ctx, repository.InsertObjectParams{
		ID:   "allwise-1",
		Ra:   0.0000001,
		Dec:  0,
		Ipix: ipix,
		Cat:  "allwise",
	})
	point = healpix.RADec(0.0000002, 0)
	ipix = mapper.PixelAt(point)
	_, err = repo.InsertObject(ctx, repository.InsertObjectParams{
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

	var result []repository.Mastercat
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	require.Len(t, result, 3)
	for i := range 3 {
		require.Equal(t, fmt.Sprintf("allwise-%d", i), result[i].ID)
	}
}

func TestBulkConesearch(t *testing.T) {
	beforeTest(t)

	// insert allwise mastercat
	var db *sql.DB
	ctr.Resolve(&db)
	err := test_helpers.InsertAllwiseMastercat(100, db)
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
