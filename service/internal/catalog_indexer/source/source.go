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
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type SourceReader struct {
	io.Reader

	Url string
}

type Source struct {
	Reader      []SourceReader
	CatalogName string
	RaCol       string
	DecCol      string
	OidCol      string
	Nside       int
}

func NewSource(cfg *config.SourceConfig) (*Source, error) {
	slog.Debug("Creating new Source", "config", cfg)
	if !validateSourceType(cfg.Type) {
		return nil, fmt.Errorf("Can't create source with type %s.", cfg.Type)
	}
	readers, err := sourceReader(cfg.Url)
	if err != nil {
		return nil, err
	}
	return &Source{
		Reader:      readers,
		CatalogName: cfg.CatalogName,
		RaCol:       cfg.RaCol,
		DecCol:      cfg.DecCol,
		OidCol:      cfg.OidCol,
		Nside:       cfg.Nside,
	}, nil
}

func validateSourceType(stype string) bool {
	allowedTypes := []string{"csv", "parquet"}
	for _, t := range allowedTypes {
		if t == stype {
			return true
		}
	}
	return false
}

func sourceReader(url string) ([]SourceReader, error) {
	if strings.HasPrefix(url, "file:") {
		parsedUrl := strings.Split(url, "file:")[1]
		reader, err := os.Open(parsedUrl)
		if err != nil {
			return nil, err
		}
		return []SourceReader{{Reader: reader, Url: parsedUrl}}, nil
	}

	if strings.HasPrefix(url, "buffer:") {
		return []SourceReader{{Reader: &bytes.Buffer{}, Url: url}}, nil
	}

	if strings.HasPrefix(url, "files:") {
		parsedUrl := strings.Split(url, "files:")[1]
		entries, err := os.ReadDir(parsedUrl)
		if err != nil {
			return nil, err
		}
		readers := []SourceReader{}
		for _, entry := range entries {
			if entry.IsDir() {
				rdrs, err := sourceReader("files:" + filepath.Join(parsedUrl, entry.Name()))
				if err != nil {
					slog.Error("Could not create reader", "entry", entry.Name())
					return nil, err
				}
				readers = append(readers, rdrs...)
				continue
			}
			fileUrl := filepath.Join(parsedUrl, entry.Name())
			file, err := os.Open(fileUrl)
			if err != nil {
				slog.Error("Could not create reader", "entry", entry.Name())
				return nil, err
			}
			readers = append(readers, SourceReader{Reader: file, Url: fileUrl})
		}
		return readers, nil
	}

	return nil, fmt.Errorf("Could not parse URL: %s", url)
}
