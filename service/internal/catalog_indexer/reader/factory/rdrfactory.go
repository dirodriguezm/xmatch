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

// Package rdrfactory provides a factory for creating readers
package rdrfactory

import (
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

func ReaderFactory(
	src *source.Source,
	cfg config.ReaderConfig,
) (reader.Reader, error) {
	if cfg.BatchSize <= 0 {
		return nil, fmt.Errorf("batch size must be greater than 0")
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
		return nil, fmt.Errorf("reader type not allowed")
	}
}

func parquetFactory(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	adapter, err := catalog.GetFactory(src.CatalogName)
	if err != nil {
		return nil, err
	}
	return adapter.NewParquetReader(src, cfg)
}

func fitsFactory(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	adapter, err := catalog.GetFactory(src.CatalogName)
	if err != nil {
		return nil, err
	}
	return adapter.NewFitsReader(src, cfg)
}
