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
	"context"
	"log/slog"
)

type WriterInput[T any] struct {
	Error error
	Rows  []T
}

type Writer[I any, O any] interface {
	Start()
	Done()
	Stop()
	Receive(WriterInput[I])
}

type BaseWriter[I any, O any] struct {
	Ctx          context.Context
	Writer       Writer[I, O]
	DoneChannel  chan struct{}
	InboxChannel chan WriterInput[I]
}

// Start starts the writer goroutine
func (w BaseWriter[I, O]) Start() {
	slog.Debug("Starting Writer")

	go func() {
		for {
			select {
			case <-w.Ctx.Done():
				// Case for context cancellation
				slog.Debug("Writer context cancellation")
				w.Writer.Stop()
				return
			case msg, ok := <-w.InboxChannel:
				// Case for closed channel
				if !ok {
					slog.Debug("Writer Done")
					w.Writer.Stop()
					return
				}
				// Write the received message
				w.Writer.Receive(msg)
			}
		}
	}()
}

// Done blocks until the writer has finished processing all messages
func (w *BaseWriter[I, O]) Done() {
	<-w.DoneChannel
}
