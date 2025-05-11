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

	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/suite"
)

type WorkerSuite struct {
	suite.Suite
}

func TestWorkerSuite(t *testing.T) {
	suite.Run(t, new(WorkerSuite))
}

func (s *WorkerSuite) TestReadDirectory() {
	dirChan := make(chan string, 1)
	worker := NewWorker(dirChan, config.TestSchema, make(chan Records, 1))
	dir := s.T().TempDir()

	s.T().Run("empty directory", func(t *testing.T) {
		records, err := worker.readDirectory(dir)
		s.Nil(err)

		s.Len(records, 0)
	})

	s.T().Run("with file", func(t *testing.T) {
		file, err := os.Create(path.Join(dir, "test.parquet"))
		s.Require().Nil(err)
		defer func() {
			if err := file.Close(); err != nil {
				t.Error(err)
			}
			os.Remove(file.Name())
		}()
		writeParquet(s.T(), file, []TestInputSchema{{Id: stringPtr("test")}})
		records, err := worker.readDirectory(dir)
		s.Nil(err)

		s.Len(records, 1)
		s.Equal(records[0].GetId(), "test")
	})
}

func (s *WorkerSuite) TestGroupByOid() {
	testCases := []struct {
		records  []TestInputSchema
		expected map[string]Records
	}{
		{
			records: []TestInputSchema{{Id: stringPtr("test")}, {Id: stringPtr("test")}},
			expected: map[string]Records{
				"test": {
					&TestInputSchema{Id: stringPtr("test")},
					&TestInputSchema{Id: stringPtr("test")},
				},
			},
		},
		{
			records:  []TestInputSchema{},
			expected: map[string]Records{},
		},
		{
			records: []TestInputSchema{{Id: stringPtr("test")}, {Id: stringPtr("test")}, {Id: stringPtr("test2")}},
			expected: map[string]Records{
				"test": {
					&TestInputSchema{Id: stringPtr("test")},
					&TestInputSchema{Id: stringPtr("test")},
				},
				"test2": {
					&TestInputSchema{Id: stringPtr("test2")},
				},
			},
		},
	}

	worker := NewWorker(make(chan string, 1), config.TestSchema, make(chan Records, 1))
	for _, tc := range testCases {
		records := make(Records, len(tc.records))
		for i := 0; i < len(tc.records); i++ {
			records[i] = &tc.records[i]
		}
		groupedRecords := worker.groupByOid(records)

		s.Equal(tc.expected, groupedRecords)
	}
}

func (s *WorkerSuite) TestStart() {
	dir := s.T().TempDir()

	s.T().Run("empty directory", func(t *testing.T) {
		dirChan := make(chan string, 1)
		output := make(chan Records, 1)
		worker := NewWorker(dirChan, config.TestSchema, output)
		go worker.Start()
		dirChan <- dir
		close(dirChan)

		records := <-output
		s.Len(records, 0)
	})

	s.T().Run("with file", func(t *testing.T) {
		file, err := os.Create(path.Join(dir, "test.parquet"))
		s.Require().Nil(err)
		defer func() {
			if err := file.Close(); err != nil {
				t.Error(err)
			}
			os.Remove(file.Name())
		}()
		writeParquet(s.T(), file, []TestInputSchema{{Id: stringPtr("test")}})

		dirChan := make(chan string, 1)
		output := make(chan Records, 1)
		worker := NewWorker(dirChan, config.TestSchema, output)
		go worker.Start()
		dirChan <- dir
		close(dirChan)

		records := <-output
		s.Len(records, 1)
		s.Equal(records[0].GetId(), "test")
	})
}
