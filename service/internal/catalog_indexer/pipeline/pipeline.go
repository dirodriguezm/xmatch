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

// Package pipeline implements the healpix indexer pipeline.
package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	fits_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/fits"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// PipelineConfig holds the catalog pipeline components composed at the
// boundary. The metadata components are nil together when the source has no
// metadata branch.
type PipelineConfig[Raw, Md any] struct {
	MastercatWriter  *actor.Actor[repository.Mastercat]
	MastercatIndexer *actor.Actor[Raw]
	MetadataWriter   *actor.Actor[Md]
	MetadataIndexer  *actor.Actor[Raw]
	SourceReader     reader.SourceReader[Raw]
}

type Pipeline[Raw, Md any] struct {
	sourceReader     reader.SourceReader[Raw]
	mastercatWriter  *actor.Actor[repository.Mastercat]
	metadataWriter   *actor.Actor[Md]
	mastercatIndexer *actor.Actor[Raw]
	metadataIndexer  *actor.Actor[Raw]
}

func New[Raw, Md any](cfg PipelineConfig[Raw, Md]) (*Pipeline[Raw, Md], error) {
	if cfg.MastercatWriter == nil {
		return nil, fmt.Errorf("building pipeline: mastercat writer is required")
	}
	if cfg.MastercatIndexer == nil {
		return nil, fmt.Errorf("building pipeline: mastercat indexer is required")
	}
	if (cfg.MetadataWriter == nil) != (cfg.MetadataIndexer == nil) {
		return nil, fmt.Errorf("building pipeline: metadata writer and indexer must be set together")
	}

	cfg.MastercatWriter.Start()
	if cfg.MetadataWriter != nil {
		cfg.MetadataWriter.Start()
	}

	cfg.MastercatIndexer.Start()
	if cfg.MetadataIndexer != nil {
		cfg.MetadataIndexer.Start()
	}

	return &Pipeline[Raw, Md]{
		sourceReader:     cfg.SourceReader,
		mastercatWriter:  cfg.MastercatWriter,
		metadataWriter:   cfg.MetadataWriter,
		mastercatIndexer: cfg.MastercatIndexer,
		metadataIndexer:  cfg.MetadataIndexer,
	}, nil
}

func (p *Pipeline[Raw, Md]) Run() {
	p.sourceReader.Read()
}

func (p *Pipeline[Raw, Md]) Stop() {
	p.mastercatIndexer.Stop()
	p.mastercatWriter.Stop()
	if p.metadataIndexer != nil {
		p.metadataIndexer.Stop()
		p.metadataWriter.Stop()
	}
}

func (p *Pipeline[Raw, Md]) CloseSource() error {
	return p.sourceReader.Close()
}

// DefaultMastercatWriter builds the mastercat sqlite writer actor.
func DefaultMastercatWriter(
	ctx context.Context,
	cfg config.CatalogIndexerConfig,
	store *repository.Queries,
) (*actor.Actor[repository.Mastercat], error) {
	w := sqlite_writer.New[repository.Mastercat](ctx, store)
	return actor.New("mastercat writer", cfg.ChannelSize, w.Write, w.Stop, ctx), nil
}

// DefaultMetadataWriter builds the metadata sqlite writer actor for a model
// that knows how to insert itself.
func DefaultMetadataWriter[Md repository.Insertable](
	ctx context.Context,
	cfg config.CatalogIndexerConfig,
	store *repository.Queries,
) (*actor.Actor[Md], error) {
	w := sqlite_writer.New[Md](ctx, store)
	return actor.New("metadata writer", cfg.ChannelSize, w.Write, w.Stop, ctx), nil
}

// DefaultMastercatIndexer builds the mastercat indexer actor from the catalog
// row type.
func DefaultMastercatIndexer[Raw mastercat_indexer.Source](
	cfg config.CatalogIndexerConfig,
	writer *actor.Actor[repository.Mastercat],
	ctx context.Context,
) (*actor.Actor[Raw], error) {
	ind, err := mastercat_indexer.New[Raw](cfg.Indexer, cfg.Source.CatalogName, writer)
	if err != nil {
		return nil, err
	}
	return actor.New("mastercat indexer", cfg.ChannelSize, ind.Index, func(*actor.Actor[Raw]) {
		writer.Stop()
	}, ctx), nil
}

// DefaultMetadataIndexer builds the metadata indexer actor from the catalog
// row type's own conversion.
func DefaultMetadataIndexer[Md any, Raw metadata.Converter[Md]](
	cfg config.CatalogIndexerConfig,
	writer *actor.Actor[Md],
	ctx context.Context,
) *actor.Actor[Raw] {
	ind := metadata.New[Md, Raw](writer)
	return actor.New("metadata indexer", cfg.ChannelSize, ind.Index, func(*actor.Actor[Raw]) {
		writer.Stop()
	}, ctx)
}

// DefaultSourceReader builds the source reader for the configured source
// format. The metadata indexer is only appended when the source has metadata.
func DefaultSourceReader[Raw any](
	src *source.Source,
	readerCfg config.ReaderConfig,
	srcCfg config.SourceConfig,
	mastercatIndexer *actor.Actor[Raw],
	metadataIndexer *actor.Actor[Raw],
) (reader.SourceReader[Raw], error) {
	if readerCfg.BatchSize <= 0 {
		return reader.SourceReader[Raw]{}, fmt.Errorf("batch size must be greater than 0")
	}

	var r reader.Reader[Raw]
	var err error
	switch strings.ToLower(readerCfg.Type) {
	case "csv":
		r, err = csv_reader.NewCsvReader(
			src,
			csv_reader.WithHeader[Raw](readerCfg.Header),
			csv_reader.WithFirstLineHeader[Raw](readerCfg.FirstLineHeader),
			csv_reader.WithComment[Raw](readerCfg.Comment),
			csv_reader.WithCsvBatchSize[Raw](readerCfg.BatchSize),
		)
	case "parquet":
		r, err = parquet_reader.NewParquetReader(
			src,
			parquet_reader.WithParquetBatchSize[Raw](readerCfg.BatchSize),
		)
	case "fits":
		r, err = fits_reader.NewFitsReader(
			src,
			fits_reader.WithBatchSize[Raw](readerCfg.BatchSize),
		)
	default:
		return reader.SourceReader[Raw]{}, fmt.Errorf("reader type not allowed")
	}
	if err != nil {
		return reader.SourceReader[Raw]{}, err
	}

	sourceReader := reader.SourceReader[Raw]{
		Reader:    r,
		BatchSize: readerCfg.BatchSize,
		Receivers: []*actor.Actor[Raw]{mastercatIndexer},
	}
	if srcCfg.Metadata {
		sourceReader.Receivers = append(sourceReader.Receivers, metadataIndexer)
	}
	return sourceReader, nil
}
