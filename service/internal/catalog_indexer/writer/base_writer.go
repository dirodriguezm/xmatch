// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	DoneChannel  chan struct{}
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
