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

package pipeline

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	mastercat_indexer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/mastercat"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer/metadata"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	rdrfactory "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/factory"
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
	Adapter catalog.CatalogAdapter
	Store   conesearch.MastercatStore
}

type Pipeline struct {
	sourceReader     reader.SourceReader
	mastercatWriter  *actor.Actor
	metadataWriter   *actor.Actor
	mastercatIndexer *actor.Actor
	metadataIndexer  *actor.Actor
}

// Factory function types expose constructor injection points for tests.
type (
	MastercatWriterFn func(context.Context, config.CatalogIndexerConfig, *sql.DB, conesearch.MastercatStore) (*actor.Actor, error)
	MetadataWriterFn  func(context.Context, config.CatalogIndexerConfig, *sql.DB, catalog.CatalogAdapter) (*actor.Actor, error)
	MastercatIndexerFn func(config.CatalogIndexerConfig, *actor.Actor, context.Context, func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error)
	MetadataIndexerFn  func(config.CatalogIndexerConfig, *actor.Actor, context.Context, func(any) repository.Metadata) *actor.Actor
	SourceReaderFn     func(*source.Source, config.ReaderConfig, config.SourceConfig, *actor.Actor, *actor.Actor) (reader.SourceReader, error)
)

// Overrideable factory functions for test injection.
var (
	NewMastercatWriter  MastercatWriterFn  = defaultMastercatWriter
	NewMetadataWriter   MetadataWriterFn   = defaultMetadataWriter
	NewMastercatIndexer MastercatIndexerFn = defaultMastercatIndexer
	NewMetadataIndexer  MetadataIndexerFn  = defaultMetadataIndexer
	NewSourceReader     SourceReaderFn     = defaultSourceReader
)

func New(cfg PipelineConfig) (*Pipeline, error) {
	ciCfg := cfg.Config.CatalogIndexer

	mWriter, err := NewMastercatWriter(cfg.Context, ciCfg, cfg.DB, cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("building mastercat writer: %w", err)
	}
	mWriter.Start()

	var mdWriter *actor.Actor
	if ciCfg.Source.Metadata {
		mdWriter, err = NewMetadataWriter(cfg.Context, ciCfg, cfg.DB, cfg.Adapter)
		if err != nil {
			return nil, fmt.Errorf("building metadata writer: %w", err)
		}
		mdWriter.Start()
	}

	fillMastercat := func(raw any, mapper *healpix.HEALPixMapper) repository.Mastercat {
		mc, _ := cfg.Adapter.ConvertToMastercat(raw, mapper)
		return mc
	}
	mIndexer, err := NewMastercatIndexer(ciCfg, mWriter, cfg.Context, fillMastercat)
	if err != nil {
		return nil, fmt.Errorf("building mastercat indexer: %w", err)
	}
	mIndexer.Start()

	var mdIndexer *actor.Actor
	if ciCfg.Source.Metadata {
		fillMetadata := func(raw any) repository.Metadata {
			md, _ := cfg.Adapter.ConvertToMetadataFromRaw(raw)
			return md
		}
		mdIndexer = NewMetadataIndexer(ciCfg, mdWriter, cfg.Context, fillMetadata)
		mdIndexer.Start()
	}

	srcReader, err := NewSourceReader(cfg.Source, ciCfg.Reader, ciCfg.Source, mIndexer, mdIndexer)
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
	adapter catalog.CatalogAdapter,
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
	fillMetadata func(any) repository.Metadata,
) *actor.Actor {
	ind := metadata.New(fillMetadata)
	return actor.New("metadata indexer", cfg.ChannelSize, ind.Index, nil, []*actor.Actor{writer}, ctx)
}

func defaultSourceReader(
	src *source.Source,
	readerCfg config.ReaderConfig,
	srcCfg config.SourceConfig,
	mastercatIndexer *actor.Actor,
	metadataIndexer *actor.Actor,
) (reader.SourceReader, error) {
	r, err := rdrfactory.ReaderFactory(src, readerCfg)
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
