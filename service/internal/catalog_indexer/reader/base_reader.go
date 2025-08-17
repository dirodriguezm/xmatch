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
	GetOutbox() []chan ReaderResult
}

type BaseReader struct {
	Reader    Reader
	Src       *source.Source
	BatchSize int
	Outbox    []chan ReaderResult
}

func (r BaseReader) Start() {
	slog.Debug("Starting Reader", "catalog", r.Src.CatalogName, "nside", r.Src.Nside, "numreaders", len(r.Src.Sources))
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
				// If the error is not EOF, it means that something went wrong reading the file
				readResult := ReaderResult{
					Rows:  nil,
					Error: err,
				}
				// We send the message containing the error
				for i := range r.Outbox {
					r.Outbox[i] <- readResult
				}
				return
			}
			// We update the eof variable so that we can stop the loop when all files are read
			eof = err == io.EOF
			// Now we can send the actual rows to all the receivers
			readResult := ReaderResult{
				Rows:  rows,
				Error: nil,
			}
			slog.Debug("Reader sending message")
			for i := range r.Outbox {
				r.Outbox[i] <- readResult
			}
			rows = nil
			readResult.Rows = nil
		}
	}()
}

func (r BaseReader) GetOutbox() []chan ReaderResult {
	return r.Outbox
}
