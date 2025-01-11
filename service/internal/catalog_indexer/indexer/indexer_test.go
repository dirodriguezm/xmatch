package indexer

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestIndexActor(t *testing.T) {
	inbox := make(chan ReaderResult)
	outbox := make(chan WriterInput)
	rows := []Row{
		{
			"ra":      1.0,
			"dec":     1.0,
			"id":      "o1",
			"catalog": "catalog",
		},
		{
			"ra":      2.0,
			"dec":     2.0,
			"id":      "o2",
			"catalog": "catalog",
		},
	}
	src, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:",
		Type:        "csv",
		CatalogName: "catalog",
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
	})
	require.NoError(t, err)
	indexer, err := New(src, inbox, outbox, &config.IndexerConfig{OrderingScheme: "nested", Nside: 18})
	require.NoError(t, err)

	indexer.Start()
	inbox <- ReaderResult{Rows: rows, Error: nil}
	close(inbox)
	results := make([]Row, 2)
	for msg := range outbox {
		for i, obj := range msg.Rows {
			results[i] = obj
		}
	}

	require.Len(t, results, 2)
}
