package parquet_writer

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	pwriter "github.com/xitongsys/parquet-go/writer"
)

type ParquetWriter[T any] struct {
	*writer.BaseWriter
	parquetWriter *pwriter.ParquetWriter
	pfile         *os.File
	OutputFile    string
}

func NewParquetWriter[T any](
	inbox chan indexer.WriterInput,
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
		BaseWriter: &writer.BaseWriter{
			InboxChannel: inbox,
			DoneChannel:  done,
		},
	}
	w.Writer = w
	return w, nil
}

func (w *ParquetWriter[T]) Receive(msg indexer.WriterInput) {
	slog.Debug("ParquetWriter received message")
	if msg.Error != nil {
		slog.Error("ParquetWriter received error message")
		panic(msg.Error)
	}

	for i := 0; i < len(msg.Rows); i++ {
		object := convertMapToStruct[T](msg.Rows[i])
		if err := w.parquetWriter.Write(object); err != nil {
			panic(fmt.Errorf("ParquetWriter could not write object %v\n%w", object, err))
		}
	}
}

func (w *ParquetWriter[T]) Stop() {
	if err := w.parquetWriter.WriteStop(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not stop %w", err))
	}
	if err := w.pfile.Close(); err != nil {
		panic(fmt.Errorf("ParquetWriter could not close parquet file %w", err))
	}
	w.DoneChannel <- true
	close(w.DoneChannel)
}

func convertMapToStruct[T any](data indexer.Row) T {
	var result T
	resultValue := reflect.ValueOf(&result).Elem()

	// Helper function to convert snake_case to PascalCase
	toPascalCase := func(s string) string {
		parts := strings.Split(s, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
		}
		return strings.Join(parts, "")
	}

	// Convert each map key to a potential struct field name and set the value
	for key, value := range data {
		fieldName := toPascalCase(key)

		// Try to find and set the field
		if field := resultValue.FieldByName(fieldName); field.IsValid() && field.CanSet() {
			converted := reflect.ValueOf(value)
			if converted.Type().ConvertibleTo(field.Type()) {
				field.Set(converted.Convert(field.Type()))
			}
		}
	}

	return result
}
