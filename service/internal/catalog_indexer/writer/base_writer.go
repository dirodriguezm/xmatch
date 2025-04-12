package writer

import (
	"log/slog"
)

type WriterInput[T any] struct {
	Error error
	Rows  []T
}

type Writer[T any] interface {
	Start()
	Done()
	Stop()
	Receive(WriterInput[T])
}

type BaseWriter[T any] struct {
	Writer       Writer[T]
	DoneChannel  chan bool
	InboxChannel chan WriterInput[T]
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
