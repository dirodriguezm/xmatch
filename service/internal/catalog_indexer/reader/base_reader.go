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
	Outbox    []chan indexer.ReaderResult
}

func (r BaseReader) Start() {
	slog.Debug("Starting Reader", "catalog", r.Src.CatalogName, "nside", r.Src.Nside, "numreaders", len(r.Src.Reader))
	go func() {
		defer func() {
			for i := range r.Outbox {
				close(r.Outbox[i])
			}
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
				for i := range r.Outbox {
					r.Outbox[i] <- readResult
				}
				return
			}
			eof = err == io.EOF
			readResult := indexer.ReaderResult{
				Rows:  rows,
				Error: nil,
			}
			slog.Debug("Reader sending message")
			for i := range r.Outbox {
				r.Outbox[i] <- readResult
			}
		}
	}()
}
