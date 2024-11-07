package indexer

import (
	"context"
	"io"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/pkg/repository"
)

type Row map[string]any

type Reader interface {
	Read() ([]Row, error)
	ReadBatch() ([]Row, error)
	ObjectIdCol() string
	RaCol() string
	DecCol() string
	Catalog() string
}

type Indexer struct {
	DBRepository *repository.Queries
	Reader       Reader
	Mapper       *healpix.HEALPixMapper
}

func New(dbRepository *repository.Queries, reader Reader, nside int, scheme healpix.OrderingScheme) (*Indexer, error) {
	mapper, err := healpix.NewHEALPixMapper(nside, scheme)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		DBRepository: dbRepository,
		Reader:       reader,
		Mapper:       mapper,
	}, nil
}

func (ix *Indexer) Row2Mastercat(row Row) (repository.Mastercat, error) {
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

func (ix *Indexer) Index() error {
	eof := false
	ctx := context.Background()
	for !eof {
		// TODO: reader could read more batches while repository writes to DB
		batch, err := ix.Reader.ReadBatch()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			// we need to set eof true, but continue the execution
			// so that we can process remaining elements, even if EOF was reached
			eof = true
		}
		// Here we are going to insert one by one the items in the batch
		// TODO: create a bulk insert function
		for _, row := range batch {
			mastercat, err := ix.Row2Mastercat(row)
			if err != nil {
				return err
			}
			ix.DBRepository.InsertObject(ctx, repository.InsertObjectParams{
				ID:   mastercat.ID,
				Ipix: mastercat.Ipix,
				Ra:   mastercat.Ra,
				Dec:  mastercat.Dec,
				Cat:  mastercat.Cat,
			})
		}
	}
	return nil
}
