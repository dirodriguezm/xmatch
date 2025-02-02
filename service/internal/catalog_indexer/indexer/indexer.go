package indexer

import (
	"log/slog"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type Row map[string]any

type ReaderResult struct {
	Rows  []repository.InputSchema
	Error error
}

type IndexerResult struct {
	Objects []repository.Mastercat
	Error   error
}

type WriterInput[T any] struct {
	Error error
	Rows  []T
}

type Reader interface {
	Start()
	Read() ([]repository.InputSchema, error)
	ReadBatch() ([]repository.InputSchema, error)
}

type Writer[T any] interface {
	Start()
	Done()
	Stop()
	Receive(WriterInput[T])
}

type Indexer struct {
	source *source.Source
	mapper *healpix.HEALPixMapper
	inbox  chan ReaderResult
	outbox chan WriterInput[repository.ParquetMastercat]
}

func New(
	src *source.Source,
	inbox chan ReaderResult,
	outbox chan WriterInput[repository.ParquetMastercat],
	cfg *config.IndexerConfig,
) (*Indexer, error) {
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
		ix.outbox <- WriterInput[repository.ParquetMastercat]{
			Rows:  nil,
			Error: msg.Error,
		}
		return
	}
	outputBatch := make([]repository.ParquetMastercat, len(msg.Rows))
	for i, row := range msg.Rows {
		mastercat := row.ToMastercat()
		point := healpix.RADec(*mastercat.Ra, *mastercat.Dec)
		ipix := ix.mapper.PixelAt(point)

		outputBatch[i] = repository.ParquetMastercat{
			Ra:   mastercat.Ra,
			Dec:  mastercat.Dec,
			ID:   mastercat.ID,
			Cat:  &ix.source.CatalogName,
			Ipix: &ipix,
		}
	}
	slog.Debug("Indexer sending message")
	ix.outbox <- WriterInput[repository.ParquetMastercat]{
		Rows:  outputBatch,
		Error: nil,
	}
}

func sendError(outbox chan WriterInput[repository.Mastercat], err error) {
	outbox <- WriterInput[repository.Mastercat]{
		Rows:  nil,
		Error: err,
	}
}
