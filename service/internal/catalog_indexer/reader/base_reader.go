package reader

import (
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
)

type BaseReader struct {
	Reader    indexer.Reader
	Src       *source.Source
	BatchSize int
	Outbox    chan indexer.ReaderResult
}

func (r BaseReader) Start() {
	slog.Debug("Starting Reader", "catalog", r.Src.CatalogName, "nside", r.Src.Nside, "numreaders", len(r.Src.Reader))
	go func() {
		defer func() {
			close(r.Outbox)
			slog.Debug("Closing Reader")
		}()
		eof := false
		for !eof {
			rows, err := r.Reader.ReadBatch()
			if err != nil && err != io.EOF {
				readResult := indexer.ReaderResult{
					Rows:  nil,
					Error: err,
				}
				r.Outbox <- readResult
				return
			}
			eof = err == io.EOF
			readResult := indexer.ReaderResult{
				Rows:  rows,
				Error: nil,
			}
			slog.Debug("Reader sending message")
			r.Outbox <- readResult
		}
	}()
}
