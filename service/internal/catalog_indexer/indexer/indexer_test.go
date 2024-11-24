package indexer

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRow2Mastercat(t *testing.T) {
	row := Row{
		"ra":      1.0,
		"dec":     1.0,
		"id":      "test",
		"catalog": "catalog",
	}
	var result repository.Mastercat
	src, err := source.NewSource(&config.SourceConfig{
		Url:         "buffer:",
		Type:        "csv",
		CatalogName: "catalog",
		RaCol:       "ra",
		DecCol:      "dec",
		OidCol:      "id",
	})
	require.NoError(t, err)
	ix, err := New(src, nil, nil, &config.IndexerConfig{OrderingScheme: "nested", Nside: 18})
	require.NoError(t, err)
	result, err = ix.Row2Mastercat(row)
	require.NoError(t, err)
	require.Equal(t, result.ID, "test")
	require.Equal(t, result.Ra, 1.0)
	require.Equal(t, result.Dec, 1.0)
	require.Equal(t, result.Cat, "catalog")
	require.NotNil(t, result.Ipix)
}

func TestIndexActor(t *testing.T) {
	inbox := make(chan ReaderResult)
	outbox := make(chan IndexerResult)
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
	results := make([]repository.Mastercat, 2)
	for msg := range outbox {
		for i, obj := range msg.Objects {
			results[i] = obj
		}
	}

	require.Len(t, results, 2)
}
