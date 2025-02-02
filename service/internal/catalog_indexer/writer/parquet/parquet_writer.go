package parquet_writer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	pwriter "github.com/xitongsys/parquet-go/writer"
)

type ParquetWriter[T any] struct {
	*writer.BaseWriter[T]
	parquetWriter *pwriter.ParquetWriter
	pfile         *os.File
	OutputFile    string
}

func NewParquetWriter[T any](
	inbox chan indexer.WriterInput[T],
	done chan bool,
	cfg *config.WriterConfig,
) (*ParquetWriter[T], error) {
	slog.Debug("Creating new ParquetWriter")

	file, err := os.Create(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("ParquetReader could not create file %s\n%w", cfg.OutputFile, err)
	}

	schema := new(T)
	parquetWriter, err := pwriter.NewParquetWriterFromWriter(file, schema, 1)
	if err != nil {
		return nil, fmt.Errorf("ParquetReader could now create writer %w", err)
	}

	w := &ParquetWriter[T]{
		parquetWriter: parquetWriter,
		pfile:         file,
		BaseWriter: &writer.BaseWriter[T]{
			InboxChannel: inbox,
			DoneChannel:  done,
		},
	}
	w.Writer = w
	return w, nil
}

func (w *ParquetWriter[T]) Receive(msg indexer.WriterInput[T]) {
	slog.Debug("ParquetWriter received message")
	if msg.Error != nil {
		slog.Error("ParquetWriter received error message")
		panic(msg.Error)
	}

	for i := 0; i < len(msg.Rows); i++ {
		obj := msg.Rows[i]
		if err := w.parquetWriter.Write(obj); err != nil {
			panic(fmt.Errorf("ParquetWriter could not write object %v\n%w", obj, err))
		}
	}
	slog.Debug("ParquetWriter wrote messages", "messages", len(msg.Rows))
}

func (w *ParquetWriter[T]) Stop() {
	if err := w.parquetWriter.WriteStop(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not stop. Error: %w", err))
	}
	if err := w.pfile.Close(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not close parquet file %w", err))
	}
	w.DoneChannel <- true
	close(w.DoneChannel)
}
