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

package fits_reader

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestReadBatch(t *testing.T) {
	testData := []map[string]any{
		{
			"col1": "hola",
			"col2": 18,
		},
		{
			"col1": "adios",
			"col2": 99,
		},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestRead.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestRead.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	fitsReader, err := NewFitsReader(src, WithBatchSize(2))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)

	require.Equal(t, 2, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithDifferentDataTypes(t *testing.T) {
	testData := []map[string]any{
		{
			"string_col": "test1",
			"int_col":    int32(42),
			"float_col":  3.14159,
			"bool_col":   true,
			"int64_col":  int64(9999999999),
		},
		{
			"string_col": "test2",
			"int_col":    int32(-100),
			"float_col":  -2.71828,
			"bool_col":   false,
			"int64_col":  int64(-9999999999),
		},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithDifferentDataTypes.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithDifferentDataTypes.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	fitsReader, err := NewFitsReader(src, WithBatchSize(2))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 2, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithEmptyFile(t *testing.T) {
	testData := []map[string]any{}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithEmptyFile.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithEmptyFile.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	fitsReader, err := NewFitsReader(src, WithBatchSize(2))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.EqualError(t, err, "EOF")
	require.Equal(t, 0, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithLargeBatchSize(t *testing.T) {
	testData := []map[string]any{
		{"col1": "row1", "col2": 1},
		{"col1": "row2", "col2": 2},
		{"col1": "row3", "col2": 3},
		{"col1": "row4", "col2": 4},
		{"col1": "row5", "col2": 5},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithLargeBatchSize.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithLargeBatchSize.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	// Test with batch size larger than available rows
	fitsReader, err := NewFitsReader(src, WithBatchSize(10))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.EqualError(t, err, "EOF")
	require.Equal(t, 5, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithSmallBatchSize(t *testing.T) {
	testData := []map[string]any{
		{"col1": "row1", "col2": 1},
		{"col1": "row2", "col2": 2},
		{"col1": "row3", "col2": 3},
		{"col1": "row4", "col2": 4},
		{"col1": "row5", "col2": 5},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithSmallBatchSize.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithSmallBatchSize.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	// Test with batch size smaller than available rows
	fitsReader, err := NewFitsReader(src, WithBatchSize(2))
	require.NoError(t, err)

	// First batch
	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 2, len(rows))

	// Second batch
	rows, err = fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 2, len(rows))

	// Third batch (remaining 1 row)
	rows, err = fitsReader.ReadBatch()
	require.EqualError(t, err, "EOF")
	require.Equal(t, 1, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithNilValues(t *testing.T) {
	testData := []map[string]any{
		{
			"col1":         "row1",
			"col2":         1,
			"nullable_col": nil,
		},
		{
			"col1":         "row2",
			"col2":         2,
			"nullable_col": "not null",
		},
		{
			"col1":         "row3",
			"col2":         3,
			"nullable_col": nil,
		},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithNilValues.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithNilValues.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	fitsReader, err := NewFitsReader(src, WithBatchSize(3))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 3, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithDefaultBatchSize(t *testing.T) {
	testData := []map[string]any{
		{"col1": "row1", "col2": 1},
		{"col1": "row2", "col2": 2},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithDefaultBatchSize.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithDefaultBatchSize.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	// Test without specifying batch size (should use default)
	fitsReader, err := NewFitsReader(src)
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	// Default batch size should be 1 based on WithBatchSize logic
	require.Equal(t, 1, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithInvalidBatchSize(t *testing.T) {
	testData := []map[string]any{
		{"col1": "row1", "col2": 1},
		{"col1": "row2", "col2": 2},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithInvalidBatchSize.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithInvalidBatchSize.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	// Test with zero batch size (should be normalized to 1)
	fitsReader, err := NewFitsReader(src, WithBatchSize(0))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 1, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}

func TestReadBatchWithNegativeBatchSize(t *testing.T) {
	testData := []map[string]any{
		{"col1": "row1", "col2": 1},
		{"col1": "row2", "col2": 2},
	}

	writeFitsFile(t, "/tmp/fits_reader_test_TestReadBatchWithNegativeBatchSize.fits", testData)

	cfg := config.SourceConfig{
		Url:         "file:/tmp/fits_reader_test_TestReadBatchWithNegativeBatchSize.fits",
		Type:        "fits",
		CatalogName: "test",
	}
	src, err := source.NewSource(cfg)
	require.NoError(t, err)

	// Test with negative batch size (should be normalized to 1)
	fitsReader, err := NewFitsReader(src, WithBatchSize(-5))
	require.NoError(t, err)

	rows, err := fitsReader.ReadBatch()
	require.NoError(t, err)
	require.Equal(t, 1, len(rows))

	err = fitsReader.Close()
	require.NoError(t, err)
}
