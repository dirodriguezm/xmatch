package repository

import (
	"context"
	"database/sql"
)

type ErositaInputSchema struct {
	IAUNAME   string  `parquet:"name=IAUNAME, type=BYTE_ARRAY"`
	RA        float64 `parquet:"name=RA, type=DOUBLE"`
	DEC       float64 `parquet:"name=DEC, type=DOUBLE"`
	MJD       float64 `parquet:"name=MJD, type=DOUBLE"`
	ML_FLUX_1 float64 `parquet:"name=ML_FLUX_1, type=DOUBLE"`
}

func (schema ErositaInputSchema) GetCoordinates() (float64, float64) {
	return schema.RA, schema.DEC
}

func (schema ErositaInputSchema) GetId() string {
	return schema.IAUNAME
}

func (schema ErositaInputSchema) FillMetadata() Metadata {
	return Erosita{
		ID:      schema.GetId(),
		Mjd:     sql.NullFloat64{Float64: schema.MJD, Valid: true},
		MlFlux1: sql.NullFloat64{Float64: schema.ML_FLUX_1, Valid: true},
	}
}

func (schema ErositaInputSchema) FillMastercat(ipix int64) Mastercat {
	ra, dec := schema.GetCoordinates()
	return Mastercat{
		ID:   schema.GetId(),
		Ipix: ipix,
		Ra:   ra,
		Dec:  dec,
		Cat:  "erosita",
	}
}

func (e Erosita) GetId() string {
	return e.ID
}

func (e Erosita) GetCatalog() string {
	return "eROSITA"
}

func (m InsertErositaParams) GetId() string {
	return m.ID
}

func (m GetErositaFromPixelsRow) GetId() string {
	return m.ID
}

func (m GetErositaFromPixelsRow) GetCoordinates() (float64, float64) {
	return m.Ra, m.Dec
}

func (q *Queries) InsertErositaWithoutParams(ctx context.Context, arg Erosita) error {
	_, err := q.db.ExecContext(ctx, insertErosita, arg.ID, arg.Mjd, arg.MlFlux1)
	return err
}
