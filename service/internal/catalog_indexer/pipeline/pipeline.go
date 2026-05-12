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
	"database/sql"
	"fmt"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	csv_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/csv"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	sqlite_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/sqlite"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type PipelineConfig struct {
	Context context.Context
	Config  config.Config
	DB      *sql.DB
	Source  *source.Source
	Adapter catalog.CatalogIndexAdapter
	Store   conesearch.MastercatStore
}

type Pipeline struct {
	sourceReader     reader.SourceReader
	mastercatWriter  *actor.Actor
	metadataWriter   *actor.Actor
	mastercatIndexer *actor.Actor
	metadataIndexer  *actor.Actor
}

func New(cfg PipelineConfig) (*Pipeline, error) {
	ciCfg := cfg.Config.CatalogIndexer

	mWriter, err := defaultMastercatWriter(cfg.Context, ciCfg, cfg.DB, cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("building mastercat writer: %w", err)
	}
	mWriter.Start()

	var mdWriter *actor.Actor
	if ciCfg.Source.Metadata {
		mdWriter, err = defaultMetadataWriter(cfg.Context, ciCfg, cfg.DB, cfg.Adapter)
		if err != nil {
			return nil, fmt.Errorf("building metadata writer: %w", err)
		}
		mdWriter.Start()
	}

	fillMastercat := func(raw any, mapper *healpix.HEALPixMapper) repository.Mastercat {
		mc, _ := cfg.Adapter.ConvertToMastercat(raw, mapper)
		return mc
	}
	mIndexer, err := defaultMastercatIndexer(ciCfg, mWriter, cfg.Context, fillMastercat)
	if err != nil {
		return nil, fmt.Errorf("building mastercat indexer: %w", err)
	}
	mIndexer.Start()

	var mdIndexer *actor.Actor
	if ciCfg.Source.Metadata {
		fillMetadata := func(raw any) any {
			md, _ := cfg.Adapter.ConvertToMetadataFromRaw(raw)
			return md
		}
		mdIndexer = defaultMetadataIndexer(ciCfg, mdWriter, cfg.Context, fillMetadata)
		mdIndexer.Start()
	}

	srcReader, err := defaultSourceReader(cfg.Source, cfg.Adapter, ciCfg.Reader, ciCfg.Source, mIndexer, mdIndexer)
	if err != nil {
		return nil, fmt.Errorf("building reader: %w", err)
	}

	return &Pipeline{
		sourceReader:     srcReader,
		mastercatWriter:  mWriter,
		metadataWriter:   mdWriter,
		mastercatIndexer: mIndexer,
		metadataIndexer:  mdIndexer,
	}, nil
}

func (p *Pipeline) Run() {
	p.sourceReader.Read()
}

func (p *Pipeline) Stop() {
	p.mastercatIndexer.Stop()
	p.mastercatWriter.Stop()
	if p.metadataIndexer != nil {
		p.metadataIndexer.Stop()
		p.metadataWriter.Stop()
	}
}

func (p *Pipeline) CloseSource() error {
	return p.sourceReader.Close()
}

func defaultMastercatWriter(
	ctx context.Context,
	cfg config.CatalogIndexerConfig,
	db *sql.DB,
	store conesearch.MastercatStore,
) (*actor.Actor, error) {
	switch cfg.IndexerWriter.Type {
	case "parquet":
		w, err := parquet_writer.New[repository.Mastercat](cfg.IndexerWriter, ctx)
		if err != nil {
			return nil, err
		}
		return actor.New("mastercat writer", cfg.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	case "sqlite":
		w := sqlite_writer.New(db, ctx, store.BulkInsertObject)
		return actor.New("mastercat writer", cfg.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	default:
		return nil, fmt.Errorf("writer type not allowed")
	}
}

func defaultMetadataWriter(
	ctx context.Context,
	cfg config.CatalogIndexerConfig,
	db *sql.DB,
	adapter catalog.CatalogIndexAdapter,
) (*actor.Actor, error) {
	switch cfg.MetadataWriter.Type {
	case "parquet":
		w, err := adapter.NewParquetWriter(cfg.MetadataWriter, ctx)
		if err != nil {
			return nil, err
		}
		return actor.New("metadata writer", cfg.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	case "sqlite":
		w := sqlite_writer.New(db, ctx, adapter.BulkInsertFn())
		return actor.New("metadata writer", cfg.ChannelSize, w.Write, w.Stop, nil, ctx), nil
	default:
		return nil, fmt.Errorf("unknown Metadata Writer Type: %s", cfg.MetadataWriter.Type)
	}
}

func defaultMastercatIndexer(
	cfg config.CatalogIndexerConfig,
	writer *actor.Actor,
	ctx context.Context,
	fillMastercat func(any, *healpix.HEALPixMapper) repository.Mastercat,
) (*actor.Actor, error) {
	ind, err := mastercat_indexer.New(cfg.Indexer, fillMastercat)
	if err != nil {
		return nil, err
	}
	return actor.New("mastercat indexer", cfg.ChannelSize, ind.Index, nil, []*actor.Actor{writer}, ctx), nil
}

func defaultMetadataIndexer(
	cfg config.CatalogIndexerConfig,
	writer *actor.Actor,
	ctx context.Context,
	fillMetadata func(any) any,
) *actor.Actor {
	ind := metadata.New(fillMetadata)
	return actor.New("metadata indexer", cfg.ChannelSize, ind.Index, nil, []*actor.Actor{writer}, ctx)
}

func defaultSourceReader(
	src *source.Source,
	adapter catalog.CatalogIndexAdapter,
	readerCfg config.ReaderConfig,
	srcCfg config.SourceConfig,
	mastercatIndexer *actor.Actor,
	metadataIndexer *actor.Actor,
) (reader.SourceReader, error) {
	if readerCfg.BatchSize <= 0 {
		return reader.SourceReader{}, fmt.Errorf("batch size must be greater than 0")
	}

	var r reader.Reader
	var err error
	switch strings.ToLower(readerCfg.Type) {
	case "csv":
		r, err = csv_reader.NewCsvReader(
			src,
			adapter,
			csv_reader.WithHeader(readerCfg.Header),
			csv_reader.WithFirstLineHeader(readerCfg.FirstLineHeader),
			csv_reader.WithComment(readerCfg.Comment),
			csv_reader.WithCsvBatchSize(readerCfg.BatchSize),
		)
	case "parquet":
		r, err = adapter.NewParquetReader(src, readerCfg)
	case "fits":
		r, err = adapter.NewFitsReader(src, readerCfg)
	default:
		return reader.SourceReader{}, fmt.Errorf("reader type not allowed")
	}
	if err != nil {
		return reader.SourceReader{}, err
	}

	sourceReader := reader.SourceReader{
		Reader:    r,
		BatchSize: readerCfg.BatchSize,
		Receivers: []*actor.Actor{mastercatIndexer},
	}
	if srcCfg.Metadata {
		sourceReader.Receivers = append(sourceReader.Receivers, metadataIndexer)
	}
	return sourceReader, nil
}
