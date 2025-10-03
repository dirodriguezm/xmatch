package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var BINARIES_PATH = "/home/drodriguez/deployment/binaries"
var PROD_BINARY = "/home/drodriguez/deployment/production/bin/prod"
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

	err = downloadRelease(filepath.Join(BINARIES_PATH, tag), client, assetData.BrowserDownloadURL)
	if err != nil {
		panic(fmt.Errorf("Could not download new binary from the release. %w", err))
	}
	slog.Info("Downloaded binary from release", "tag", tag)

	var previousBinary string
	if _, err = os.Stat(PROD_BINARY); !os.IsNotExist(err) {
		previousBinary, err = filepath.EvalSymlinks(PROD_BINARY)
		if err != nil {
			panic(fmt.Errorf("Could not resolve symlink for %s", PROD_BINARY))
		}
		slog.Info("Binaries", "previous binary", previousBinary, "new binary", filepath.Join(BINARIES_PATH, tag))
	}

	if previousBinary == filepath.Join(BINARIES_PATH, tag) {
		slog.Info("Current production binary is the latest", "binary", previousBinary)
		os.Exit(0)
	}

	// Promote binary to prod
	err = os.Symlink(filepath.Join(BINARIES_PATH, tag), PROD_BINARY)
	if err != nil {
		panic(fmt.Errorf("Could not create symlink %s -> %s. %w", filepath.Join(BINARIES_PATH, tag), PROD_BINARY, err))
	}
	slog.Info("Promoted binary to prod", "binary", filepath.Join(BINARIES_PATH, tag))

	// Make binary executable
	err = os.Chmod(PROD_BINARY, 0755)
	if err != nil {
		panic(fmt.Errorf("Could not make binary executable %s. %w", PROD_BINARY, err))
	}

	// Restart every instance of the service (systemd)
	slog.Info("Restarting services")
	err = restartServiceInstances(instances)
	if err != nil {
		slog.Error("Failed to restart services", "error", err)
		slog.Info("Rolling back")
		// Rollback: restore previous symlink
		if rollbackErr := os.Symlink(previousBinary, PROD_BINARY); rollbackErr != nil {
			slog.Error("Failed to rollback symlink", "error", rollbackErr)
		}
		panic(fmt.Errorf("Could not restart services: %w", err))
	}

	// Successful restart means we are done
	slog.Info("Successfully restarted services")
}

// Downloads a binary file from the release URL (github releases).
// It places the downloaded file in the specified directory.
func downloadRelease(path string, client *http.Client, url string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		slog.Info("Release binary already exists. Skipping download")
		return nil
	}
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

func restartServiceInstances(instances int) error {
	for i := range instances {
		serviceName := fmt.Sprintf("xmatch@%d.service", i+1)
		slog.Info("Restarting service", "service", serviceName)
		cmd := exec.Command("systemctl", "--user", "restart", serviceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to restart service %s: %w", serviceName, err)
		}
		slog.Info("Restarted service", "service", serviceName)
	}
	return nil
}

type ReleaseResponse struct {
	TagName string `json:"tag_name"`
	Assets  []Asset
}

type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
}
