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
	"io"
	"sync"
	"testing"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func restoreDefaults() {
	NewMastercatWriter = defaultMastercatWriter
	NewMetadataWriter = defaultMetadataWriter
	NewMastercatIndexer = defaultMastercatIndexer
	NewMetadataIndexer = defaultMetadataIndexer
	NewSourceReader = defaultSourceReader
}

func baseConfig() PipelineConfig {
	return PipelineConfig{
		Context: context.Background(),
		Config: config.Config{
			CatalogIndexer: config.CatalogIndexerConfig{
				Source: config.SourceConfig{
					CatalogName: "test",
					Nside:       32,
					Metadata:    false,
				},
				Reader: config.ReaderConfig{
					BatchSize: 10,
					Type:      "csv",
				},
				Indexer: config.IndexerConfig{
					OrderingScheme: "nested",
					Nside:          32,
				},
				IndexerWriter: config.WriterConfig{
					Type: "sqlite",
				},
				MetadataWriter: config.WriterConfig{
					Type: "sqlite",
				},
				ChannelSize: 10,
			},
		},
		DB:      &sql.DB{},
		Source:  &source.Source{CatalogName: "test", Nside: 32},
		Adapter: newMockAdapter(),
		Store:   newMockStore(),
	}
}

func newMockStore() *mockMastercatStore {
	m := &mockMastercatStore{}
	m.On("BulkInsertObject", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return m
}

func newMockAdapter() *mockCatalogAdapter {
	m := &mockCatalogAdapter{}
	m.On("ConvertToMastercat", mock.Anything, mock.Anything).Return(repository.Mastercat{}, (error)(nil))
	m.On("ConvertToMetadataFromRaw", mock.Anything).Return(nil, (error)(nil))
	m.On("BulkInsertFn").Return(func(context.Context, *sql.DB, []any) error { return nil })
	return m
}

type mockMastercatStore struct {
	mock.Mock
}

func (m *mockMastercatStore) FindObjects(ctx context.Context, pixels []int64) ([]repository.Mastercat, error) {
	args := m.Called(ctx, pixels)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.Mastercat), args.Error(1)
}

func (m *mockMastercatStore) InsertMastercat(ctx context.Context, arg repository.Mastercat) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockMastercatStore) GetAllObjects(ctx context.Context) ([]repository.Mastercat, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.Mastercat), args.Error(1)
}

func (m *mockMastercatStore) RemoveAllObjects(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockMastercatStore) BulkInsertObject(ctx context.Context, db *sql.DB, rows []any) error {
	args := m.Called(ctx, db, rows)
	return args.Error(0)
}

type mockCatalogAdapter struct {
	mock.Mock
}

func (m *mockCatalogAdapter) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockCatalogAdapter) NewRawRecord() any {
	args := m.Called()
	return args.Get(0)
}

func (m *mockCatalogAdapter) NewParquetWriter(cfg config.WriterConfig, ctx context.Context) (writer.Writer, error) {
	args := m.Called(cfg, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(writer.Writer), args.Error(1)
}

func (m *mockCatalogAdapter) NewParquetReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	args := m.Called(src, cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(reader.Reader), args.Error(1)
}

func (m *mockCatalogAdapter) NewFitsReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	args := m.Called(src, cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(reader.Reader), args.Error(1)
}

func (m *mockCatalogAdapter) BulkInsertFn() func(context.Context, *sql.DB, []any) error {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(func(context.Context, *sql.DB, []any) error)
}

func (m *mockCatalogAdapter) GetByID(ctx context.Context, id string) (any, error) {
	args := m.Called(ctx, id)
	return args.Get(0), args.Error(1)
}

func (m *mockCatalogAdapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	args := m.Called(ctx, ids)
	return args.Get(0), args.Error(1)
}

func (m *mockCatalogAdapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.MetadataWithCoordinates, error) {
	args := m.Called(ctx, pixels)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.MetadataWithCoordinates), args.Error(1)
}

func (m *mockCatalogAdapter) ConvertToMetadata(obj repository.MetadataWithCoordinates) repository.Metadata {
	args := m.Called(obj)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(repository.Metadata)
}

func (m *mockCatalogAdapter) ConvertToMastercat(raw any, mapper *healpix.HEALPixMapper) (repository.Mastercat, error) {
	args := m.Called(raw, mapper)
	return args.Get(0).(repository.Mastercat), args.Error(1)
}

func (m *mockCatalogAdapter) ConvertToMetadataFromRaw(raw any) (repository.Metadata, error) {
	args := m.Called(raw)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(repository.Metadata), args.Error(1)
}

// noopReader is a Reader implementation that returns EOF immediately.
type noopReader struct{}

func (r *noopReader) Read() ([]any, error)     { return nil, io.EOF }
func (r *noopReader) ReadBatch() ([]any, error) { return nil, io.EOF }
func (r *noopReader) Close() error              { return nil }

func TestPipeline_WiringWithoutMetadata(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	ctx := t.Context()

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMetadataWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ catalog.CatalogAdapter) (*actor.Actor, error) {
		return actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, nil, []*actor.Actor{writer}, ctx), nil
	}
	NewMetadataIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any) repository.Metadata) *actor.Actor {
		return actor.New("md-indexer", 10, func(*actor.Actor, actor.Message) {}, nil, []*actor.Actor{writer}, ctx)
	}
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, mcIndexer, mdIndexer *actor.Actor) (reader.SourceReader, error) {
		receivers := []*actor.Actor{mcIndexer}
		if mdIndexer != nil {
			receivers = append(receivers, mdIndexer)
		}
		return reader.SourceReader{Reader: &noopReader{}, BatchSize: 10, Receivers: receivers}, nil
	}

	pipeline, err := New(cfg)
	require.NoError(t, err)

	require.NotNil(t, pipeline.mastercatWriter)
	require.NotNil(t, pipeline.mastercatIndexer)

	require.Nil(t, pipeline.metadataWriter, "metadata writer should be nil when metadata disabled")
	require.Nil(t, pipeline.metadataIndexer, "metadata indexer should be nil when metadata disabled")

	require.Len(t, pipeline.sourceReader.Receivers, 1)
	assert.Same(t, pipeline.mastercatIndexer, pipeline.sourceReader.Receivers[0])

	recvs := pipeline.mastercatIndexer.Receivers()
	require.Len(t, recvs, 1)
	assert.Same(t, pipeline.mastercatWriter, recvs[0])

	pipeline.Stop()
}

func TestPipeline_WiringWithMetadata(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	cfg.Config.CatalogIndexer.Source.Metadata = true
	ctx := t.Context()

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMetadataWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ catalog.CatalogAdapter) (*actor.Actor, error) {
		return actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, nil, []*actor.Actor{writer}, ctx), nil
	}
	NewMetadataIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any) repository.Metadata) *actor.Actor {
		return actor.New("md-indexer", 10, func(*actor.Actor, actor.Message) {}, nil, []*actor.Actor{writer}, ctx)
	}
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, mcIndexer, mdIndexer *actor.Actor) (reader.SourceReader, error) {
		return reader.SourceReader{Reader: &noopReader{}, BatchSize: 10, Receivers: []*actor.Actor{mcIndexer, mdIndexer}}, nil
	}

	pipeline, err := New(cfg)
	require.NoError(t, err)

	require.NotNil(t, pipeline.mastercatWriter)
	require.NotNil(t, pipeline.mastercatIndexer)
	require.NotNil(t, pipeline.metadataWriter, "metadata writer should be non-nil when metadata enabled")
	require.NotNil(t, pipeline.metadataIndexer, "metadata indexer should be non-nil when metadata enabled")

	require.Len(t, pipeline.sourceReader.Receivers, 2)
	assert.Same(t, pipeline.mastercatIndexer, pipeline.sourceReader.Receivers[0])
	assert.Same(t, pipeline.metadataIndexer, pipeline.sourceReader.Receivers[1])

	recvs := pipeline.mastercatIndexer.Receivers()
	require.Len(t, recvs, 1)
	assert.Same(t, pipeline.mastercatWriter, recvs[0])

	recvs = pipeline.metadataIndexer.Receivers()
	require.Len(t, recvs, 1)
	assert.Same(t, pipeline.metadataWriter, recvs[0])

	pipeline.Stop()
}

func TestPipeline_Stop(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	ctx := t.Context()

	var mcWriterStopped, mcIndexerStopped, mdWriterStopped, mdIndexerStopped bool
	var mu sync.Mutex

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mu.Lock()
		mcWriterStopped = true
		mu.Unlock()
	}, nil, ctx)

	mdWriter := actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mu.Lock()
		mdWriterStopped = true
		mu.Unlock()
	}, nil, ctx)

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return mcWriter, nil
	}
	NewMetadataWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ catalog.CatalogAdapter) (*actor.Actor, error) {
		return mdWriter, nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
			mu.Lock()
			mcIndexerStopped = true
			mu.Unlock()
		}, []*actor.Actor{writer}, ctx), nil
	}
	NewMetadataIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any) repository.Metadata) *actor.Actor {
		return actor.New("md-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
			mu.Lock()
			mdIndexerStopped = true
			mu.Unlock()
		}, []*actor.Actor{writer}, ctx)
	}
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, mcIndexer, mdIndexer *actor.Actor) (reader.SourceReader, error) {
		return reader.SourceReader{Reader: &noopReader{}, BatchSize: 10, Receivers: []*actor.Actor{mcIndexer}}, nil
	}

	mcWriter.Start()
	mdWriter.Start()
	pipeline, err := New(cfg)
	require.NoError(t, err)

	pipeline.Stop()

	assert.True(t, mcIndexerStopped, "mastercat indexer should be stopped")
	assert.True(t, mcWriterStopped, "mastercat writer should be stopped")
	assert.False(t, mdWriterStopped, "metadata writer should not be stopped when metadata disabled")
	assert.False(t, mdIndexerStopped, "metadata indexer should not be stopped when metadata disabled")
}

func TestPipeline_Stop_WithMetadata(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	cfg.Config.CatalogIndexer.Source.Metadata = true
	ctx := t.Context()

	var mcWriterStopped, mcIndexerStopped, mdWriterStopped, mdIndexerStopped bool
	var mu sync.Mutex

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mu.Lock()
		mcWriterStopped = true
		mu.Unlock()
	}, nil, ctx)

	mdWriter := actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mu.Lock()
		mdWriterStopped = true
		mu.Unlock()
	}, nil, ctx)

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return mcWriter, nil
	}
	NewMetadataWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ catalog.CatalogAdapter) (*actor.Actor, error) {
		return mdWriter, nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
			mu.Lock()
			mcIndexerStopped = true
			mu.Unlock()
		}, []*actor.Actor{writer}, ctx), nil
	}
	NewMetadataIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any) repository.Metadata) *actor.Actor {
		return actor.New("md-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
			mu.Lock()
			mdIndexerStopped = true
			mu.Unlock()
		}, []*actor.Actor{writer}, ctx)
	}
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, mcIndexer, mdIndexer *actor.Actor) (reader.SourceReader, error) {
		return reader.SourceReader{Reader: &noopReader{}, BatchSize: 10, Receivers: []*actor.Actor{mcIndexer, mdIndexer}}, nil
	}

	mcWriter.Start()
	mdWriter.Start()
	pipeline, err := New(cfg)
	require.NoError(t, err)

	pipeline.Stop()

	assert.True(t, mcIndexerStopped, "mastercat indexer should be stopped")
	assert.True(t, mcWriterStopped, "mastercat writer should be stopped")
	assert.True(t, mdWriterStopped, "metadata writer should be stopped when metadata enabled")
	assert.True(t, mdIndexerStopped, "metadata indexer should be stopped when metadata enabled")
}

func TestPipeline_ErrorPropagation(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	ctx := t.Context()

	receivedErrors := make(chan error, 10)

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return actor.New("mc-writer", 10, func(_ *actor.Actor, msg actor.Message) {
			if msg.Error != nil {
				receivedErrors <- msg.Error
			}
		}, nil, nil, ctx), nil
	}
	NewMetadataWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ catalog.CatalogAdapter) (*actor.Actor, error) {
		return actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(a *actor.Actor, msg actor.Message) {
			a.Broadcast(msg)
		}, nil, []*actor.Actor{writer}, ctx), nil
	}
	NewMetadataIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any) repository.Metadata) *actor.Actor {
		return actor.New("md-indexer", 10, func(a *actor.Actor, msg actor.Message) {
			a.Broadcast(msg)
		}, nil, []*actor.Actor{writer}, ctx)
	}
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, mcIndexer, mdIndexer *actor.Actor) (reader.SourceReader, error) {
		return reader.SourceReader{Reader: &noopReader{}, BatchSize: 10, Receivers: []*actor.Actor{mcIndexer}}, nil
	}

	pipeline, err := New(cfg)
	require.NoError(t, err)

	testErr := io.ErrUnexpectedEOF
	mcIndexer := pipeline.sourceReader.Receivers[0]
	mcIndexer.Send(actor.Message{Error: testErr})

	errReceived := <-receivedErrors
	assert.Equal(t, testErr, errReceived)

	pipeline.Stop()
}

func TestPipeline_CloseSource(t *testing.T) {
	defer restoreDefaults()

	cfg := baseConfig()
	ctx := t.Context()

	closed := false
	NewSourceReader = func(_ *source.Source, _ config.ReaderConfig, _ config.SourceConfig, _, _ *actor.Actor) (reader.SourceReader, error) {
		return reader.SourceReader{
			Reader: &closeTrackingReader{closed: &closed},
		}, nil
	}

	NewMastercatWriter = func(_ context.Context, _ config.CatalogIndexerConfig, _ *sql.DB, _ conesearch.MastercatStore) (*actor.Actor, error) {
		return actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, nil, nil, ctx), nil
	}
	NewMastercatIndexer = func(_ config.CatalogIndexerConfig, writer *actor.Actor, _ context.Context, _ func(any, *healpix.HEALPixMapper) repository.Mastercat) (*actor.Actor, error) {
		return actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, nil, []*actor.Actor{writer}, ctx), nil
	}

	pipeline, err := New(cfg)
	require.NoError(t, err)

	err = pipeline.CloseSource()
	require.NoError(t, err)
	assert.True(t, closed)

	pipeline.Stop()
}

type closeTrackingReader struct {
	noopReader
	closed *bool
}

func (r *closeTrackingReader) Close() error {
	*r.closed = true
	return nil
}
