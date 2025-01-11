package writer

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
)

type BaseWriter struct {
	Writer       indexer.Writer
	DoneChannel  chan bool
	InboxChannel chan indexer.WriterInput
}

// Start starts the writer goroutine
func (w BaseWriter) Start() {
	slog.Debug("Starting Writer")

	go func() {
		defer func() {
			slog.Debug("Writer Done")
			w.Writer.Stop()
		}()
		for msg := range w.InboxChannel {
			w.Writer.Receive(msg)
		}
	}()
}

// Done blocks until the writer has finished processing all messages
func (w *BaseWriter) Done() {
	<-w.DoneChannel
}
