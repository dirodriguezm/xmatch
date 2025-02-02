package writer

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
)

type BaseWriter[T any] struct {
	Writer       indexer.Writer[T]
	DoneChannel  chan bool
	InboxChannel chan indexer.WriterInput[T]
}

// Start starts the writer goroutine
func (w BaseWriter[T]) Start() {
	slog.Debug("Starting Writer")

	go func() {
		for msg := range w.InboxChannel {
			w.Writer.Receive(msg)
		}
		w.Writer.Stop()
		slog.Debug("Writer Done")
	}()
}

// Done blocks until the writer has finished processing all messages
func (w *BaseWriter[T]) Done() {
	<-w.DoneChannel
}
