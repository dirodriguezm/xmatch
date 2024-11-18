package reader

import (
	"io"
	"strings"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	reader := strings.NewReader(csv)

	source := source.Source{
		Reader:      reader,
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, make(chan indexer.ReaderResult))
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
		receivedOids[i] = row["oid"].(string)
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadWithHeader(t *testing.T) {
	csv := `o1,1,1
o2,2,2
o3,3,3
`

	reader := strings.NewReader(csv)

	source := source.Source{
		Reader:      reader,
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, make(chan indexer.ReaderResult), WithHeader([]string{"oid", "ra", "dec"}))
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
		receivedOids[i] = row["oid"].(string)
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadBatch(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
o4,4,4
`
	reader := strings.NewReader(csv)

	source := source.Source{
		Reader:      reader,
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, make(chan indexer.ReaderResult), WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}

	var rows []indexer.Row

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
		receivedOids[i] = row["oid"].(string)
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestReadBatchWithDelta(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	reader := strings.NewReader(csv)

	source := source.Source{
		Reader:      reader,
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, make(chan indexer.ReaderResult), WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}

	var rows []indexer.Row

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
		receivedOids[i] = row["oid"].(string)
	}
	require.Equal(t, expectedOids, receivedOids)
}

func TestActor(t *testing.T) {
	csv := `oid,ra,dec
o1,1,1
o2,2,2
o3,3,3
`
	reader := strings.NewReader(csv)

	source := source.Source{
		Reader:      reader,
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "oid",
		CatalogName: "vlass",
	}

	csvReader, err := NewCsvReader(&source, make(chan indexer.ReaderResult), WithBatchSize(2))
	if err != nil {
		t.Fatal(err)
	}
	csvReader.Start()

	var rows []indexer.Row
	for msg := range csvReader.outbox {
		for _, row := range msg.Rows {
			rows = append(rows, row)
		}
	}
	require.Equal(t, 3, len(rows))
	expectedOids := []string{"o1", "o2", "o3"}
	receivedOids := make([]string, 3, 3)
	for i, row := range rows {
		receivedOids[i] = row["oid"].(string)
	}
	require.Equal(t, expectedOids, receivedOids)
}
