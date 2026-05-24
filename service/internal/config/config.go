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

// Package config contains the configuration for the service and catalog indexer.
package config

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var embeddedConfig embed.FS

type Config struct {
	CatalogIndexer CatalogIndexerConfig `yaml:"catalog_indexer"`
	Service        ServiceConfig        `yaml:"service"`
	Preprocessor   PreprocessorConfig   `yaml:"preprocessor"`
}

type CatalogIndexerConfig struct {
	Database       DatabaseConfig `yaml:"database"`
	Source         SourceConfig   `yaml:"source"`
	Reader         ReaderConfig   `yaml:"reader"`
	Indexer        IndexerConfig  `yaml:"indexer"`
	IndexerWriter  WriterConfig   `yaml:"indexer_writer"`
	MetadataWriter WriterConfig   `yaml:"metadata_writer"`
	ChannelSize    int            `yaml:"channel_size"`
}

type PreprocessorConfig struct {
	Source        SourceConfig        `yaml:"source"`
	Reader        ReaderConfig        `yaml:"reader"`
	ReducerWriter ReducerWriterConfig `yaml:"reducer_writer"`
}

type ReducerWriterConfig struct {
	WriterConfig `yaml:",inline"`
	BatchSize    int `yaml:"batch_size"`
}

type SourceConfig struct {
	Url         string `yaml:"url"`
	Type        string `yaml:"type"`
	CatalogName string `yaml:"catalog_name"`
	Nside       int    `yaml:"nside"`
	Metadata    bool   `yaml:"metadata"`
}

type ReaderConfig struct {
	BatchSize int    `yaml:"batch_size"`
	Type      string `yaml:"type"`

	// CSV config
	Header          []string `yaml:"header"`
	FirstLineHeader bool     `yaml:"first_line_header"`
	Comment         string   `yaml:"comment"`
}

type IndexerConfig struct {
	OrderingScheme string `yaml:"ordering_scheme"`
	Nside          int    `yaml:"nside"`
}

type WriterConfig struct {
	Type string `yaml:"type"`

	// parquet config
	OutputFile string `yaml:"output_file"`
}

type ServiceConfig struct {
	Host                    string                  `yaml:"host"`
	BasePath                string                  `yaml:"base_path"`
	Database                DatabaseConfig          `yaml:"database"`
	BulkChunkSize           int                     `yaml:"bulk_chunk_size"`
	MaxBulkConcurrency      int                     `yaml:"max_bulk_concurrency"`
	LightcurveServiceConfig LightcurveServiceConfig `yaml:"lightcurve_service"`
}

type DatabaseConfig struct {
	Url string `yaml:"url"`
}

type LightcurveServiceConfig struct {
	NeowiseConfig NeowiseConfig `yaml:"neowise"`
	ZtfDrConfig   ZtfDrConfig   `yaml:"ztf_dr"`
}

type NeowiseConfig struct {
	UseIdFilter   bool `yaml:"use_id_filter"`
	UseCntrFilter bool `yaml:"use_cntr_filter"`
}

type ZtfDrConfig struct {
	UseIdFilter bool `yaml:"use_id_filter"`
}

func Load(getEnv func(string) string) (Config, error) {
	defaultConfig, err := loadDefaultConfig()
	if err != nil {
		return Config{}, err
	}

	customConfigPath := getEnv("CONFIG_PATH")
	if customConfigPath == "" {
		return defaultConfig, nil
	}

	data, err := embeddedConfig.ReadFile("config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("reading embedded config: %w", err)
	}

	customData, err := os.ReadFile(customConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("loading custom config: reading config file: %w", err)
	}

	customConfig, err := mergeYAMLConfig(data, customData)
	if err != nil {
		return Config{}, fmt.Errorf("merging config: %w", err)
	}

	return customConfig, nil
}

func mergeConfig(defaultConfig Config, customConfig Config) (Config, error) {
	result, err := utils.MergeStructs(defaultConfig, customConfig)
	if err != nil {
		return Config{}, err
	}
	return result, err
}

func loadDefaultConfig() (Config, error) {
	return loadEmbeddedConfig()
}

func loadEmbeddedConfig() (Config, error) {
	data, err := embeddedConfig.ReadFile("config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("reading embedded config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing embedded config: %w", err)
	}

	return cfg, nil
}

func mergeYAMLConfig(defaultData, customData []byte) (Config, error) {
	defaultMap, err := unmarshalConfigMap(defaultData)
	if err != nil {
		return Config{}, fmt.Errorf("parsing embedded config: %w", err)
	}

	customMap, err := unmarshalConfigMap(customData)
	if err != nil {
		return Config{}, fmt.Errorf("parsing config file: %w", err)
	}

	mergedMap := mergeConfigMaps(defaultMap, customMap)
	mergedData, err := yaml.Marshal(mergedMap)
	if err != nil {
		return Config{}, fmt.Errorf("marshaling merged config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(mergedData, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing merged config: %w", err)
	}

	return cfg, nil
}

func unmarshalConfigMap(data []byte) (map[string]any, error) {
	if strings.TrimSpace(string(data)) == "" {
		return map[string]any{}, nil
	}

	var cfg map[string]any
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return map[string]any{}, nil
	}

	return cfg, nil
}

func mergeConfigMaps(defaultMap, customMap map[string]any) map[string]any {
	merged := make(map[string]any, len(defaultMap)+len(customMap))
	for key, value := range defaultMap {
		merged[key] = value
	}

	for key, customValue := range customMap {
		customNested, customIsMap := customValue.(map[string]any)
		defaultNested, defaultIsMap := merged[key].(map[string]any)
		if customIsMap && defaultIsMap {
			merged[key] = mergeConfigMaps(defaultNested, customNested)
			continue
		}

		merged[key] = customValue
	}

	return merged
}

func LoadFile(path string) (Config, error) {
	slog.Info("Loading configuration", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}
