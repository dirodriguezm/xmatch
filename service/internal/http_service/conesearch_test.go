package httpservice_test

import (
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

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		require.Truef(t, maps.EqualFunc(expected.Error, result, func(a string, b interface{}) bool {
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

	for i := 0; i < 10; i++ {
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
	for i := 0; i < 3; i++ {
		require.Equal(t, fmt.Sprintf("allwise-%d", i), result[i].ID)
	}
}
