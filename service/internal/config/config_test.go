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
    url: "sqlite://test.db"
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
  writer:
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, &DatabaseConfig{Url: "sqlite://test.db"}, cfg.CatalogIndexer.Database)
				require.Equal(t, &SourceConfig{
					Url:         "file:test.csv",
					CatalogName: "test",
					RaCol:       "ra",
					DecCol:      "dec",
					OidCol:      "oid",
					Nside:       18,
				}, cfg.CatalogIndexer.Source)
				require.Equal(t, &ReaderConfig{BatchSize: 500, Type: "csv"}, cfg.CatalogIndexer.Reader)
				require.Equal(t, &IndexerConfig{OrderingScheme: "nested"}, cfg.CatalogIndexer.Indexer)
				require.Nil(t, cfg.CatalogIndexer.Writer)
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
		setupEnv func(t *testing.T) (cleanup func())
		wantErr  bool
	}{
		{
			name:    "load config path from env variable",
			wantErr: false,
			setupEnv: func(t *testing.T) (cleanup func()) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "custom_path.yaml")
				os.Setenv("CONFIG_PATH", configPath)
				err := os.WriteFile(configPath, []byte(""), 0644)
				require.NoError(t, err)

				return func() {
					os.Setenv("CONFIG_PATH", "")
					os.Remove(configPath)
				}
			},
		},
		{
			name:    "load config path from default location",
			wantErr: false,
			setupEnv: func(t *testing.T) (cleanup func()) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "config.yaml")
				os.Unsetenv("CONFIG_PATH")
				os.WriteFile(configPath, []byte(""), 0644)

				return func() {
					os.Remove(configPath)
				}
			},
		},
		{
			name:    "no config found",
			wantErr: true,
			setupEnv: func(t *testing.T) (cleanup func()) {
				os.Unsetenv("CONFIG_PATH")
				os.Setenv("CONFIG_PATH", t.TempDir())
				return func() {
					os.Unsetenv("CONFIG_PATH")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup := test.setupEnv(t)
			defer cleanup()

			cfg, err := Load()
			if test.wantErr {
				require.Error(t, err, "expected test case %s to have error", test.name)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cfg)
		})
	}
}
