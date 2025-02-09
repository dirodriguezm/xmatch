package mastercat_indexer

import (
	"log/slog"
	"strings"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type IndexerActor struct {
	source *source.Source
	mapper *healpix.HEALPixMapper
	inbox  chan indexer.ReaderResult
	outbox chan indexer.WriterInput[repository.ParquetMastercat]
}

func New(
	src *source.Source,
	inbox chan indexer.ReaderResult,
	outbox chan indexer.WriterInput[repository.ParquetMastercat],
	cfg *config.IndexerConfig,
) (*IndexerActor, error) {
	slog.Debug("Creating new Indexer")
	orderingScheme := healpix.Ring
	if strings.ToLower(cfg.OrderingScheme) == "nested" {
		orderingScheme = healpix.Nest
	}
	mapper, err := healpix.NewHEALPixMapper(src.Nside, orderingScheme)
	if err != nil {
		return nil, err
	}
	return &IndexerActor{
		source: src,
		mapper: mapper,
		inbox:  inbox,
		outbox: outbox,
	}, nil
}

func (actor *IndexerActor) Start() {
	slog.Debug("Starting Indexer")
	go func() {
		defer func() {
			close(actor.outbox)
			slog.Debug("Closing Indexer")
		}()
		for msg := range actor.inbox {
			actor.receive(msg)
		}
	}()
}

func (actor *IndexerActor) receive(msg indexer.ReaderResult) {
	slog.Debug("Indexer Received Message")
	if msg.Error != nil {
		actor.outbox <- indexer.WriterInput[repository.ParquetMastercat]{
			Rows:  nil,
			Error: msg.Error,
		}
		return
	}
	outputBatch := make([]repository.ParquetMastercat, len(msg.Rows))
	for i, row := range msg.Rows {
		mastercat := row.ToMastercat()
		point := healpix.RADec(*mastercat.Ra, *mastercat.Dec)
		ipix := actor.mapper.PixelAt(point)

		outputBatch[i] = repository.ParquetMastercat{
			Ra:   mastercat.Ra,
			Dec:  mastercat.Dec,
			ID:   mastercat.ID,
			Cat:  &actor.source.CatalogName,
			Ipix: &ipix,
		}
	}
	slog.Debug("Indexer sending message")
	actor.outbox <- indexer.WriterInput[repository.ParquetMastercat]{
		Rows:  outputBatch,
		Error: nil,
	}
}

func sendError(outbox chan indexer.WriterInput[repository.Mastercat], err error) {
	outbox <- indexer.WriterInput[repository.Mastercat]{
		Rows:  nil,
		Error: err,
	}
}
