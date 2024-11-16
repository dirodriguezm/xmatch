package indexer

import (
	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/assertions"
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
	ObjectIdCol() string
	RaCol() string
	DecCol() string
	Catalog() string
}

type Writer interface {
	Start()
}

type Indexer struct {
	DBRepository *repository.Queries
	Reader       Reader
	Mapper       *healpix.HEALPixMapper
	inbox        chan ReaderResult
	outbox       chan IndexerResult
}

func New(reader Reader, nside int, scheme healpix.OrderingScheme, inbox chan ReaderResult, outbox chan IndexerResult) (*Indexer, error) {
	mapper, err := healpix.NewHEALPixMapper(nside, scheme)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		Reader: reader,
		Mapper: mapper,
		inbox:  inbox,
		outbox: outbox,
	}, nil
}

func (ix *Indexer) Start() {
	go func() {
		defer close(ix.outbox)
		for msg := range ix.inbox {
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
			ix.outbox <- IndexerResult{
				Objects: outputBatch,
				Error:   nil,
			}
		}
	}()
}

func (ix *Indexer) Row2Mastercat(row Row) (repository.Mastercat, error) {
	assertions.HasKey(row, ix.Reader.RaCol(), "Row didn't contain %s", ix.Reader.RaCol())
	assertions.HasKey(row, ix.Reader.DecCol(), "Row didn't contain %s", ix.Reader.DecCol())
	assertions.HasKey(row, ix.Reader.ObjectIdCol(), "Row didn't contain %s", ix.Reader.ObjectIdCol())
	assertions.HasKey(row, ix.Reader.Catalog(), "Row didn't contain %s", ix.Reader.Catalog())

	ra := row[ix.Reader.RaCol()].(float64)
	dec := row[ix.Reader.DecCol()].(float64)
	point := healpix.RADec(ra, dec)
	ipix := ix.Mapper.PixelAt(point)
	mastercat := repository.Mastercat{
		ID:   row[ix.Reader.ObjectIdCol()].(string),
		Ra:   ra,
		Dec:  dec,
		Cat:  row[ix.Reader.Catalog()].(string),
		Ipix: ipix,
	}
	return mastercat, nil
}
