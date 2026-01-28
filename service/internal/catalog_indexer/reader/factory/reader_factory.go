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

package reader_factory

import (
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	fits_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/fits"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

func ReaderFactory(
	src *source.Source,
	cfg config.ReaderConfig,
) (reader.Reader, error) {
	if cfg.BatchSize <= 0 {
		return nil, fmt.Errorf("Batch size must be greater than 0")
	}
	readerType := strings.ToLower(cfg.Type)
	switch readerType {
	case "csv":
		return csv_reader.NewCsvReader(
			src,
			csv_reader.WithHeader(cfg.Header),
			csv_reader.WithFirstLineHeader(cfg.FirstLineHeader),
			csv_reader.WithComment(cfg.Comment),
			csv_reader.WithCsvBatchSize(cfg.BatchSize),
		)
	case "parquet":
		return parquetFactory(src, cfg)
	case "fits":
		return fitsFactory(src, cfg)
	default:
		return nil, fmt.Errorf("Reader type not allowed")
	}
}

func parquetFactory(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	switch strings.ToLower(src.CatalogName) {
	case "allwise":
		return parquet_reader.NewParquetReader(
			src,
			parquet_reader.WithParquetBatchSize[repository.AllwiseInputSchema](cfg.BatchSize),
		)
	case "gaia":
		return parquet_reader.NewParquetReader(
			src,
			parquet_reader.WithParquetBatchSize[repository.GaiaInputSchema](cfg.BatchSize),
		)
	default:
		return nil, fmt.Errorf("Schema not found for catalog %s", src.CatalogName)
	}
}

func fitsFactory(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	fitsReader, err := fits_reader.NewFitsReader(src, fits_reader.WithBatchSize(cfg.BatchSize))
	if err != nil {
		return nil, err
	}
	return &fitsReader, nil
}
