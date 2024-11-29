package indexer

import (
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Row map[string]any

type ReaderResult struct {
	Rows  []Row
	Error error
}

type IndexerResult struct {
	Objects []repository.Mastercat
	Error   error
}

type Reader interface {
	Start()
	Read() ([]Row, error)
	ReadBatch() ([]Row, error)
}

type Writer interface {
	Start()
	Done()
}

type Indexer struct {
	source *source.Source
	mapper *healpix.HEALPixMapper
	inbox  chan ReaderResult
	outbox chan IndexerResult
}

func New(src *source.Source, inbox chan ReaderResult, outbox chan IndexerResult, cfg *config.IndexerConfig) (*Indexer, error) {
	slog.Debug("Creating new Indexer")
	orderingScheme := healpix.Ring
	if strings.ToLower(cfg.OrderingScheme) == "nested" {
		orderingScheme = healpix.Nest
	}
	mapper, err := healpix.NewHEALPixMapper(src.Nside, orderingScheme)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		source: src,
		mapper: mapper,
		inbox:  inbox,
		outbox: outbox,
	}, nil
}

func (ix *Indexer) Start() {
	slog.Debug("Starting Indexer")
	go func() {
		defer func() {
			close(ix.outbox)
			slog.Debug("Closing Indexer")
		}()
		for msg := range ix.inbox {
			ix.receive(msg)
		}
	}()
}

func (ix *Indexer) receive(msg ReaderResult) {
	slog.Debug("Indexer Received Message")
	if msg.Error != nil {
		ix.outbox <- IndexerResult{
			Objects: nil,
			Error:   msg.Error,
		}
		return
	}
	outputBatch := make([]repository.Mastercat, 0)
	for _, row := range msg.Rows {
		mastercat, err := ix.Row2Mastercat(row)
		if err != nil {
			ix.outbox <- IndexerResult{
				Objects: nil,
				Error:   msg.Error,
			}
			return
		}
		outputBatch = append(outputBatch, repository.Mastercat{
			ID:   mastercat.ID,
			Ipix: mastercat.Ipix,
			Ra:   mastercat.Ra,
			Dec:  mastercat.Dec,
			Cat:  mastercat.Cat,
		})
	}
	slog.Debug("Indexer sending message")
	ix.outbox <- IndexerResult{
		Objects: outputBatch,
		Error:   nil,
	}
}

func (ix *Indexer) Row2Mastercat(row Row) (repository.Mastercat, error) {
	ra, err := ix.convertRa(row[ix.source.RaCol])
	if err != nil {
		return repository.Mastercat{}, err
	}
	dec, err := ix.convertDec(row[ix.source.DecCol])
	if err != nil {
		return repository.Mastercat{}, err
	}
	cat := ix.catalogName(row[ix.source.CatalogName])
	point := healpix.RADec(ra, dec)
	ipix := ix.mapper.PixelAt(point)
	mastercat := repository.Mastercat{
		ID:   row[ix.source.OidCol].(string),
		Ra:   ra,
		Dec:  dec,
		Cat:  cat,
		Ipix: ipix,
	}
	return mastercat, nil
}

func (ix *Indexer) convertRa(ra any) (float64, error) {
	switch reflect.TypeOf(ra).Kind() {
	case reflect.String:
		return strconv.ParseFloat(ra.(string), 64)
	default:
		return ra.(float64), nil
	}

}
func (ix *Indexer) convertDec(dec any) (float64, error) {
	switch reflect.TypeOf(dec).Kind() {
	case reflect.String:
		return strconv.ParseFloat(dec.(string), 64)
	default:
		return dec.(float64), nil
	}

}
func (ix *Indexer) catalogName(name any) string {
	if name == nil {
		return ix.source.CatalogName
	}
	return name.(string)
}
