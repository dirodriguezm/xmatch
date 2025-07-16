package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

var BINARIES_PATH = "deployment/binaries"
var PROD_BINARY = "deployment/production/bin/prod"
var RELEASE_URL = "https://api.github.com/repos/dirodriguezm/xmatch/releases/latest"

func main() {
	var instances int

	flag.IntVar(&instances, "instances", 1, "Number of instances")
	flag.Parse()

	slog.Info("Starting Deploy Process")

	client := http.DefaultClient
	deploy(instances, client, RELEASE_URL)
}

func deploy(instances int, client *http.Client, url string) {
	slog.Info("Deploying", "number of instances", instances)
	// Get asset data from release
	assetData, tag, err := getAssetData(client, url)
	if err != nil {
		panic(fmt.Errorf("Could not get release data from url %s. %w", url, err))
	}

	err = downloadRelease(filepath.Join(BINARIES_PATH, tag), client, assetData.Url)
	if err != nil {
		panic(fmt.Errorf("Could not download new binary from the release"))
	}
	slog.Info("Downloaded binary from release", "tag", tag)

	previousBinary, err := filepath.EvalSymlinks(PROD_BINARY)
	if err != nil {
		panic(fmt.Errorf("Could not resolve symlink for %s", PROD_BINARY))
	}
	slog.Info("Binaries", "previous binary", previousBinary, "new binary", filepath.Join(BINARIES_PATH, tag))

	// Promote binary to prod
	err = os.Symlink(filepath.Join(BINARIES_PATH, tag), PROD_BINARY)
	if err != nil {
		panic(fmt.Errorf("Could not create symlink %s -> %s. %w", filepath.Join(BINARIES_PATH, tag), PROD_BINARY, err))
	}
	slog.Info("Promoted binary to prod", "binary", filepath.Join(BINARIES_PATH, tag))

	// Restart every instance of the service

	// Successful restart means we are done

	// Rollback
}

// Downloads a binary file from the release URL (github releases).
// It places the downloaded file in the specified directory.
func downloadRelease(path string, client *http.Client, url string) error {
	// Make the GET request to download the executable asset
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("Could not download file from url %s. %w", url, err)
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Could not create file %s. %w", path, err)
	}
	defer out.Close()

	// Copy the response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Could not copy file from url %s to %s. %w", url, path, err)
	}

	return nil
}

// Obtains Asset metadata and Tag from the release.
// Parses the releases API to get the tag and asset URL.
// The actual binary is downloaded from the resulting URL,
// so this function is an intermediate step for downloading the binary.
func getAssetData(client *http.Client, url string) (Asset, string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return Asset{}, "", fmt.Errorf("Could not download file from url %s. %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Asset{}, "", fmt.Errorf("Could not download file from url %s. Status code: %d", url, resp.StatusCode)
	}

	var result ReleaseResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, "", fmt.Errorf("Could not read response from url %s. %w", url, err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Asset{}, "", fmt.Errorf("Could not parse response from url %s. %w\nData: %s", url, err, body)
	}

	for _, asset := range result.Assets {
		if asset.Name == "main" {
			return asset, result.TagName, nil
		}
	}

	return Asset{}, "", fmt.Errorf("Could not find main binary in release %s", url)
}

type ReleaseResponse struct {
	TagName string `json:"tag_name"`
	Assets  []Asset
}

type Asset struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func promoteBinary(binaryPath string, prodPath string) {
}
