package erosita

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/healpix"
	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	fits_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/fits"
	parquet_reader "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
	parquet_writer "github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer/parquet"
	"github.com/dirodriguezm/xmatch/service/internal/config"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type ErositaStore interface {
	InsertErositaWithoutParams(context.Context, repository.Erosita) error
	GetErosita(context.Context, string) (repository.GetErositaRow, error)
	BulkInsertErosita(context.Context, *sql.DB, []any) error
	BulkGetErosita(context.Context, []string) ([]repository.BulkGetErositaRow, error)
	GetErositaFromPixels(context.Context, []int64) ([]repository.GetErositaFromPixelsRow, error)
}

type Adapter struct {
	store ErositaStore
}

func init() {
	catalog.Register("erosita", func(store any) (catalog.CatalogAdapter, error) {
		s, _ := store.(ErositaStore)
		return &Adapter{store: s}, nil
	})
}

func (a Adapter) Name() string {
	return "erosita"
}

func (a Adapter) NewRawRecord() any {
	return repository.ErositaInputSchema{}
}

func (a Adapter) NewParquetWriter(cfg config.WriterConfig, ctx context.Context) (writer.Writer, error) {
	return parquet_writer.New[repository.Erosita](cfg, ctx)
}

func (a Adapter) NewParquetReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return parquet_reader.NewParquetReader(
		src,
		parquet_reader.WithParquetBatchSize[repository.ErositaInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) NewFitsReader(src *source.Source, cfg config.ReaderConfig) (reader.Reader, error) {
	return fits_reader.NewFitsReader(
		src,
		fits_reader.WithBatchSize[repository.ErositaInputSchema](cfg.BatchSize),
	)
}

func (a Adapter) BulkInsertFn() func(context.Context, *sql.DB, []any) error {
	if a.store == nil {
		return func(ctx context.Context, db *sql.DB, rows []any) error {
			return fmt.Errorf("erosita adapter has no store")
		}
	}
	return a.store.BulkInsertErosita
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("erosita adapter has no store")
	}
	return a.store.GetErosita(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.store == nil {
		return nil, fmt.Errorf("erosita adapter has no store")
	}
	return a.store.BulkGetErosita(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.MetadataWithCoordinates, error) {
	if a.store == nil {
		return nil, fmt.Errorf("erosita adapter has no store")
	}
	rows, err := a.store.GetErositaFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.MetadataWithCoordinates, len(rows))
	for i, r := range rows {
		result[i] = r
	}
	return result, nil
}

func (a Adapter) ConvertToMetadata(obj repository.MetadataWithCoordinates) repository.Metadata {
	row := obj.(repository.GetErositaFromPixelsRow)
	return repository.Erosita{
		ID:             row.ID,
		Detuid:         row.Detuid,
		Skytile:        row.Skytile,
		IDSrc:          row.IDSrc,
		Uid:            row.Uid,
		UidHard:        row.UidHard,
		IDCluster:      row.IDCluster,
		Ra:             row.Ra,
		Dec:            row.Dec,
		RaLowerr:       row.RaLowerr,
		RaUperr:        row.RaUperr,
		DecLowerr:      row.DecLowerr,
		DecUperr:       row.DecUperr,
		PosErr:         row.PosErr,
		Mjd:            row.Mjd,
		MjdMin:         row.MjdMin,
		MjdMax:         row.MjdMax,
		Ext:            row.Ext,
		ExtErr:         row.ExtErr,
		ExtLike:        row.ExtLike,
		DetLike0:       row.DetLike0,
		MlCts1:         row.MlCts1,
		MlCtsErr1:      row.MlCtsErr1,
		MlRate1:        row.MlRate1,
		MlRateErr1:     row.MlRateErr1,
		MlFlux1:        row.MlFlux1,
		MlFluxErr1:     row.MlFluxErr1,
		MlBkg1:         row.MlBkg1,
		MlExp1:         row.MlExp1,
		ApeBkg1:        row.ApeBkg1,
		ApeRadius1:     row.ApeRadius1,
		ApePois1:       row.ApePois1,
		DetLikeP1:      row.DetLikeP1,
		MlCtsP1:        row.MlCtsP1,
		MlCtsErrP1:     row.MlCtsErrP1,
		MlRateP1:       row.MlRateP1,
		MlRateErrP1:    row.MlRateErrP1,
		MlFluxP1:       row.MlFluxP1,
		MlFluxErrP1:    row.MlFluxErrP1,
		MlBkgP1:        row.MlBkgP1,
		MlExpP1:        row.MlExpP1,
		ApeBkgP1:       row.ApeBkgP1,
		ApeRadiusP1:    row.ApeRadiusP1,
		ApePoisP1:      row.ApePoisP1,
		DetLikeP2:      row.DetLikeP2,
		MlCtsP2:        row.MlCtsP2,
		MlCtsErrP2:     row.MlCtsErrP2,
		MlRateP2:       row.MlRateP2,
		MlRateErrP2:    row.MlRateErrP2,
		MlFluxP2:       row.MlFluxP2,
		MlFluxErrP2:    row.MlFluxErrP2,
		MlBkgP2:        row.MlBkgP2,
		MlExpP2:        row.MlExpP2,
		ApeBkgP2:       row.ApeBkgP2,
		ApeRadiusP2:    row.ApeRadiusP2,
		ApePoisP2:      row.ApePoisP2,
		DetLikeP3:      row.DetLikeP3,
		MlCtsP3:        row.MlCtsP3,
		MlCtsErrP3:     row.MlCtsErrP3,
		MlRateP3:       row.MlRateP3,
		MlRateErrP3:    row.MlRateErrP3,
		MlFluxP3:       row.MlFluxP3,
		MlFluxErrP3:    row.MlFluxErrP3,
		MlBkgP3:        row.MlBkgP3,
		MlExpP3:        row.MlExpP3,
		ApeBkgP3:       row.ApeBkgP3,
		ApeRadiusP3:    row.ApeRadiusP3,
		ApePoisP3:      row.ApePoisP3,
		DetLikeP4:      row.DetLikeP4,
		MlCtsP4:        row.MlCtsP4,
		MlCtsErrP4:     row.MlCtsErrP4,
		MlRateP4:       row.MlRateP4,
		MlRateErrP4:    row.MlRateErrP4,
		MlFluxP4:       row.MlFluxP4,
		MlFluxErrP4:    row.MlFluxErrP4,
		MlBkgP4:        row.MlBkgP4,
		MlExpP4:        row.MlExpP4,
		ApeBkgP4:       row.ApeBkgP4,
		ApeRadiusP4:    row.ApeRadiusP4,
		ApePoisP4:      row.ApePoisP4,
		DetLikeP5:      row.DetLikeP5,
		MlCtsP5:        row.MlCtsP5,
		MlCtsErrP5:     row.MlCtsErrP5,
		MlRateP5:       row.MlRateP5,
		MlRateErrP5:    row.MlRateErrP5,
		MlFluxP5:       row.MlFluxP5,
		MlFluxErrP5:    row.MlFluxErrP5,
		MlBkgP5:        row.MlBkgP5,
		MlExpP5:        row.MlExpP5,
		ApeBkgP5:       row.ApeBkgP5,
		ApeRadiusP5:    row.ApeRadiusP5,
		ApePoisP5:      row.ApePoisP5,
		DetLikeP6:      row.DetLikeP6,
		MlCtsP6:        row.MlCtsP6,
		MlCtsErrP6:     row.MlCtsErrP6,
		MlRateP6:       row.MlRateP6,
		MlRateErrP6:    row.MlRateErrP6,
		MlFluxP6:       row.MlFluxP6,
		MlFluxErrP6:    row.MlFluxErrP6,
		MlBkgP6:        row.MlBkgP6,
		MlExpP6:        row.MlExpP6,
		ApeBkgP6:       row.ApeBkgP6,
		ApeRadiusP6:    row.ApeRadiusP6,
		ApePoisP6:      row.ApePoisP6,
		FlagSpSnr:      row.FlagSpSnr,
		FlagSpBps:      row.FlagSpBps,
		FlagSpScl:      row.FlagSpScl,
		FlagSpLga:      row.FlagSpLga,
		FlagSpGcCons:   row.FlagSpGcCons,
		FlagNoRadecErr: row.FlagNoRadecErr,
		FlagNoExtErr:   row.FlagNoExtErr,
		FlagNoCtsErr:   row.FlagNoCtsErr,
		FlagOpt:        row.FlagOpt,
	}
}

func (a Adapter) ConvertToMastercat(raw any, mapper *healpix.HEALPixMapper) (repository.Mastercat, error) {
	schema := raw.(repository.ErositaInputSchema)
	ipix := mapper.PixelAt(healpix.RADec(schema.RA, schema.DEC))
	return repository.Mastercat{
		ID:   schema.IAUNAME,
		Ipix: ipix,
		Ra:   schema.RA,
		Dec:  schema.DEC,
		Cat:  "erosita",
	}, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (repository.Metadata, error) {
	schema := raw.(repository.ErositaInputSchema)
	return repository.Erosita{
		ID:             schema.IAUNAME,
		Detuid:         repository.NullString{sql.NullString{String: schema.DETUID, Valid: true}},
		Skytile:        repository.NullInt64{sql.NullInt64{Int64: int64(schema.SKYTILE), Valid: true}},
		IDSrc:          repository.NullInt64{sql.NullInt64{Int64: int64(schema.ID_SRC), Valid: true}},
		Uid:            repository.NullInt64{sql.NullInt64{Int64: schema.UID, Valid: true}},
		UidHard:        repository.NullInt64{sql.NullInt64{Int64: schema.UID_HARD, Valid: true}},
		IDCluster:      repository.NullInt64{sql.NullInt64{Int64: int64(schema.ID_CLUSTER), Valid: true}},
		Ra:             repository.NullFloat64{sql.NullFloat64{Float64: schema.RA, Valid: true}},
		Dec:            repository.NullFloat64{sql.NullFloat64{Float64: schema.DEC, Valid: true}},
		RaLowerr:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.RA_LOWERR), Valid: true}},
		RaUperr:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.RA_UPERR), Valid: true}},
		DecLowerr:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DEC_LOWERR), Valid: true}},
		DecUperr:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DEC_UPERR), Valid: true}},
		PosErr:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.POS_ERR), Valid: true}},
		Mjd:            repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.MJD), Valid: true}},
		MjdMin:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.MJD_MIN), Valid: true}},
		MjdMax:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.MJD_MAX), Valid: true}},
		Ext:            repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.EXT), Valid: true}},
		ExtErr:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.EXT_ERR), Valid: true}},
		ExtLike:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.EXT_LIKE), Valid: true}},
		DetLike0:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_0), Valid: true}},
		MlCts1:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_1), Valid: true}},
		MlCtsErr1:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_1), Valid: true}},
		MlRate1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_1), Valid: true}},
		MlRateErr1:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_1), Valid: true}},
		MlFlux1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_1), Valid: true}},
		MlFluxErr1:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_1), Valid: true}},
		MlBkg1:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_1), Valid: true}},
		MlExp1:         repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_1), Valid: true}},
		ApeBkg1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_1), Valid: true}},
		ApeRadius1:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_1), Valid: true}},
		ApePois1:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_1), Valid: true}},
		DetLikeP1:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P1), Valid: true}},
		MlCtsP1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P1), Valid: true}},
		MlCtsErrP1:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P1), Valid: true}},
		MlRateP1:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P1), Valid: true}},
		MlRateErrP1:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P1), Valid: true}},
		MlFluxP1:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P1), Valid: true}},
		MlFluxErrP1:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P1), Valid: true}},
		MlBkgP1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P1), Valid: true}},
		MlExpP1:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P1), Valid: true}},
		ApeBkgP1:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P1), Valid: true}},
		ApeRadiusP1:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P1), Valid: true}},
		ApePoisP1:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P1), Valid: true}},
		DetLikeP2:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P2), Valid: true}},
		MlCtsP2:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P2), Valid: true}},
		MlCtsErrP2:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P2), Valid: true}},
		MlRateP2:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P2), Valid: true}},
		MlRateErrP2:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P2), Valid: true}},
		MlFluxP2:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P2), Valid: true}},
		MlFluxErrP2:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P2), Valid: true}},
		MlBkgP2:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P2), Valid: true}},
		MlExpP2:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P2), Valid: true}},
		ApeBkgP2:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P2), Valid: true}},
		ApeRadiusP2:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P2), Valid: true}},
		ApePoisP2:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P2), Valid: true}},
		DetLikeP3:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P3), Valid: true}},
		MlCtsP3:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P3), Valid: true}},
		MlCtsErrP3:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P3), Valid: true}},
		MlRateP3:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P3), Valid: true}},
		MlRateErrP3:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P3), Valid: true}},
		MlFluxP3:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P3), Valid: true}},
		MlFluxErrP3:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P3), Valid: true}},
		MlBkgP3:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P3), Valid: true}},
		MlExpP3:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P3), Valid: true}},
		ApeBkgP3:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P3), Valid: true}},
		ApeRadiusP3:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P3), Valid: true}},
		ApePoisP3:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P3), Valid: true}},
		DetLikeP4:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P4), Valid: true}},
		MlCtsP4:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P4), Valid: true}},
		MlCtsErrP4:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P4), Valid: true}},
		MlRateP4:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P4), Valid: true}},
		MlRateErrP4:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P4), Valid: true}},
		MlFluxP4:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P4), Valid: true}},
		MlFluxErrP4:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P4), Valid: true}},
		MlBkgP4:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P4), Valid: true}},
		MlExpP4:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P4), Valid: true}},
		ApeBkgP4:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P4), Valid: true}},
		ApeRadiusP4:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P4), Valid: true}},
		ApePoisP4:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P4), Valid: true}},
		DetLikeP5:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P5), Valid: true}},
		MlCtsP5:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P5), Valid: true}},
		MlCtsErrP5:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P5), Valid: true}},
		MlRateP5:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P5), Valid: true}},
		MlRateErrP5:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P5), Valid: true}},
		MlFluxP5:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P5), Valid: true}},
		MlFluxErrP5:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P5), Valid: true}},
		MlBkgP5:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P5), Valid: true}},
		MlExpP5:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P5), Valid: true}},
		ApeBkgP5:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P5), Valid: true}},
		ApeRadiusP5:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P5), Valid: true}},
		ApePoisP5:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P5), Valid: true}},
		DetLikeP6:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.DET_LIKE_P6), Valid: true}},
		MlCtsP6:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_P6), Valid: true}},
		MlCtsErrP6:     repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P6), Valid: true}},
		MlRateP6:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_P6), Valid: true}},
		MlRateErrP6:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P6), Valid: true}},
		MlFluxP6:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_P6), Valid: true}},
		MlFluxErrP6:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P6), Valid: true}},
		MlBkgP6:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_BKG_P6), Valid: true}},
		MlExpP6:        repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.ML_EXP_P6), Valid: true}},
		ApeBkgP6:       repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_BKG_P6), Valid: true}},
		ApeRadiusP6:    repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P6), Valid: true}},
		ApePoisP6:      repository.NullFloat64{sql.NullFloat64{Float64: float64(schema.APE_POIS_P6), Valid: true}},
		FlagSpSnr:      repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_SP_SNR), Valid: true}},
		FlagSpBps:      repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_SP_BPS), Valid: true}},
		FlagSpScl:      repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_SP_SCL), Valid: true}},
		FlagSpLga:      repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_SP_LGA), Valid: true}},
		FlagSpGcCons:   repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_SP_GC_CONS), Valid: true}},
		FlagNoRadecErr: repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_NO_RADEC_ERR), Valid: true}},
		FlagNoExtErr:   repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_NO_EXT_ERR), Valid: true}},
		FlagNoCtsErr:   repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_NO_CTS_ERR), Valid: true}},
		FlagOpt:        repository.NullInt64{sql.NullInt64{Int64: int64(schema.FLAG_OPT), Valid: true}},
	}, nil
}
