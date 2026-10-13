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
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/actor"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRow struct {
	ID  string
	RA  float64
	Dec float64
}

func (r testRow) GetCoordinates() (float64, float64) { return r.RA, r.Dec }

func (r testRow) MastercatID() string { return r.ID }

func (r testRow) ConvertToMetadata() (testMetadata, error) {
	return testMetadata{ID: r.ID}, nil
}

type testMetadata struct {
	ID string
}

func (m testMetadata) Insert(context.Context, *repository.Queries) error { return nil }

type fakeReader[T any] struct {
	batches [][]T
	err     error
	current int
}

func (r *fakeReader[T]) Read() ([]T, error) {
	rows := make([]T, 0)
	for {
		batch, err := r.ReadBatch()
		rows = append(rows, batch...)
		if err != nil {
			return rows, err
		}
	}
}

func (r *fakeReader[T]) ReadBatch() ([]T, error) {
	if r.current < len(r.batches) {
		batch := r.batches[r.current]
		r.current++
		return batch, nil
	}
	if r.err != nil {
		return nil, r.err
	}
	return nil, io.EOF
}

func (r *fakeReader[T]) Close() error { return nil }

type noopReader[T any] struct{}

func (r *noopReader[T]) Read() ([]T, error)      { return nil, io.EOF }
func (r *noopReader[T]) ReadBatch() ([]T, error) { return nil, io.EOF }
func (r *noopReader[T]) Close() error            { return nil }

type recorder[T any] struct {
	mu   sync.Mutex
	rows [][]T
	errs []error
}

func (r *recorder[T]) handler(_ *actor.Actor[T], msg actor.Message[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if msg.Error != nil {
		r.errs = append(r.errs, msg.Error)
	}
	r.rows = append(r.rows, msg.Rows)
}

func (r *recorder[T]) flat() []T {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows := make([]T, 0)
	for _, batch := range r.rows {
		rows = append(rows, batch...)
	}
	return rows
}

func (r *recorder[T]) errors() []error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]error(nil), r.errs...)
}

func fakeIndexer[Raw, Out any](
	writer *actor.Actor[Out],
	convert func(Raw) Out,
	ctx context.Context,
) *actor.Actor[Raw] {
	return actor.New("fake indexer", 10, func(_ *actor.Actor[Raw], msg actor.Message[Raw]) {
		if msg.Error != nil {
			writer.Send(actor.Message[Out]{Error: msg.Error})
			return
		}
		rows := make([]Out, len(msg.Rows))
		for i, row := range msg.Rows {
			rows[i] = convert(row)
		}
		writer.Send(actor.Message[Out]{Rows: rows})
	}, nil, ctx)
}

func baseConfig(metadata bool) config.CatalogIndexerConfig {
	return config.CatalogIndexerConfig{
		Source: config.SourceConfig{
			CatalogName: "Gaia",
			Nside:       5,
			Metadata:    metadata,
		},
		Reader: config.ReaderConfig{
			BatchSize: 10,
			Type:      "csv",
		},
		Indexer: config.IndexerConfig{
			OrderingScheme: "nested",
			Nside:          5,
		},
		ChannelSize: 10,
	}
}

func TestPipeline_RunBroadcastsToBothIndexers(t *testing.T) {
	ctx := t.Context()

	mcRecorder := &recorder[repository.Mastercat]{}
	mcWriter := actor.New("mc-writer", 10, mcRecorder.handler, nil, ctx)
	mcIndexer := fakeIndexer(mcWriter, func(row testRow) repository.Mastercat {
		return repository.Mastercat{ID: row.ID, Ra: row.RA, Dec: row.Dec}
	}, ctx)

	mdRecorder := &recorder[testMetadata]{}
	mdWriter := actor.New("md-writer", 10, mdRecorder.handler, nil, ctx)
	mdIndexer := fakeIndexer(mdWriter, func(row testRow) testMetadata {
		return testMetadata{ID: row.ID}
	}, ctx)

	p := &Pipeline[testRow, testMetadata]{
		sourceReader: reader.SourceReader[testRow]{
			Reader: &fakeReader[testRow]{batches: [][]testRow{
				{{ID: "1", RA: 1, Dec: 1}, {ID: "2", RA: 2, Dec: 2}},
				{{ID: "3", RA: 3, Dec: 3}},
			}},
			BatchSize: 2,
			Receivers: []*actor.Actor[testRow]{mcIndexer, mdIndexer},
		},
		mastercatWriter:  mcWriter,
		metadataWriter:   mdWriter,
		mastercatIndexer: mcIndexer,
		metadataIndexer:  mdIndexer,
	}

	mcWriter.Start()
	mdWriter.Start()
	mcIndexer.Start()
	mdIndexer.Start()

	p.Run()
	p.Stop()

	require.Len(t, mcRecorder.flat(), 3)
	assert.Equal(t, "1", mcRecorder.flat()[0].ID)
	assert.Equal(t, "3", mcRecorder.flat()[2].ID)

	require.Len(t, mdRecorder.flat(), 3)
	assert.Equal(t, "1", mdRecorder.flat()[0].ID)
	assert.Equal(t, "3", mdRecorder.flat()[2].ID)
}

func TestPipeline_RunForwardsReadError(t *testing.T) {
	ctx := t.Context()
	readErr := errors.New("read failed")

	mcRecorder := &recorder[repository.Mastercat]{}
	mcWriter := actor.New("mc-writer", 10, mcRecorder.handler, nil, ctx)
	mcIndexer := fakeIndexer(mcWriter, func(row testRow) repository.Mastercat {
		return repository.Mastercat{ID: row.ID}
	}, ctx)

	p := &Pipeline[testRow, testMetadata]{
		sourceReader: reader.SourceReader[testRow]{
			Reader:    &fakeReader[testRow]{batches: [][]testRow{{{ID: "1"}}}, err: readErr},
			BatchSize: 1,
			Receivers: []*actor.Actor[testRow]{mcIndexer},
		},
		mastercatWriter:  mcWriter,
		mastercatIndexer: mcIndexer,
	}

	mcWriter.Start()
	mcIndexer.Start()

	p.Run()
	p.Stop()

	require.Len(t, mcRecorder.flat(), 1)
	require.Len(t, mcRecorder.errors(), 1)
	assert.Equal(t, readErr, mcRecorder.errors()[0])
}

func TestPipeline_StopWithoutMetadata(t *testing.T) {
	ctx := t.Context()
	var order []string

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor[repository.Mastercat], actor.Message[repository.Mastercat]) {}, func(*actor.Actor[repository.Mastercat]) {
		order = append(order, "mc-writer")
	}, ctx)
	mcIndexer := actor.New("mc-indexer", 10, func(*actor.Actor[testRow], actor.Message[testRow]) {}, func(*actor.Actor[testRow]) {
		order = append(order, "mc-indexer")
		mcWriter.Stop()
	}, ctx)

	p := &Pipeline[testRow, testMetadata]{mastercatWriter: mcWriter, mastercatIndexer: mcIndexer}
	p.Stop()

	assert.Equal(t, []string{"mc-indexer", "mc-writer"}, order)
}

func TestPipeline_StopWithMetadata(t *testing.T) {
	ctx := t.Context()
	var order []string

	mcWriter := actor.New("mc-writer", 10, func(*actor.Actor[repository.Mastercat], actor.Message[repository.Mastercat]) {}, func(*actor.Actor[repository.Mastercat]) {
		order = append(order, "mc-writer")
	}, ctx)
	mdWriter := actor.New("md-writer", 10, func(*actor.Actor[testMetadata], actor.Message[testMetadata]) {}, func(*actor.Actor[testMetadata]) {
		order = append(order, "md-writer")
	}, ctx)
	mcIndexer := actor.New("mc-indexer", 10, func(*actor.Actor[testRow], actor.Message[testRow]) {}, func(*actor.Actor[testRow]) {
		order = append(order, "mc-indexer")
		mcWriter.Stop()
	}, ctx)
	mdIndexer := actor.New("md-indexer", 10, func(*actor.Actor[testRow], actor.Message[testRow]) {}, func(*actor.Actor[testRow]) {
		order = append(order, "md-indexer")
		mdWriter.Stop()
	}, ctx)

	p := &Pipeline[testRow, testMetadata]{
		mastercatWriter:  mcWriter,
		metadataWriter:   mdWriter,
		mastercatIndexer: mcIndexer,
		metadataIndexer:  mdIndexer,
	}
	p.Stop()

	assert.Equal(t, []string{"mc-indexer", "mc-writer", "md-indexer", "md-writer"}, order)
}

func TestNew_WithoutMetadata(t *testing.T) {
	ctx := t.Context()
	cfg := baseConfig(false)
	store := repository.New(nil)

	mcWriter, err := DefaultMastercatWriter(ctx, cfg, store)
	require.NoError(t, err)
	mcIndexer, err := DefaultMastercatIndexer[testRow](cfg, mcWriter, ctx)
	require.NoError(t, err)

	src := &source.Source{Sources: []string{"buffer:id"}, CatalogName: "Gaia", Nside: 5}
	srcReader, err := DefaultSourceReader(src, cfg.Reader, cfg.Source, mcIndexer, nil)
	require.NoError(t, err)

	p, err := New(PipelineConfig[testRow, testMetadata]{
		MastercatWriter:  mcWriter,
		MastercatIndexer: mcIndexer,
		SourceReader:     srcReader,
	})
	require.NoError(t, err)
	defer p.Stop()

	require.Nil(t, p.metadataWriter, "metadata writer should be nil when metadata disabled")
	require.Nil(t, p.metadataIndexer, "metadata indexer should be nil when metadata disabled")
	require.Len(t, p.sourceReader.Receivers, 1)
	assert.Same(t, mcIndexer, p.sourceReader.Receivers[0])
}

func TestNew_WithMetadata(t *testing.T) {
	ctx := t.Context()
	cfg := baseConfig(true)
	store := repository.New(nil)

	mcWriter, err := DefaultMastercatWriter(ctx, cfg, store)
	require.NoError(t, err)
	mcIndexer, err := DefaultMastercatIndexer[testRow](cfg, mcWriter, ctx)
	require.NoError(t, err)

	mdWriter, err := DefaultMetadataWriter[testMetadata](ctx, cfg, store)
	require.NoError(t, err)
	mdIndexer := DefaultMetadataIndexer[testMetadata, testRow](cfg, mdWriter, ctx)

	src := &source.Source{Sources: []string{"buffer:id"}, CatalogName: "Gaia", Nside: 5}
	srcReader, err := DefaultSourceReader(src, cfg.Reader, cfg.Source, mcIndexer, mdIndexer)
	require.NoError(t, err)

	p, err := New(PipelineConfig[testRow, testMetadata]{
		MastercatWriter:  mcWriter,
		MastercatIndexer: mcIndexer,
		MetadataWriter:   mdWriter,
		MetadataIndexer:  mdIndexer,
		SourceReader:     srcReader,
	})
	require.NoError(t, err)
	defer p.Stop()

	require.NotNil(t, p.metadataWriter, "metadata writer should be non-nil when metadata enabled")
	require.NotNil(t, p.metadataIndexer, "metadata indexer should be non-nil when metadata enabled")
	require.Len(t, p.sourceReader.Receivers, 2)
	assert.Same(t, mcIndexer, p.sourceReader.Receivers[0])
	assert.Same(t, mdIndexer, p.sourceReader.Receivers[1])
}

func TestNew_RejectsHalfMetadataBranch(t *testing.T) {
	ctx := t.Context()
	cfg := baseConfig(true)
	store := repository.New(nil)

	mcWriter, err := DefaultMastercatWriter(ctx, cfg, store)
	require.NoError(t, err)
	mcIndexer, err := DefaultMastercatIndexer[testRow](cfg, mcWriter, ctx)
	require.NoError(t, err)

	_, err = New(PipelineConfig[testRow, testMetadata]{
		MastercatWriter:  mcWriter,
		MastercatIndexer: mcIndexer,
		MetadataIndexer:  mcIndexer,
	})
	require.Error(t, err)
}

func TestPipeline_CloseSource(t *testing.T) {
	closed := false
	pipeline := &Pipeline[int, int]{sourceReader: reader.SourceReader[int]{Reader: &closeTrackingReader[int]{closed: &closed}}}

	err := pipeline.CloseSource()
	require.NoError(t, err)
	assert.True(t, closed)
}

type closeTrackingReader[T any] struct {
	noopReader[T]
	closed *bool
}

func (r *closeTrackingReader[T]) Close() error {
	*r.closed = true
	return nil
}
