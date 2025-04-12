package csv_reader

import (
	"io"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	r := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
	csvReader, err := NewCsvReader(&source, outputs)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3, 3)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadWithHeader(t *testing.T) {
	csv := `o1,1,1
o2,2,2
o3,3,3
`

	r := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}
	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)

	csvReader, err := NewCsvReader(&source, outputs, WithHeader([]string{"oid", "ra", "dec"}))
	if err != nil {
		t.Fatal(err)
	}

	rows, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3, 3)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadWithHeader_Error(t *testing.T) {
	csv := ""

	r := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}
	outputs := []chan reader.ReaderResult{make(chan reader.ReaderResult)}

	csvReader, err := NewCsvReader(&source, outputs)
	require.NoError(t, err)

	rows, err := csvReader.Read()
	require.Error(t, err)
	require.Nil(t, rows)
}

func TestReadBatch(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
o4,4,4
`
	r := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
	csvReader, err := NewCsvReader(
		&source,
		outputs,
		WithCsvBatchSize(2),
		WithFirstLineHeader(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	var rows []repository.InputSchema

	for {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			t.Log("EOF")
			break
		}
		require.Len(t, batch, csvReader.BatchSize)
		for _, row := range batch {
			rows = append(rows, row)
		}
	}

	require.Equal(t, 4, len(rows))
	expectedOids := []string{"o1", "o2", "o3", "o4"}
	receivedOids := make([]string, 4, 4)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadBatchWithDelta(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	r := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
	csvReader, err := NewCsvReader(
		&source,
		outputs,
		WithCsvBatchSize(2),
		WithFirstLineHeader(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	var rows []repository.InputSchema

	eof := false
	for !eof {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			eof = true
		}
		for _, row := range batch {
			rows = append(rows, row)
		}
	}

	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3, 3)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadMultipleReaders(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	r := strings.NewReader(csv)
	r2 := strings.NewReader(csv)

	source := source.Source{
		Reader:      []source.SourceReader{{Reader: r}, {Reader: r2}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)
	csvReader, err := NewCsvReader(
		&source,
		outputs,
		WithCsvBatchSize(2),
		WithFirstLineHeader(true))
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 2, len(csvReader.csvReaders))

	var rows []repository.InputSchema

	eof := false
	for !eof {
		batch, err := csvReader.ReadBatch()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if err == io.EOF {
			eof = true
		}
		for _, row := range batch {
			rows = append(rows, row)
		}
	}

	require.Equal(t, 6, len(rows))
	expectedOids := []string{"o1", "o2", "o3", "o1", "o2", "o3"}
	receivedOids := make([]string, 6, 6)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestActor(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	r := strings.NewReader(csv)

	src := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 1)
	outputs[0] = make(chan reader.ReaderResult)

	csvReader, err := NewCsvReader(
		&src,
		outputs,
		WithCsvBatchSize(2),
	)
	if err != nil {
		t.Fatal(err)
	}
	csvReader.Start()

	var rows []repository.InputSchema
	for i := range csvReader.Outbox {
		for msg := range csvReader.Outbox[i] {
			for _, row := range msg.Rows {
				rows = append(rows, row)
			}
		}
	}
	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3, 3)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestActor_WithMultipleOutputs(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	r := strings.NewReader(csv)

	src := source.Source{
		Reader:      []source.SourceReader{{Reader: r}},
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	outputs := make([]chan reader.ReaderResult, 2)
	outputs[0] = make(chan reader.ReaderResult)
	outputs[1] = make(chan reader.ReaderResult)

	csvReader, err := NewCsvReader(
		&src,
		outputs,
		WithCsvBatchSize(2),
	)
	if err != nil {
		t.Fatal(err)
	}
	csvReader.Start()

	// read from all outbox concurrently
	// we need a lock to ensure safe access to the rows slice when appending
	// and a done channel to signal when all go routines are done
	var rows []repository.InputSchema
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := range csvReader.Outbox {
		wg.Add(1)
		go func() {
			for msg := range csvReader.Outbox[i] {
				mu.Lock()
				rows = append(rows, msg.Rows...)
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	require.Equal(t, 6, len(rows))
	expectedOids := []string{"o1", "o2", "o3", "o1", "o2", "o3"}
	receivedOids := make([]string, 6, 6)
	for i, row := range rows {
		receivedOids[i] = *row.ToMastercat(0).ID
	}

	// sort to be able to compare
	sort.Strings(receivedOids)
	sort.Strings(expectedOids)
	// compare
	require.Equal(t, expectedOids, receivedOids)
}
