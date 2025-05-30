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

package di

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	reader_factory "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/factory"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	partition_reader "github.com/dirodriguezm/xmatch/service/internal/preprocessor/reader"
	"github.com/dirodriguezm/xmatch/service/internal/preprocessor/reducer"
	partition_writer "github.com/dirodriguezm/xmatch/service/internal/preprocessor/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/golobby/container/v3"
)

func BuildPreprocessorContainer() container.Container {
	ctr := container.New()

	ctr.Singleton(func() *config.Config {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		return cfg
	})

	ctr.Singleton(func() *slog.LevelVar {
		levels := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"info":  slog.LevelInfo,
			"error": slog.LevelError,
			"warn":  slog.LevelWarn,
			"":      slog.LevelInfo,
		}
		var programLevel = new(slog.LevelVar)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
		slog.SetDefault(logger)
		programLevel.Set(levels[os.Getenv("LOG_LEVEL")])
		return programLevel
	})

	// Register Source
	ctr.Singleton(func(cfg *config.Config) *source.Source {
		src, err := source.NewSource(cfg.Preprocessor.Source)
		if err != nil {
			slog.Error("Could not register Source")
			panic(err)
		}
		return src
	})

	// Register Preprocessor Reader
	ctr.Singleton(func(src *source.Source, cfg *config.Config) reader.Reader {
		readerResults := make(chan reader.ReaderResult)
		r, err := reader_factory.ReaderFactory(src, []chan reader.ReaderResult{readerResults}, cfg.Preprocessor.Reader)
		if err != nil {
			panic(fmt.Errorf("Could not register preprocessor reader: %w", err))
		}
		return r
	})

	// Register partition writer if is configured
	ctr.Singleton(func(cfg *config.Config, src *source.Source, r reader.Reader) *partition_writer.PartitionWriter {
		switch strings.ToLower(cfg.Preprocessor.Source.CatalogName) {
		case "allwise":
			cfg.Preprocessor.PartitionWriter.Schema = config.AllwiseSchema
		case "vlass":
			cfg.Preprocessor.PartitionWriter.Schema = config.VlassSchema
		case "test":
			cfg.Preprocessor.PartitionWriter.Schema = config.TestSchema
		default:
			panic(fmt.Errorf("Can't register partition writer: unknown catalog name %s", cfg.Preprocessor.Source.CatalogName))
		}

		inputChan := make(chan writer.WriterInput[repository.InputSchema])
		doneChan := make(chan struct{})
		return partition_writer.New(cfg.Preprocessor.PartitionWriter, inputChan, doneChan)
	})

	// Register PartitionReader Workers
	dirChannel := make(chan string)
	reducerChannel := make(chan partition_reader.Records)

	ctr.Singleton(func(cfg *config.Config) []*partition_reader.Worker {
		workers := make([]*partition_reader.Worker, cfg.Preprocessor.PartitionReader.NumWorkers)
		for i := range cfg.Preprocessor.PartitionReader.NumWorkers {
			workers[i] = partition_reader.NewWorker(
				dirChannel,
				cfg.Preprocessor.PartitionWriter.Schema,
				reducerChannel,
			)
		}
		return workers
	})

	// Register PartitionReader
	ctr.Singleton(
		func(workers []*partition_reader.Worker, cfg *config.Config) *partition_reader.PartitionReader {
			return partition_reader.NewPartitionReader(
				dirChannel,
				workers,
				cfg.Preprocessor.PartitionWriter.BaseDir,
			)
		},
	)

	// Register Reducer Workers
	processedObjectsChannel := make(chan writer.WriterInput[repository.InputSchema])
	ctr.Singleton(func(cfg *config.Config) []*reducer.Worker {
		workers := make([]*reducer.Worker, cfg.Preprocessor.PartitionReader.NumWorkers)
		for i := range cfg.Preprocessor.PartitionReader.NumWorkers {
			workers[i] = reducer.NewWorker(
				reducerChannel,
				processedObjectsChannel,
				cfg.Preprocessor.PartitionWriter.Schema,
				cfg.Preprocessor.ReducerWriter.BatchSize,
			)
		}
		return workers
	})

	// Register Reducer Writer
	ctr.Singleton(func(cfg *config.Config) *parquet_writer.ParquetWriter[repository.InputSchema] {
		switch strings.ToLower(cfg.Preprocessor.Source.CatalogName) {
		case "vlass":
			cfg.Preprocessor.ReducerWriter.WriterConfig.Schema = config.VlassSchema
		case "default":
			panic("Catalog name not configured or unknown")
		}

		writer, err := parquet_writer.NewParquetWriter(processedObjectsChannel, make(chan struct{}), &cfg.Preprocessor.ReducerWriter.WriterConfig)
		if err != nil {
			panic(err)
		}
		return writer
	})

	// Register Reducer
	ctr.Singleton(func(workers []*reducer.Worker, writer *parquet_writer.ParquetWriter[repository.InputSchema]) *reducer.Reducer {
		return reducer.NewReducer(workers, writer)
	})

	return ctr
}
