package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CatalogIndexer *CatalogIndexerConfig `yaml:"catalog_indexer"`
	Service        *ServiceConfig        `yaml:"service"`
}

type CatalogIndexerConfig struct {
	Database *DatabaseConfig `yaml:"database"`
	Source   *SourceConfig   `yaml:"source"`
	Reader   *ReaderConfig   `yaml:"reader"`
	Indexer  *IndexerConfig  `yaml:"indexer"`
	Writer   *WriterConfig   `yaml:"writer"`
}

type SourceConfig struct {
	Url         string `yaml:"url"`
	Type        string `yaml:"type"`
	CatalogName string `yaml:"catalog_name"`
	RaCol       string `yaml:"ra_col"`
	DecCol      string `yaml:"dec_col"`
	OidCol      string `yaml:"oid_col"`
	Nside       int    `yaml:"nside"`
}

type ReaderConfig struct {
	BatchSize int    `yaml:"batch_size"`
	Type      string `yaml:"type"`

	// CSV config
	Header []string `yaml:"header"`
}

type IndexerConfig struct {
	OrderingScheme string `yaml:"ordering_scheme"`
}

type WriterConfig struct{}

type ServiceConfig struct {
	Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Url string `yaml:"url"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		locations := []string{"./config.yml", "./config.yaml"}
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