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

package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CatalogIndexer *CatalogIndexerConfig `yaml:"catalog_indexer"`
	Service        *ServiceConfig        `yaml:"service"`
}

type CatalogIndexerConfig struct {
	Database        *DatabaseConfig `yaml:"database"`
	Source          *SourceConfig   `yaml:"source"`
	Reader          *ReaderConfig   `yaml:"reader"`
	Indexer         *IndexerConfig  `yaml:"indexer"`
	PartitionWriter *WriterConfig   `yaml:"partition_writer"`
	ReducerWriter   *WriterConfig   `yaml:"reducer_writer"`
	IndexerWriter   *WriterConfig   `yaml:"indexer_writer"`
	MetadataWriter  *WriterConfig   `yaml:"metadata_writer"`
}

type SourceConfig struct {
	Url         string `yaml:"url"`
	Type        string `yaml:"type"`
	CatalogName string `yaml:"catalog_name"`
	RaCol       string `yaml:"ra_col"`
	DecCol      string `yaml:"dec_col"`
	OidCol      string `yaml:"oid_col"`
	Nside       int    `yaml:"nside"`
	Metadata    bool   `yaml:"metadata"`
}

type ReaderConfig struct {
	BatchSize int    `yaml:"batch_size"`
	Type      string `yaml:"type"`

	// CSV config
	Header          []string `yaml:"header"`
	FirstLineHeader bool     `yaml:"first_line_header"`
}

type IndexerConfig struct {
	OrderingScheme string `yaml:"ordering_scheme"`
	Nside          int    `yaml:"nside"`
}

type ParquetWriterSchema int

const (
	AllwiseSchema ParquetWriterSchema = iota
	MastercatSchema
	TestSchema
)

type WriterConfig struct {
	Type string `yaml:"type"`

	// parquet config
	OutputFile string `yaml:"output_file"`
	Schema     ParquetWriterSchema
}

type ServiceConfig struct {
	Database           *DatabaseConfig `yaml:"database"`
	BulkChunkSize      int             `yaml:"bulk_chunk_size"`
	MaxBulkConcurrency int             `yaml:"max_bulk_concurrency"`
}

type DatabaseConfig struct {
	Url string `yaml:"url"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	slog.Info("Loading configuration", "path", configPath)
	if configPath == "" {
		rootPath, err := utils.FindRootModulePath(5)
		if err != nil {
			return nil, err
		}
		locations := []string{
			"./config.yml",
			"./config.yaml",
			filepath.Join(rootPath, "config.yml"),
			filepath.Join(rootPath, "config.yaml"),
		}
		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				configPath = loc
				break
			}
		}
	}
	if configPath == "" {
		return nil, fmt.Errorf("Could not find configuration file")
	}
	return LoadFile(configPath)
}

func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &cfg, nil
}
