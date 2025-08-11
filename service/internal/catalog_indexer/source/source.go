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

package source

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type Source struct {
	Sources       []string
	CurrentSource int
	CatalogName   string
	Nside         int
}

func NewSource(cfg *config.SourceConfig) (*Source, error) {
	slog.Debug("Creating new Source", "config", cfg)
	if !validateSourceType(cfg.Type) {
		return nil, fmt.Errorf("Can't create source with type %s.", cfg.Type)
	}
	sources, err := urlSources(cfg.Url)
	if err != nil {
		return nil, err
	}
	return &Source{
		Sources:       sources,
		CurrentSource: 0,
		CatalogName:   cfg.CatalogName,
		Nside:         cfg.Nside,
	}, nil
}

// Next reads the next source file or buffer in the list of sources.
//
// Returns:
// - An io.Reader for the next source in the list.
// - io.EOF if there are no more sources to read.
// - An error if a file cannot be opened.
//
// It increments the currentSource index after successfully opening a file.
// Buffers with "buffer" prefix are handled separately.
func (src *Source) Next() (io.Reader, error) {
	// All sources read
	if src.CurrentSource >= len(src.Sources) {
		return nil, io.EOF
	}

	// In memory source from the source string itself
	if strings.HasPrefix(src.Sources[src.CurrentSource], "buffer:") {
		content := strings.Split(src.Sources[src.CurrentSource], "buffer:")[1]
		src.CurrentSource++
		return bytes.NewBufferString(content), nil
	}

	filepath := src.Sources[src.CurrentSource]
	reader, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Could not open file: %s", filepath)
	}
	src.CurrentSource++

	return reader, nil
}

func validateSourceType(stype string) bool {
	allowedTypes := []string{"csv", "parquet"}
	return slices.Contains(allowedTypes, stype)
}

func urlSources(url string) ([]string, error) {
	if strings.HasPrefix(url, "file:") {
		parsedUrl := strings.Split(url, "file:")[1]
		return []string{parsedUrl}, nil
	}

	if strings.HasPrefix(url, "buffer:") {
		return []string{url}, nil
	}

	if strings.HasPrefix(url, "files:") {
		parsedUrl := strings.Split(url, "files:")[1]
		entries, err := os.ReadDir(parsedUrl)
		if err != nil {
			return nil, err
		}
		files := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				rdrs, err := urlSources("files:" + filepath.Join(parsedUrl, entry.Name()))
				if err != nil {
					slog.Error("Could not create reader", "entry", entry.Name())
					return nil, err
				}
				files = append(files, rdrs...)
				continue
			}
			fileUrl := filepath.Join(parsedUrl, entry.Name())
			files = append(files, fileUrl)
		}
		return files, nil
	}

	return nil, fmt.Errorf("Could not parse URL: %s", url)
}
