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

package partition_writer_test

import (
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/di"
	partition_writer "github.com/dirodriguezm/xmatch/service/internal/preprocessor/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/suite"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

type IntegrationTestSuite struct {
	suite.Suite
	w          *partition_writer.PartitionWriter
	configPath string
	tmpDir     string
	ctr        container.Container
}

func (s *IntegrationTestSuite) SetupSuite() {
	os.Setenv("LOG_LEVEL", "debug")

	// create a config file
	tmpDir, err := os.MkdirTemp("", "partition_writer_integration_test_*")
	if err != nil {
		slog.Error("could not make temp dir")
		panic(err)
	}
	s.tmpDir = tmpDir

	s.configPath = partition_writer.CreateTestConfig()

	s.ctr = di.BuildPreprocessorContainer()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	os.Remove(s.configPath)
	os.Remove(s.tmpDir)
}

func (s *IntegrationTestSuite) SetupTest() {
	err := s.ctr.Resolve(&s.w)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestActor() {
	s.w.Start()

	messages := []writer.WriterInput[repository.InputSchema]{
		{
			Rows: []repository.InputSchema{
				partition_writer.NewTestRow("1", 1, "1"),
				partition_writer.NewTestRow("2", 2, "2"),
				partition_writer.NewTestRow("3", 3, "3"),
			},
		},
	}

	for _, msg := range messages {
		s.w.InboxChannel <- msg
	}
	close(s.w.InboxChannel)
	s.w.Done()

	// check the parquet file
	var cfg *config.Config
	err := s.ctr.Resolve(&cfg)
	s.Require().NoError(err)

	// read the partition file
	dir := cfg.Preprocessor.PartitionWriter.BaseDir
	s.T().Log(os.ReadDir(dir))
	result := s.read_helper(path.Join(dir, "0", "001.parquet"))
	s.Require().Len(result, 3)
	s.Require().Equal("1", *result[0].Id)
	s.Require().Equal("2", *result[1].Id)
	s.Require().Equal("3", *result[2].Id)
}

func (s *IntegrationTestSuite) read_helper(file string) []partition_writer.TestInputSchema {
	s.T().Helper()

	fr, err := local.NewLocalFileReader(file)
	s.Require().NoError(err, "could not create local file reader")

	pr, err := reader.NewParquetReader(fr, new(partition_writer.TestInputSchema), 4)
	s.Require().NoError(err, "could not create parquet reader")

	num := int(pr.GetNumRows())

	rows := make([]partition_writer.TestInputSchema, num)
	err = pr.Read(&rows)
	s.Require().NoError(err, "could not read rows")

	pr.ReadStop()
	err = fr.Close()
	s.Require().NoError(err, "could not close file reader")

	return rows
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
