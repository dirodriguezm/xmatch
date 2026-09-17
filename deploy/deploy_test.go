package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
				Assets:  []Asset{{Name: "main", BrowserDownloadURL: "test/url", Digest: "sha256:abc"}},
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
	require.Equal(t, asset, Asset{Name: "main", BrowserDownloadURL: "test/url", Digest: "sha256:abc"})
	require.Equal(t, tag, "test")
}

func TestDownloadRelease(t *testing.T) {
	path := "/repos/dirodriguezm/xmatch/releases/assets/001"
	content := []byte("test content")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(content)
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client := ts.Client()
	tmpPath := t.TempDir()
	dest := filepath.Join(tmpPath, "test")
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(content))

	err := downloadRelease(dest, client, ts.URL+path, digest)

	require.NoError(t, err)
	require.FileExists(t, dest)
	require.NoFileExists(t, dest+".part")

	got, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, content, got)
}

func TestDownloadReleaseDigestMismatch(t *testing.T) {
	path := "/repos/dirodriguezm/xmatch/releases/assets/001"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("test content"))
		require.NoError(t, err)
	}))
	defer ts.Close()

	tmpPath := t.TempDir()
	dest := filepath.Join(tmpPath, "test")

	err := downloadRelease(dest, ts.Client(), ts.URL+path, "sha256:deadbeef")

	require.ErrorContains(t, err, "sha256 mismatch")
	require.NoFileExists(t, dest)
	require.NoFileExists(t, dest+".part")
}

func TestDownloadReleaseServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	tmpPath := t.TempDir()
	dest := filepath.Join(tmpPath, "test")

	err := downloadRelease(dest, ts.Client(), ts.URL+"/missing", "")

	require.ErrorContains(t, err, "Status code: 404")
	require.NoFileExists(t, dest)
	require.NoFileExists(t, dest+".part")
}

func TestDownloadReleaseExistingFile(t *testing.T) {
	tmpPath := t.TempDir()
	dest := filepath.Join(tmpPath, "test")
	require.NoError(t, os.WriteFile(dest, []byte("already here"), 0o644))

	err := downloadRelease(dest, http.DefaultClient, "http://127.0.0.1:1/should-not-be-called", "")

	require.NoError(t, err)
	got, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, []byte("already here"), got)
}

func TestVerifyDigest(t *testing.T) {
	content := []byte("test content")
	sum := sha256.Sum256(content)
	valid := fmt.Sprintf("sha256:%x", sum)

	require.NoError(t, verifyDigest(sum[:], ""))
	require.NoError(t, verifyDigest(sum[:], "md5:whatever"))
	require.NoError(t, verifyDigest(sum[:], valid))
	require.NoError(t, verifyDigest(sum[:], fmt.Sprintf("sha256:%X", sum)))
	require.ErrorContains(t, verifyDigest(sum[:], "sha256:deadbeef"), "sha256 mismatch")
}

func TestReplaceSymlink(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first")
	second := filepath.Join(dir, "second")
	require.NoError(t, os.WriteFile(first, []byte("first"), 0o755))
	require.NoError(t, os.WriteFile(second, []byte("second"), 0o755))

	link := filepath.Join(dir, "prod")
	require.NoError(t, os.WriteFile(link, []byte("stale regular file"), 0o644))

	require.NoError(t, replaceSymlink(first, link))
	target, err := filepath.EvalSymlinks(link)
	require.NoError(t, err)
	require.Equal(t, first, target)

	require.NoError(t, replaceSymlink(second, link))
	target, err = filepath.EvalSymlinks(link)
	require.NoError(t, err)
	require.Equal(t, second, target)

	require.NoFileExists(t, link+".tmp")
}
