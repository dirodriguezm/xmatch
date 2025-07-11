package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAssetData(t *testing.T) {
	releasePath := "/repos/dirodriguezm/xmatch/releases/latest"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == releasePath {
			fakeData := ReleaseResponse{
				TagName: "test",
				Assets:  []Asset{{Name: "main", Url: "test/url"}},
			}
			jsonData, err := json.Marshal(fakeData)
			require.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonData)
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := ts.Client()
	asset, tag, err := getAssetData(client, ts.URL+releasePath)

	require.NoError(t, err)
	require.Equal(t, asset, Asset{Name: "main", Url: "test/url"})
	require.Equal(t, tag, "test")
}

func TestDownloadRelease(t *testing.T) {
	path := "/repos/dirodriguezm/xmatch/releases/assets/001"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("test content"))
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := ts.Client()
	tmpPath := t.TempDir()
	err := downloadRelease(filepath.Join(tmpPath, "test"), client, ts.URL+path)

	require.NoError(t, err)
	require.FileExists(t, filepath.Join(tmpPath, "test"))

	content, err := os.ReadFile(filepath.Join(tmpPath, "test"))
	require.NoError(t, err)

	require.Equal(t, content, []byte("test content"))
}
