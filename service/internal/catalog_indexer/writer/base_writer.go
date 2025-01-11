package writer

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
)

type BaseWriter struct {
	Writer indexer.Writer
	Done   chan bool
	Inbox  chan indexer.WriterInput
}

func (w BaseWriter) Start() {
	slog.Debug("Starting Writer")

	go func() {
		defer func() {
			slog.Debug("Writer Done")
			w.Done <- true
			close(w.Done)
		}()
		for msg := range w.Inbox {
			w.Writer.Receive(msg)
		}
	}()
}
