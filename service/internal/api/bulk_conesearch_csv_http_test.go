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
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/stretchr/testify/require"
)

func TestBulkConesearch_CSV(t *testing.T) {
	beforeTest(t)
	db := newConesearchTestDB(t)
	insertMastercatObjects(t, db, []repository.Mastercat{
		{ID: "allwise-0", Ra: 0, Dec: 0, Cat: "allwise"},
	})

	const csvBody = "RA,dec,name\n0,0,annotated\n"
	tests := map[string]func(t *testing.T) (body, contentType string){
		"raw body": func(t *testing.T) (string, string) {
			return csvBody, "text/csv"
		},
		"multipart file": func(t *testing.T) (string, string) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", "coordinates.csv")
			require.NoError(t, err)
			_, err = part.Write([]byte(csvBody))
			require.NoError(t, err)
			require.NoError(t, writer.Close())
			return body.String(), writer.FormDataContentType()
		},
	}

	for name, build := range tests {
		t.Run(name, func(t *testing.T) {
			body, contentType := build(t)
			req, err := http.NewRequest(
				"POST",
				"/v1/bulk-conesearch?radius=1&catalog=allwise&nneighbor=1",
				strings.NewReader(body),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code, w.Body.String())

			var result []conesearch.MastercatResult
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
			require.Len(t, result, 1)
			require.Equal(t, 0, result[0].Index)
			require.Equal(t, "allwise-0", result[0].Data[0].ID)
		})
	}
}

func TestBulkConesearch_CSVMultipartMissingFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("ra", "0"))
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(
		"POST",
		"/v1/bulk-conesearch?radius=1",
		strings.NewReader(body.String()),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, "file", result["Field"])
}

func TestBulkConesearch_CSVInvalidHeader(t *testing.T) {
	req, err := http.NewRequest(
		"POST",
		"/v1/bulk-conesearch?radius=1",
		strings.NewReader("ra,name\n1,first\n"),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/csv")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, "header", result["Field"])
}

func TestBulkConesearch_MissingRadius(t *testing.T) {
	req, err := http.NewRequest(
		"POST",
		"/v1/bulk-conesearch",
		strings.NewReader(`{"ra":[0],"dec":[0]}`),
	)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, "radius", result["Field"])
}
