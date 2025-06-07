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
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestMetadata_FindByID(t *testing.T) {
	beforeTest(t)

	var db *sql.DB
	ctr.Resolve(&db)
	test_helpers.InsertAllwiseMetadata(10, db)

	for i := range 10 {
		recorder := httptest.NewRecorder()
		endpoint := fmt.Sprintf("/v1/metadata?id=allwise-%v&catalog=allwise", i)
		req, _ := http.NewRequest("GET", endpoint, nil)
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var result repository.AllwiseMetadata
		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatalf("could not unmarshal response: %v\n%v\nOn id: %v", err, recorder.Body.String(), i)
		}
		require.Equal(t, fmt.Sprintf("allwise-%v", i), *result.Source_id)
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

	var db *sql.DB
	ctr.Resolve(&db)
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

	var result []repository.AllwiseMetadata

	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("could not unmarshal response: %v\n%v", err, w.Body.String())
	}
	require.Len(t, result, 10)
}
