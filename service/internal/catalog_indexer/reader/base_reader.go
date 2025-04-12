package reader

import (
	"io"
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type ReaderResult struct {
	Rows  []repository.InputSchema
	Error error
}

type Reader interface {
	Start()
	Read() ([]repository.InputSchema, error)
	ReadBatch() ([]repository.InputSchema, error)
}

type BaseReader struct {
	Reader    Reader
	Src       *source.Source
	BatchSize int
	Outbox    []chan ReaderResult
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
				readResult := ReaderResult{
					Rows:  nil,
					Error: err,
				}
				for i := range r.Outbox {
					r.Outbox[i] <- readResult
				}
				return
			}
			eof = err == io.EOF
			readResult := ReaderResult{
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
