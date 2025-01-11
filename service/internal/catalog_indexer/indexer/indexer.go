package indexer

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type Row map[string]any

type ReaderResult struct {
	Rows  []Row
	Error error
}

type IndexerResult struct {
	Objects []Row
	Error   error
}

type WriterInput struct {
	Error error
	Rows  []Row
}

type Reader interface {
	Start()
	Read() ([]Row, error)
	ReadBatch() ([]Row, error)
}

type Writer interface {
	Start()
	Done()
	Stop()
	Receive(WriterInput)
}

type Indexer struct {
	source *source.Source
	mapper *healpix.HEALPixMapper
	inbox  chan ReaderResult
	outbox chan WriterInput
}

func New(src *source.Source, inbox chan ReaderResult, outbox chan WriterInput, cfg *config.IndexerConfig) (*Indexer, error) {
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
		ix.outbox <- WriterInput{
			Rows:  nil,
			Error: msg.Error,
		}
		return
	}
	outputBatch := make([]Row, len(msg.Rows))
	for i, row := range msg.Rows {
		ra, err := ix.convertRa(row[ix.source.RaCol])
		if err != nil {
			sendError(ix.outbox, err)
		}
		dec, err := ix.convertDec(row[ix.source.DecCol])
		if err != nil {
			sendError(ix.outbox, err)
		}

		point := healpix.RADec(ra, dec)
		ipix := ix.mapper.PixelAt(point)

		outputBatch[i] = Row{
			ix.source.RaCol:  ra,
			ix.source.DecCol: dec,
			ix.source.OidCol: row[ix.source.OidCol],
			"cat":            ix.source.CatalogName,
			"ipix":           ipix,
		}
	}
	slog.Debug("Indexer sending message")
	ix.outbox <- WriterInput{
		Rows:  outputBatch,
		Error: nil,
	}
}

func sendError(outbox chan WriterInput, err error) {
	outbox <- WriterInput{
		Rows:  nil,
		Error: err,
	}
}

func (ix *Indexer) convertRa(ra any) (float64, error) {
	switch v := ra.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case *float64:
		return *v, nil
	case float64:
		return v, nil
	default:
		return v.(float64), nil
	}
}
func (ix *Indexer) convertDec(dec any) (float64, error) {
	switch v := dec.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case *float64:
		return *v, nil
	case float64:
		return v, nil
	default:
		return v.(float64), nil
	}
}
