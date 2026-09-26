package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		panic(fmt.Errorf("could not get release data from url %s. %w", url, err))
	}

	binaryPath := filepath.Join(BINARIES_PATH, tag)
	err = downloadRelease(binaryPath, client, assetData.BrowserDownloadURL, assetData.Digest)
	if err != nil {
		panic(fmt.Errorf("could not download new binary from the release. %w", err))
	}
	slog.Info("Downloaded binary from release", "tag", tag)

	var previousBinary string
	if _, err = os.Stat(PROD_BINARY); !os.IsNotExist(err) {
		previousBinary, err = filepath.EvalSymlinks(PROD_BINARY)
		if err != nil {
			panic(fmt.Errorf("could not resolve symlink for %s", PROD_BINARY))
		}
		slog.Info("Binaries", "previous binary", previousBinary, "new binary", binaryPath)
	}

	if previousBinary == binaryPath {
		slog.Info("Current production binary is the latest", "binary", previousBinary)
		os.Exit(0)
	}

	// Promote binary to prod
	err = replaceSymlink(binaryPath, PROD_BINARY)
	if err != nil {
		panic(fmt.Errorf("could not create symlink %s -> %s. %w", binaryPath, PROD_BINARY, err))
	}
	slog.Info("Promoted binary to prod", "binary", binaryPath)

	// Make binary executable
	err = os.Chmod(PROD_BINARY, 0755)
	if err != nil {
		panic(fmt.Errorf("could not make binary executable %s. %w", PROD_BINARY, err))
	}

	// Restart every instance of the service (systemd)
	slog.Info("Restarting services")
	err = restartServiceInstances(instances)
	if err != nil {
		slog.Error("Failed to restart services", "error", err)
		slog.Info("Rolling back")
		// Rollback: restore the previous symlink, or remove it on a first deploy
		var rollbackErr error
		if previousBinary == "" {
			rollbackErr = os.Remove(PROD_BINARY)
		} else {
			rollbackErr = replaceSymlink(previousBinary, PROD_BINARY)
		}
		if rollbackErr != nil {
			slog.Error("Failed to rollback symlink", "error", rollbackErr)
		}
		panic(fmt.Errorf("could not restart services: %w", err))
	}

	// Successful restart means we are done
	slog.Info("Successfully restarted services")
}

// Downloads a binary file from the release URL (github releases) and verifies
// its sha256 digest when the release API provides one.
// The download goes to a temporary file that is renamed into place only after
// a complete transfer, so a failed download never leaves a truncated binary.
func downloadRelease(path string, client *http.Client, url, digest string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		slog.Info("Release binary already exists. Skipping download")
		return nil
	}
	// Make the GET request to download the executable asset
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("could not download file from url %s. %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not download file from url %s. Status code: %d", url, resp.StatusCode)
	}

	// Create the temporary output file
	tmpPath := path + ".part"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("could not create file %s. %w", tmpPath, err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	// Copy the response body to file while hashing it
	hasher := sha256.New()
	_, err = io.Copy(io.MultiWriter(out, hasher), resp.Body)
	if err != nil {
		_ = out.Close()
		return fmt.Errorf("could not copy file from url %s to %s. %w", url, tmpPath, err)
	}
	if err = out.Close(); err != nil {
		return fmt.Errorf("could not write file %s. %w", tmpPath, err)
	}

	if err = verifyDigest(hasher.Sum(nil), digest); err != nil {
		return fmt.Errorf("downloaded file from url %s failed verification. %w", url, err)
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("could not move %s to %s. %w", tmpPath, path, err)
	}

	return nil
}

// verifyDigest compares the sha256 sum of a download against the digest
// reported by the GitHub releases API ("sha256:<hex>"). Verification is
// skipped when the API does not provide a supported digest.
func verifyDigest(sum []byte, digest string) error {
	if digest == "" {
		slog.Warn("Release asset has no digest, skipping checksum verification")
		return nil
	}
	algo, expected, ok := strings.Cut(digest, ":")
	if !ok || algo != "sha256" {
		slog.Warn("Unsupported release asset digest, skipping checksum verification", "digest", digest)
		return nil
	}
	actual := hex.EncodeToString(sum)
	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf("sha256 mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

// replaceSymlink atomically points link at target, replacing any existing
// symlink or file. os.Symlink alone fails when link already exists.
func replaceSymlink(target, link string) error {
	tmp := link + ".tmp"
	if err := os.Remove(tmp); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Symlink(target, tmp); err != nil {
		return err
	}
	return os.Rename(tmp, link)
}

// Obtains Asset metadata and Tag from the release.
// Parses the releases API to get the tag and asset URL.
// The actual binary is downloaded from the resulting URL,
// so this function is an intermediate step for downloading the binary.
func getAssetData(client *http.Client, url string) (Asset, string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return Asset{}, "", fmt.Errorf("could not download file from url %s. %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Asset{}, "", fmt.Errorf("could not download file from url %s. Status code: %d", url, resp.StatusCode)
	}

	var result ReleaseResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Asset{}, "", fmt.Errorf("could not read response from url %s. %w", url, err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return Asset{}, "", fmt.Errorf("could not parse response from url %s. %w\nData: %s", url, err, body)
	}

	for _, asset := range result.Assets {
		if asset.Name == "main" {
			return asset, result.TagName, nil
		}
	}

	return Asset{}, "", fmt.Errorf("could not find main binary in release %s", url)
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
	Digest             string `json:"digest"`
}
