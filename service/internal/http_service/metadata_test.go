package httpservice_test

import (
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

	for i := 0; i < 10; i++ {
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
