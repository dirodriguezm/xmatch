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
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseConfig(t *testing.T) PipelineConfig {
	t.Helper()

	return PipelineConfig{
		Context: context.Background(),
		Config: config.Config{
			CatalogIndexer: config.CatalogIndexerConfig{
				Source: config.SourceConfig{
					CatalogName: "test",
					Nside:       5,
					Metadata:    false,
				},
				Reader: config.ReaderConfig{
					BatchSize: 10,
					Type:      "csv",
				},
				Indexer: config.IndexerConfig{
					OrderingScheme: "nested",
					Nside:          5,
				},
				IndexerWriter:  config.WriterConfig{Type: "sqlite"},
				MetadataWriter: config.WriterConfig{Type: "sqlite"},
				ChannelSize:    10,
			},
		},
		DB:      newTestDB(t),
		Source:  &source.Source{Sources: []string{"buffer:id\n1"}, CatalogName: "test", Nside: 5},
		Adapter: catalog.NewMockCatalogIndexAdapter(t),
		Store:   conesearch.NewMockMastercatStore(t),
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", "file:pipeline-test?mode=memory&cache=shared")
	require.NoError(t, err)
	return db
}

type noopReader struct{}

func (r *noopReader) Read() ([]any, error)      { return nil, io.EOF }
func (r *noopReader) ReadBatch() ([]any, error) { return nil, io.EOF }
func (r *noopReader) Close() error              { return nil }

func TestPipeline_WiringWithoutMetadata(t *testing.T) {
	cfg := baseConfig(t)

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
	cfg := baseConfig(t)
	cfg.Config.CatalogIndexer.Source.Metadata = true
	cfg.Adapter.(*catalog.MockCatalogIndexAdapter).
		EXPECT().
		BulkInsertFn().
		Return(func(context.Context, *sql.DB, []any) error { return nil })

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
	ctx := t.Context()
	var mcWriterStopped, mcIndexerStopped bool

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mcWriterStopped = true
	}, nil, ctx)
	mcIndexer := actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mcIndexerStopped = true
	}, []*actor.Actor{mcWriter}, ctx)

	pipeline := &Pipeline{mastercatWriter: mcWriter, mastercatIndexer: mcIndexer}
	pipeline.Stop()

	assert.True(t, mcIndexerStopped, "mastercat indexer should be stopped")
	assert.True(t, mcWriterStopped, "mastercat writer should be stopped")
}

func TestPipeline_Stop_WithMetadata(t *testing.T) {
	ctx := t.Context()
	var mcWriterStopped, mcIndexerStopped, mdWriterStopped, mdIndexerStopped bool

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mcWriterStopped = true
	}, nil, ctx)
	mdWriter := actor.New("md-writer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mdWriterStopped = true
	}, nil, ctx)
	mcIndexer := actor.New("mc-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mcIndexerStopped = true
	}, []*actor.Actor{mcWriter}, ctx)
	mdIndexer := actor.New("md-indexer", 10, func(*actor.Actor, actor.Message) {}, func(*actor.Actor) {
		mdIndexerStopped = true
	}, []*actor.Actor{mdWriter}, ctx)

	pipeline := &Pipeline{
		mastercatWriter:  mcWriter,
		metadataWriter:   mdWriter,
		mastercatIndexer: mcIndexer,
		metadataIndexer:  mdIndexer,
	}
	pipeline.Stop()

	assert.True(t, mcIndexerStopped, "mastercat indexer should be stopped")
	assert.True(t, mcWriterStopped, "mastercat writer should be stopped")
	assert.True(t, mdWriterStopped, "metadata writer should be stopped when metadata enabled")
	assert.True(t, mdIndexerStopped, "metadata indexer should be stopped when metadata enabled")
}

func TestPipeline_ErrorPropagation(t *testing.T) {
	ctx := t.Context()
	receivedErrors := make(chan error, 1)

	mcWriter := actor.New("mc-writer", 10, func(_ *actor.Actor, msg actor.Message) {
		if msg.Error != nil {
			receivedErrors <- msg.Error
		}
	}, nil, nil, ctx)
	mcIndexer := actor.New("mc-indexer", 10, func(a *actor.Actor, msg actor.Message) {
		a.Broadcast(msg)
	}, nil, []*actor.Actor{mcWriter}, ctx)

	mcWriter.Start()
	mcIndexer.Start()
	pipeline := &Pipeline{mastercatWriter: mcWriter, mastercatIndexer: mcIndexer}

	testErr := io.ErrUnexpectedEOF
	mcIndexer.Send(actor.Message{Error: testErr})

	errReceived := <-receivedErrors
	assert.Equal(t, testErr, errReceived)

	pipeline.Stop()
}

func TestPipeline_CloseSource(t *testing.T) {
	closed := false
	pipeline := &Pipeline{sourceReader: reader.SourceReader{Reader: &closeTrackingReader{closed: &closed}}}

	err := pipeline.CloseSource()
	require.NoError(t, err)
	assert.True(t, closed)
}

type closeTrackingReader struct {
	noopReader
	closed *bool
}

func (r *closeTrackingReader) Close() error {
	*r.closed = true
	return nil
}
