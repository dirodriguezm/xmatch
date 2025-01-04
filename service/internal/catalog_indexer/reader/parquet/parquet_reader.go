package parquet_reader

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/xitongsys/parquet-go-source/local"
	preader "github.com/xitongsys/parquet-go/reader"
	psource "github.com/xitongsys/parquet-go/source"
)

type ParquetReader[T any] struct {
	*reader.BaseReader
	parquetReaders []*preader.ParquetReader
	src            *source.Source
	currentReader  int
	outbox         chan indexer.ReaderResult

	fileReaders []psource.ParquetFile
}

func NewParquetReader[T any](
	src *source.Source,
	channel chan indexer.ReaderResult,
	opts ...ParquetReaderOption[T],
) (*ParquetReader[T], error) {
	readers := []*preader.ParquetReader{}
	fileReaders := []psource.ParquetFile{}

	for _, srcReader := range src.Reader {
		fr, err := local.NewLocalFileReader(srcReader.Url)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewLocalFileReader\n%w", err)
		}
		fileReaders = append(fileReaders, fr)
		schema := new(T)
		pr, err := preader.NewParquetReader(fr, schema, 4)
		if err != nil {
			return nil, fmt.Errorf("Could not create NewParquetReader\n%w", err)
		}
		readers = append(readers, pr)
	}

	newReader := &ParquetReader[T]{
		BaseReader: &reader.BaseReader{
			Src:    src,
			Outbox: channel,
		},
		src:            src,
		currentReader:  0,
		outbox:         channel,
		fileReaders:    fileReaders,
		parquetReaders: readers,
	}

	for _, opt := range opts {
		opt(newReader)
	}

	newReader.Reader = newReader
	return newReader, nil
}

func (r *ParquetReader[T]) ReadSingleFile(currentReader *preader.ParquetReader) ([]indexer.Row, error) {
	defer currentReader.ReadStop()

	nrows := currentReader.GetNumRows()
	records := make([]T, nrows)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(r.currentReader, err, r.src, "Failed to read parquet")
	}

	parsedRecords := convertToMapSlice(records, r.src.OidCol, r.src.RaCol, r.src.DecCol)
	return parsedRecords, nil
}

func (r *ParquetReader[T]) Read() ([]indexer.Row, error) {
	defer func() {
		for i := range r.fileReaders {
			r.fileReaders[i].Close()
		}
	}()

	rows := make([]indexer.Row, 0, 0)
	for _, currentReader := range r.parquetReaders {
		currentRows, err := r.ReadSingleFile(currentReader)
		if err != nil {
			return nil, err
		}
		rows = append(rows, currentRows...)
	}
	return rows, nil
}

func (r *ParquetReader[T]) ReadBatchSingleFile(currentReader *preader.ParquetReader) ([]indexer.Row, error) {
	records := make([]T, r.BatchSize)

	if err := currentReader.Read(&records); err != nil {
		return nil, reader.NewReadError(r.currentReader, err, r.src, "Failed to read parquet in batch")
	}

	parsedRecords := convertToMapSlice(records, r.src.OidCol, r.src.RaCol, r.src.DecCol)

	if isZeroValueSlice(parsedRecords) {
		// finished reading
		currentReader.ReadStop()
		return nil, io.EOF
	}

	return parsedRecords, nil
}

func (r *ParquetReader[T]) ReadBatch() ([]indexer.Row, error) {
	currentReader := r.parquetReaders[r.currentReader]
	slog.Debug("ParquetReader reading batch")
	currentRows, err := r.ReadBatchSingleFile(currentReader)
	if err != nil {
		if err == io.EOF && r.currentReader < len(r.parquetReaders)-1 {
			// only finished current file, but more files remain
			r.currentReader += 1
			return currentRows, nil
		}
		if err == io.EOF {
			// finished reading all files
			r.currentReader += 1
			for i := range r.fileReaders {
				r.fileReaders[i].Close()
			}
			return currentRows, err
		}
		return nil, err
	}
	return currentRows, nil
}

func convertToMapSlice[T any](data []T, oidCol, raCol, decCol string) []indexer.Row {
	mapSlice := make([]indexer.Row, len(data))
	caser := cases.Title(language.English, cases.Compact)
	for i := 0; i < len(data); i++ {
		mapSlice[i] = make(indexer.Row)
		elem := reflect.ValueOf(data[i])
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		// Check if it's a struct
		if elem.Kind() != reflect.Struct {
			panic(fmt.Errorf("expected struct, got %v", elem.Kind()))
		}

		mapSlice[i][oidCol] = elem.FieldByName(caser.String(oidCol)).Interface()
		mapSlice[i][raCol] = elem.FieldByName(caser.String(raCol)).Interface()
		mapSlice[i][decCol] = elem.FieldByName(caser.String(decCol)).Interface()
	}
	return mapSlice
}

func isZeroValueSlice(s []indexer.Row) bool {
	for i := 0; i < len(s); i++ {
		if !isZeroValueMap(s[i]) {
			return false
		}
	}
	return true
}

func isZeroValueMap(m map[string]interface{}) bool {
	for _, v := range m {
		if !isZeroValue(v) {
			return false
		}
	}
	return true
}

// Helper function to check if an interface{} value is zero
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch value := v.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(value).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(value).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(value).Float() == 0
	case bool:
		return !value
	case string:
		return value == ""
	case []interface{}:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	default:
		// For struct types, compare with their zero value
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}
