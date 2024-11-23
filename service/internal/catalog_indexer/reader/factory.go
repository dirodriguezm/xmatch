package reader

import (
	"fmt"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

func ReaderFactory(src *source.Source, outbox chan indexer.ReaderResult, cfg *config.ReaderConfig) (indexer.Reader, error) {
	readerType := strings.ToLower(cfg.Type)
	switch readerType {
	case "csv":
		return NewCsvReader(
			src,
			outbox,
			WithBatchSize(cfg.BatchSize),
			WithHeader(cfg.Header),
			WithFirstLineHeader(cfg.FirstLineHeader),
		)
	default:
		return nil, fmt.Errorf("Reader type not allowed")
	}
}
