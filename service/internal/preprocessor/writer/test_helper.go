// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package partition_writer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type TestInputSchema struct {
	Id      *string `parquet:"name=id, type=BYTE_ARRAY"`
	Column1 *int    `parquet:"name=column1, type=INT64"`
	Column2 *string `parquet:"name=column2, type=BYTE_ARRAY"`
}

func (r TestInputSchema) GetCoordinates() (float64, float64) {
	return 0, 0
}

func (r TestInputSchema) ToMetadata() any {
	return nil
}

func (r TestInputSchema) ToMastercat(ipix int64) repository.ParquetMastercat {
	return repository.ParquetMastercat{}
}

func (r TestInputSchema) SetField(string, any) {}

func (r TestInputSchema) GetId() string {
	return *r.Id
}

func CreateTestConfig() string {
	// create a config file
	tmpDir, err := os.MkdirTemp("", "partition_writer_integration_test_*")
	if err != nil {
		slog.Error("could not make temp dir")
		panic(err)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	config := `
reader:
  batch_size: 500
preprocessor:
  source:
    url: "buffer:"
    type: "csv"
    catalog_name: "test"
  reader:
    batch_size: 500
    type: "csv"
  partition_writer:
    max_file_size: 104857600
    num_partitions: 1
    partition_levels: 1
    base_dir: "%s"
`
	config = fmt.Sprintf(config, tmpDir)
	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		slog.Error("could not write config file")
		panic(err)
	}
	os.Setenv("CONFIG_PATH", configPath)
	return configPath
}
