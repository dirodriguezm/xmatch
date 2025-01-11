package writer

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
)

type ParquetWriter struct {
	inbox chan indexer.WriterInput
	done  chan bool
}

func NewParquetWriter(inbox chan indexer.WriterInput, done chan bool) *ParquetWriter {
	slog.Debug("Creating new ParquetWriter")

	return &ParquetWriter{
		inbox: inbox,
		done:  done,
	}
}

func (w *ParquetWriter) Receive(msg indexer.WriterInput) {
	slog.Debug("ParquetWriter received message")
	if msg.Error != nil {

	}
}
