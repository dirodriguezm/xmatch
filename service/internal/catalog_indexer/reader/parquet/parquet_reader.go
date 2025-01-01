package parquet_reader

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"regexp"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/xitongsys/parquet-go-source/local"
	preader "github.com/xitongsys/parquet-go/reader"
	psource "github.com/xitongsys/parquet-go/source"
)

type ParquetReader struct {
	*reader.BaseReader
	parquetReaders []*preader.ParquetReader
	src            *source.Source
	currentReader  int
	outbox         chan indexer.ReaderResult

	fileReaders []psource.ParquetFile
	recordType  reflect.Type
	metadata    []string
}

func createDynamicStruct(fields map[string]reflect.StructField) reflect.Type {
	structFields := make([]reflect.StructField, 0, len(fields))
	for _, field := range fields {
		structFields = append(structFields, field)
	}
	return reflect.StructOf(structFields)
}

func NewParquetReader(
	src *source.Source,
	channel chan indexer.ReaderResult,
	opts ...ParquetReaderOption,
) (*ParquetReader, error) {
	readers := []*preader.ParquetReader{}
	fileReaders := []psource.ParquetFile{}

	for _, srcReader := range src.Reader {
		fr, err := local.NewLocalFileReader(srcReader.Url)
		if err != nil {
			return nil, err
		}
		fileReaders = append(fileReaders, fr)
		pr, err := preader.NewParquetReader(fr, nil, 4)
		if err != nil {
			return nil, err
		}
		readers = append(readers, pr)
	}

	defaultMetadata := []string{
		fmt.Sprintf("name=%s, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY", src.OidCol),
		fmt.Sprintf("name=%s, type=DOUBLE", src.RaCol),
		fmt.Sprintf("name=%s, type=DOUBLE", src.DecCol),
	}

	newReader := &ParquetReader{
		BaseReader: &reader.BaseReader{
			Src:    src,
			Outbox: channel,
		},
		parquetReaders: readers,
		src:            src,
		currentReader:  0,
		outbox:         channel,
		fileReaders:    fileReaders,
		metadata:       defaultMetadata,
	}

	for _, opt := range opts {
		opt(newReader)
	}

	fields := createStructFields(newReader.metadata)
	newReader.recordType = createDynamicStruct(fields)

	newReader.Reader = newReader

	return newReader, nil
}

func (r *ParquetReader) ReadSingleFile(currentReader *preader.ParquetReader) ([]indexer.Row, error) {
	defer currentReader.ReadStop()

	nrows := currentReader.GetNumRows()
	records := reflect.MakeSlice(reflect.SliceOf(r.recordType), int(nrows), int(nrows))
	recordsPtr := reflect.New(records.Type())
	recordsPtr.Elem().Set(records)

	for range nrows {
		if err := currentReader.Read(recordsPtr.Interface()); err != nil {
			return nil, reader.NewReadError(r.currentReader, err, r.src, "Failed to read parquet")
		}
	}
	parsedRecords := convertToMapSlice(records, r.src.OidCol, r.src.RaCol, r.src.DecCol)
	return parsedRecords, nil
}

func (r *ParquetReader) Read() ([]indexer.Row, error) {
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

func (r *ParquetReader) ReadBatchSingleFile(currentReader *preader.ParquetReader) ([]indexer.Row, error) {
	records := reflect.MakeSlice(reflect.SliceOf(r.recordType), r.BatchSize, r.BatchSize)
	recordsPtr := reflect.New(records.Type())
	recordsPtr.Elem().Set(records)
	recordsInterface := recordsPtr.Interface()

	if err := currentReader.Read(recordsInterface); err != nil {
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

func (r *ParquetReader) ReadBatch() ([]indexer.Row, error) {
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

func convertToMapSlice(data reflect.Value, oidCol, raCol, decCol string) []indexer.Row {
	mapSlice := make([]indexer.Row, data.Len())
	for i := 0; i < data.Len(); i++ {
		mapSlice[i] = make(indexer.Row)
		elem := data.Index(i)
		caser := cases.Title(language.English, cases.Compact)
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

func getTypeFromMetadata(typeStr string) reflect.Type {
	switch typeStr {
	case "DOUBLE":
		return reflect.TypeOf(float64(0))
	case "BYTE_ARRAY":
		return reflect.TypeOf("")
	case "INT32":
		return reflect.TypeOf(int32(0))
	case "INT64":
		return reflect.TypeOf(int64(0))
	case "BOOLEAN":
		return reflect.TypeOf(bool(false))
	default:
		return reflect.TypeOf("") // default to string
	}
}

func parseMetadata(metadata string) (string, reflect.Type) {
	// Extract type from metadata string
	typeMatch := regexp.MustCompile(`type=(\w+)`).FindStringSubmatch(metadata)
	if len(typeMatch) < 2 {
		return metadata, reflect.TypeOf("")
	}

	return metadata, getTypeFromMetadata(typeMatch[1])
}

func createStructFields(metadata []string) map[string]reflect.StructField {
	fields := make(map[string]reflect.StructField)

	for _, md := range metadata {
		// Extract name from metadata
		nameMatch := regexp.MustCompile(`name=(\w+)`).FindStringSubmatch(md)
		if len(nameMatch) < 2 {
			continue
		}

		fieldName := cases.Title(language.English, cases.Compact).String(nameMatch[1])
		metadata, fieldType := parseMetadata(md)

		fields[fieldName] = reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`parquet:%s`, metadata)),
		}
	}

	return fields
}
