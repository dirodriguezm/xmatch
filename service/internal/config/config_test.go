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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		validate func(*testing.T, *Config)
	}{
		{
			name: "catalog indexer complete config",
			input: `
catalog_indexer:
  database:
    url: "sqlite3://test.db"
  source:
    url: "file:test.csv"
    catalog_name: "test"
    ra_col: "ra"
    dec_col: "dec"
    oid_col: "oid"
    nside: 18
  reader:
    batch_size: 500
    type: "csv"
  indexer:
    ordering_scheme: "nested"
  indexer_writer:
    type: "sqlite"
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, &DatabaseConfig{Url: "sqlite3://test.db"}, cfg.CatalogIndexer.Database)
				require.Equal(t, &SourceConfig{
					Url:         "file:test.csv",
					CatalogName: "test",
					Nside:       18,
				}, cfg.CatalogIndexer.Source)
				require.Equal(t, &ReaderConfig{BatchSize: 500, Type: "csv"}, cfg.CatalogIndexer.Reader)
				require.Equal(t, &IndexerConfig{OrderingScheme: "nested"}, cfg.CatalogIndexer.Indexer)
				require.Equal(t, &WriterConfig{Type: "sqlite"}, cfg.CatalogIndexer.IndexerWriter)
			},
		},
		{
			name:  "catalog indexer with empty config",
			input: "",
			validate: func(t *testing.T, cfg *Config) {
				require.Nil(t, cfg.CatalogIndexer)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dir := t.TempDir()
			configPath := filepath.Join(dir, "config.yml")
			err := os.WriteFile(configPath, []byte(testCase.input), 0644)
			require.NoError(t, err)

			cfg, err := LoadFile(configPath)

			require.NoError(t, err)
			if testCase.validate != nil {
				testCase.validate(t, cfg)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func(t *testing.T) (cleanup func(), getEnv func(string) string)
		wantErr  bool
	}{
		{
			name:    "load config path from env variable",
			wantErr: false,
			setupEnv: func(t *testing.T) (cleanup func(), getEnv func(string) string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "custom_path.yaml")
				err := os.WriteFile(configPath, []byte(""), 0644)
				require.NoError(t, err)

				getEnv = func(key string) string {
					switch key {
					case "CONFIG_PATH":
						return configPath
					default:
						return ""
					}
				}

				cleanup = func() {
					os.Remove(configPath)
				}

				return cleanup, getEnv
			},
		},
		{
			name:    "load config path from default location",
			wantErr: false,
			setupEnv: func(t *testing.T) (cleanup func(), getEnv func(string) string) {
				tempDir := t.TempDir()
				t.Log("asdfsafd")
				t.Log(tempDir)
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(""), 0644)

				getEnv = func(key string) string {
					switch key {
					case "CONFIG_PATH":
						return filepath.Join(tempDir, "config.yaml")
					default:
						return ""
					}
				}

				cleanup = func() {
					os.Remove(configPath)
				}

				return cleanup, getEnv
			},
		},
		{
			name:    "no config found",
			wantErr: true,
			setupEnv: func(t *testing.T) (cleanup func(), getEnv func(string) string) {
				getEnv = func(key string) string {
					switch key {
					case "CONFIG_PATH":
						return t.TempDir()
					default:
						return ""
					}
				}
				cleanup = func() {
				}
				return cleanup, getEnv
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup, getEnv := test.setupEnv(t)
			defer cleanup()

			cfg, err := Load(getEnv)
			if test.wantErr {
				require.Error(t, err, "expected test case %s to have error", test.name)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cfg)
		})
	}
}
