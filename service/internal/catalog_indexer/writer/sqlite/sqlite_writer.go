package sqlite_writer

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type SqliteWriter struct {
	*writer.BaseWriter
	repository conesearch.Repository
	ctx        context.Context
	src        *source.Source
}

func NewSqliteWriter(repository conesearch.Repository, ch chan indexer.WriterInput, done chan bool, ctx context.Context, src *source.Source) *SqliteWriter {
	slog.Debug("Creating new SqliteWriter")
	w := &SqliteWriter{
		BaseWriter: &writer.BaseWriter{
			DoneChannel:  done,
			InboxChannel: ch,
		},
		repository: repository,
		ctx:        ctx,
		src:        src,
	}
	w.Writer = w
	return w
}

func (w *SqliteWriter) Receive(msg indexer.WriterInput) {
	slog.Debug("Writer received message")
	if msg.Error != nil {
		slog.Error("SqliteWriter received error")
		panic(msg.Error)
	}
	for _, object := range msg.Rows {
		// convert the received row to insert params needed by the repository
		params, err := row2insertParams(object, w.src)
		if err != nil {
			slog.Error("SqliteWriter could not convert received object to insert params", "object", object)
			panic(err)
		}

		// insert converted rows
		_, err = w.repository.InsertObject(w.ctx, params)
		if err != nil {
			slog.Error("SqliteWriter could not write object to database", "object", object)
			panic(err)
		}
	}
}

func (w *SqliteWriter) Stop() {
	w.DoneChannel <- true
}

func row2insertParams(obj indexer.Row, src *source.Source) (repository.InsertObjectParams, error) {
	oid, err := convertOid(obj[src.OidCol])
	if err != nil {
		return repository.InsertObjectParams{}, err
	}

	ipix, err := convertIpix(obj["ipix"])
	if err != nil {
		return repository.InsertObjectParams{}, err
	}

	ra, err := convertRa(obj[src.RaCol])
	if err != nil {
		return repository.InsertObjectParams{}, err
	}

	dec, err := convertDec(obj[src.DecCol])
	if err != nil {
		return repository.InsertObjectParams{}, err
	}

	return repository.InsertObjectParams{
		ID:   oid,
		Ipix: ipix,
		Ra:   ra,
		Dec:  dec,
		Cat:  src.CatalogName,
	}, nil
}

func convertOid(oid any) (string, error) {
	switch v := oid.(type) {
	case string:
		return oid.(string), nil
	case *string:
		return *v, nil
	default:
		return v.(string), nil
	}
}

func convertIpix(ipix any) (int64, error) {
	switch v := ipix.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int:
		return int64(v), nil
	default:
		return v.(int64), nil
	}
}

func convertRa(ra any) (float64, error) {
	switch v := ra.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case *float64:
		return *v, nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return v.(float64), nil
	}
}
func convertDec(dec any) (float64, error) {
	switch v := dec.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case *float64:
		return *v, nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return v.(float64), nil
	}
}
