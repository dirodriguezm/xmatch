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

package partition_reader

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PartitionReaderTestSuite struct {
	suite.Suite
	baseDir string
}

func (suite *PartitionReaderTestSuite) SetupSuite() {
	suite.baseDir = suite.T().TempDir()
}

func (suite *PartitionReaderTestSuite) TearDownSuite() {
}

func TestPartitionReaderTestSuite(t *testing.T) {
	suite.Run(t, new(PartitionReaderTestSuite))
}

func (suite *PartitionReaderTestSuite) TestTraversePartitions() {
	os.MkdirAll(path.Join(suite.baseDir, "a", "b", "c"), 0777)
	os.WriteFile(path.Join(suite.baseDir, "a", "b", "c", "file1"), []byte("test"), 0777)
	os.WriteFile(path.Join(suite.baseDir, "a", "b", "c", "file2"), []byte("test"), 0777)
	os.MkdirAll(path.Join(suite.baseDir, "a", "b", "d"), 0777)
	os.WriteFile(path.Join(suite.baseDir, "a", "b", "d", "file1"), []byte("test"), 0777)
	os.MkdirAll(path.Join(suite.baseDir, "a", "b", "e"), 0777)

	dirChan := make(chan string, 4)
	pr := NewPartitionReader(dirChan, nil, suite.baseDir)
	go func() {
		pr.TraversePartitions(suite.baseDir)
		close(dirChan)
	}()

	directories := make([]string, 0)
	for dir := range dirChan {
		directories = append(directories, dir)
	}

	suite.Len(directories, 2)
	suite.Equal(directories[0], path.Join(suite.baseDir, "a", "b", "c"))
	suite.Equal(directories[1], path.Join(suite.baseDir, "a", "b", "d"))
}
