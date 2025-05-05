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
package partition_writer

import (
	"os"
	"path"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/config"
	filesystemmanager "github.com/dirodriguezm/xmatch/service/internal/preprocessor/filesystem_manager"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/suite"
)

type ParquetStoreTestSuite struct {
	suite.Suite
	store    *ParquetStore
	tempDir  string
	dirMap   map[string]int
	testData []TestInputSchema
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func (s *ParquetStoreTestSuite) SetupTest() {
	s.tempDir = s.T().TempDir()

	fs := filesystemmanager.AFileSystemManager(s.T()).
		WithBaseDir(s.tempDir).
		WithNumLevels(1).
		WithNumPartitions(1).
		Build()

	s.store = &ParquetStore{
		fs:      &fs,
		maxSize: 1024 * 1024, // 1MB
		schema:  config.TestSchema,
	}

	s.dirMap = make(map[string]int)
	s.testData = []TestInputSchema{
		{Id: stringPtr("1"), Column1: intPtr(1), Column2: stringPtr("test1")},
		{Id: stringPtr("2"), Column1: intPtr(2), Column2: stringPtr("test2")},
	}
}

func (s *ParquetStoreTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ParquetStoreTestSuite) TestCreateNewFile() {
	// Test creating a regular file
	dir := path.Join(s.tempDir, "0")
	err := os.MkdirAll(dir, 0777)
	s.Require().NoError(err)

	rows := []repository.InputSchema{
		NewTestRow("1", 1, "test1"),
		NewTestRow("2", 2, "test2"),
	}

	fileName, err := s.store.createNewFile(dir, 1, rows, false)
	s.Require().NoError(err)
	s.Equal(path.Join(dir, "001.parquet"), fileName)

	// Verify file exists and contains correct data
	records := readParquet(s.T(), fileName)
	s.Require().Len(records, 2)
	s.Equal("test1", *records[0].Column2)
	s.Equal("test2", *records[1].Column2)
}

func (s *ParquetStoreTestSuite) TestCreateNewFileTmp() {
	// Test creating a temporary file
	dir := path.Join(s.tempDir, "0")
	err := os.MkdirAll(dir, 0777)
	s.Require().NoError(err)

	rows := []repository.InputSchema{
		NewTestRow("3", 3, "test3"),
		NewTestRow("4", 4, "test4"),
	}

	fileName, err := s.store.createNewFile(dir, 1, rows, true)
	s.Require().NoError(err)
	s.Contains(fileName, "tmp") // Temporary files should contain tmp in the name

	// Verify file exists and contains correct data
	records := readParquet(s.T(), fileName)
	s.Require().Len(records, 2)
	s.Equal("test3", *records[0].Column2)
	s.Equal("test4", *records[1].Column2)
}

func (s *ParquetStoreTestSuite) TestReuseFile() {

	dir := path.Join(s.tempDir, "0")
	s.dirMap[dir] = 1

	initialFileName := path.Join(dir, "001.parquet")
	os.Mkdir(path.Join(s.tempDir, "0"), 0777)
	initialFile, err := os.Create(initialFileName)
	s.Require().NoError(err)

	writeParquet(s.T(), initialFile, s.testData)

	// Open the file for reuse
	file, err := os.Open(initialFileName)
	s.Require().NoError(err)
	defer file.Close()

	// Test reusing the file by appending new data
	newData := []repository.InputSchema{
		NewTestRow("3", 3, "test3"),
		NewTestRow("4", 4, "test4"),
	}
	err = s.store.reuseFile(file, newData, s.dirMap)
	s.Require().NoError(err)

	// Read the final file and verify contents
	records := readParquet(s.T(), initialFileName)

	// Verify the contents
	s.Require().Len(records, 4) // Should have both initial and new records
}

func (s *ParquetStoreTestSuite) TestWrite() {
	// Setup test directory
	dir := path.Join(s.tempDir, "1")
	err := os.MkdirAll(dir, 0777)
	s.Require().NoError(err)

	// Create initial file with some data
	initialFileName := path.Join(dir, "001.parquet")
	initialFile, err := os.Create(initialFileName)
	s.Require().NoError(err)
	writeParquet(s.T(), initialFile, s.testData)

	// Initialize dirMap with existing file
	dirMap := map[string]int{
		dir: 1,
	}

	// Test data for both scenarios
	rowsToWrite := map[string][]repository.InputSchema{
		// Dir with existing small file that can be reused
		dir: {
			NewTestRow("3", 3, "test3"),
			NewTestRow("4", 4, "test4"),
		},
		// Dir that needs a new file
		path.Join(s.tempDir, "2"): {
			NewTestRow("5", 5, "test5"),
			NewTestRow("6", 6, "test6"),
		},
	}

	// Execute Write method
	dirMap, err = s.store.Write(rowsToWrite, dirMap)
	s.Require().NoError(err)

	// Verify reused file
	reusedRecords := readParquet(s.T(), initialFileName)
	s.Require().Len(reusedRecords, 4) // Original 2 + 2 new records
	s.Equal("test1", *reusedRecords[0].Column2)
	s.Equal("test2", *reusedRecords[1].Column2)
	s.Equal("test3", *reusedRecords[2].Column2)
	s.Equal("test4", *reusedRecords[3].Column2)

	// Verify new file
	newFileName := path.Join(s.tempDir, "2", "001.parquet")
	s.Require().FileExists(newFileName)
	newRecords := readParquet(s.T(), newFileName)
	s.Require().Len(newRecords, 2)
	s.Equal("test5", *newRecords[0].Column2)
	s.Equal("test6", *newRecords[1].Column2)

	// Verify dirMap was updated correctly
	s.Equal(1, dirMap[dir])                       // Should remain 1 as file was reused
	s.Equal(1, dirMap[path.Join(s.tempDir, "2")]) // Should be 1 for new directory
}

func TestParquetStoreSuite(t *testing.T) {
	suite.Run(t, new(ParquetStoreTestSuite))
}
