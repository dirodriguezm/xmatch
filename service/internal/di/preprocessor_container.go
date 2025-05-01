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
	"github.com/dirodriguezm/xmatch/service/internal/config"
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
		case "test":
			cfg.Preprocessor.PartitionWriter.Schema = config.TestSchema
		default:
			panic(fmt.Errorf("Can't register partition writer: unknown catalog name %s", cfg.CatalogIndexer.Source.CatalogName))
		}

		inputChan := make(chan writer.WriterInput[repository.InputSchema])
		doneChan := make(chan struct{})
		return partition_writer.New(cfg.Preprocessor.PartitionWriter, inputChan, doneChan)
	})

	return ctr
}
