package erosita

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/catalog"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

const displayName = "eROSITA"

type Adapter struct {
	repo *repository.Queries
}

func init() {
	catalog.Register("erosita", func(repo *repository.Queries) (catalog.CatalogAdapter, error) {
		return &Adapter{repo: repo}, nil
	})
}

func (a Adapter) Name() string {
	return "erosita"
}

func (a Adapter) NewRawRecord() any {
	return InputSchema{}
}

func (a Adapter) NewMetadataRecord() any {
	return repository.Erosita{}
}

func (a Adapter) BulkInsertMetadata(ctx context.Context, rows []any) error {
	if a.repo == nil {
		return fmt.Errorf("erosita adapter has no repository")
	}
	return a.repo.BulkInsertErosita(ctx, rows)
}

func (a Adapter) GetByID(ctx context.Context, id string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("erosita adapter has no repository")
	}
	return a.repo.GetErosita(ctx, id)
}

func (a Adapter) BulkGetByID(ctx context.Context, ids []string) (any, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("erosita adapter has no repository")
	}
	return a.repo.BulkGetErosita(ctx, ids)
}

func (a Adapter) GetFromPixels(ctx context.Context, pixels []int64) ([]repository.Metadata, error) {
	if a.repo == nil {
		return nil, fmt.Errorf("erosita adapter has no repository")
	}
	rows, err := a.repo.GetErositaFromPixels(ctx, pixels)
	if err != nil {
		return nil, err
	}
	result := make([]repository.Metadata, len(rows))
	for i, r := range rows {
		result[i] = repository.Metadata{
			ID:      r.ID,
			Catalog: displayName,
			Ra:      r.Ra.Float64,
			Dec:     r.Dec.Float64,
			Object: repository.Erosita{
				ID:             r.ID,
				Detuid:         r.Detuid,
				Skytile:        r.Skytile,
				IDSrc:          r.IDSrc,
				Uid:            r.Uid,
				UidHard:        r.UidHard,
				IDCluster:      r.IDCluster,
				Ra:             r.Ra,
				Dec:            r.Dec,
				RaLowerr:       r.RaLowerr,
				RaUperr:        r.RaUperr,
				DecLowerr:      r.DecLowerr,
				DecUperr:       r.DecUperr,
				PosErr:         r.PosErr,
				Mjd:            r.Mjd,
				MjdMin:         r.MjdMin,
				MjdMax:         r.MjdMax,
				Ext:            r.Ext,
				ExtErr:         r.ExtErr,
				ExtLike:        r.ExtLike,
				DetLike0:       r.DetLike0,
				MlCts1:         r.MlCts1,
				MlCtsErr1:      r.MlCtsErr1,
				MlRate1:        r.MlRate1,
				MlRateErr1:     r.MlRateErr1,
				MlFlux1:        r.MlFlux1,
				MlFluxErr1:     r.MlFluxErr1,
				MlBkg1:         r.MlBkg1,
				MlExp1:         r.MlExp1,
				ApeBkg1:        r.ApeBkg1,
				ApeRadius1:     r.ApeRadius1,
				ApePois1:       r.ApePois1,
				DetLikeP1:      r.DetLikeP1,
				MlCtsP1:        r.MlCtsP1,
				MlCtsErrP1:     r.MlCtsErrP1,
				MlRateP1:       r.MlRateP1,
				MlRateErrP1:    r.MlRateErrP1,
				MlFluxP1:       r.MlFluxP1,
				MlFluxErrP1:    r.MlFluxErrP1,
				MlBkgP1:        r.MlBkgP1,
				MlExpP1:        r.MlExpP1,
				ApeBkgP1:       r.ApeBkgP1,
				ApeRadiusP1:    r.ApeRadiusP1,
				ApePoisP1:      r.ApePoisP1,
				DetLikeP2:      r.DetLikeP2,
				MlCtsP2:        r.MlCtsP2,
				MlCtsErrP2:     r.MlCtsErrP2,
				MlRateP2:       r.MlRateP2,
				MlRateErrP2:    r.MlRateErrP2,
				MlFluxP2:       r.MlFluxP2,
				MlFluxErrP2:    r.MlFluxErrP2,
				MlBkgP2:        r.MlBkgP2,
				MlExpP2:        r.MlExpP2,
				ApeBkgP2:       r.ApeBkgP2,
				ApeRadiusP2:    r.ApeRadiusP2,
				ApePoisP2:      r.ApePoisP2,
				DetLikeP3:      r.DetLikeP3,
				MlCtsP3:        r.MlCtsP3,
				MlCtsErrP3:     r.MlCtsErrP3,
				MlRateP3:       r.MlRateP3,
				MlRateErrP3:    r.MlRateErrP3,
				MlFluxP3:       r.MlFluxP3,
				MlFluxErrP3:    r.MlFluxErrP3,
				MlBkgP3:        r.MlBkgP3,
				MlExpP3:        r.MlExpP3,
				ApeBkgP3:       r.ApeBkgP3,
				ApeRadiusP3:    r.ApeRadiusP3,
				ApePoisP3:      r.ApePoisP3,
				DetLikeP4:      r.DetLikeP4,
				MlCtsP4:        r.MlCtsP4,
				MlCtsErrP4:     r.MlCtsErrP4,
				MlRateP4:       r.MlRateP4,
				MlRateErrP4:    r.MlRateErrP4,
				MlFluxP4:       r.MlFluxP4,
				MlFluxErrP4:    r.MlFluxErrP4,
				MlBkgP4:        r.MlBkgP4,
				MlExpP4:        r.MlExpP4,
				ApeBkgP4:       r.ApeBkgP4,
				ApeRadiusP4:    r.ApeRadiusP4,
				ApePoisP4:      r.ApePoisP4,
				DetLikeP5:      r.DetLikeP5,
				MlCtsP5:        r.MlCtsP5,
				MlCtsErrP5:     r.MlCtsErrP5,
				MlRateP5:       r.MlRateP5,
				MlRateErrP5:    r.MlRateErrP5,
				MlFluxP5:       r.MlFluxP5,
				MlFluxErrP5:    r.MlFluxErrP5,
				MlBkgP5:        r.MlBkgP5,
				MlExpP5:        r.MlExpP5,
				ApeBkgP5:       r.ApeBkgP5,
				ApeRadiusP5:    r.ApeRadiusP5,
				ApePoisP5:      r.ApePoisP5,
				DetLikeP6:      r.DetLikeP6,
				MlCtsP6:        r.MlCtsP6,
				MlCtsErrP6:     r.MlCtsErrP6,
				MlRateP6:       r.MlRateP6,
				MlRateErrP6:    r.MlRateErrP6,
				MlFluxP6:       r.MlFluxP6,
				MlFluxErrP6:    r.MlFluxErrP6,
				MlBkgP6:        r.MlBkgP6,
				MlExpP6:        r.MlExpP6,
				ApeBkgP6:       r.ApeBkgP6,
				ApeRadiusP6:    r.ApeRadiusP6,
				ApePoisP6:      r.ApePoisP6,
				FlagSpSnr:      r.FlagSpSnr,
				FlagSpBps:      r.FlagSpBps,
				FlagSpScl:      r.FlagSpScl,
				FlagSpLga:      r.FlagSpLga,
				FlagSpGcCons:   r.FlagSpGcCons,
				FlagNoRadecErr: r.FlagNoRadecErr,
				FlagNoExtErr:   r.FlagNoExtErr,
				FlagNoCtsErr:   r.FlagNoCtsErr,
				FlagOpt:        r.FlagOpt,
			},
		}
	}
	return result, nil
}

func (a Adapter) GetCoordinates(raw any) (float64, float64, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return 0, 0, fmt.Errorf("expected erosita.InputSchema, got %T", raw)
	}
	return schema.RA, schema.DEC, nil
}

func (a Adapter) ConvertToMastercat(raw any, ipix int64) (repository.Mastercat, error) {
	schema, ok := raw.(InputSchema)
	if !ok {
		return repository.Mastercat{}, fmt.Errorf("expected erosita.InputSchema, got %T", raw)
	}
	return repository.Mastercat{
		ID:   schema.IAUNAME,
		Ipix: ipix,
		Ra:   schema.RA,
		Dec:  schema.DEC,
		Cat:  "erosita",
	}, nil
}

func (a Adapter) ConvertToMetadataFromRaw(raw any) (any, error) {
	schema := raw.(InputSchema)
	return repository.Erosita{
		ID:             schema.IAUNAME,
		Detuid:         repository.NullString{NullString: sql.NullString{String: schema.DETUID, Valid: true}},
		Skytile:        repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.SKYTILE), Valid: true}},
		IDSrc:          repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.ID_SRC), Valid: true}},
		Uid:            repository.NullInt64{NullInt64: sql.NullInt64{Int64: schema.UID, Valid: true}},
		UidHard:        repository.NullInt64{NullInt64: sql.NullInt64{Int64: schema.UID_HARD, Valid: true}},
		IDCluster:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.ID_CLUSTER), Valid: true}},
		Ra:             repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: schema.RA, Valid: true}},
		Dec:            repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: schema.DEC, Valid: true}},
		RaLowerr:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.RA_LOWERR), Valid: true}},
		RaUperr:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.RA_UPERR), Valid: true}},
		DecLowerr:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DEC_LOWERR), Valid: true}},
		DecUperr:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DEC_UPERR), Valid: true}},
		PosErr:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.POS_ERR), Valid: true}},
		Mjd:            repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.MJD), Valid: true}},
		MjdMin:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.MJD_MIN), Valid: true}},
		MjdMax:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.MJD_MAX), Valid: true}},
		Ext:            repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.EXT), Valid: true}},
		ExtErr:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.EXT_ERR), Valid: true}},
		ExtLike:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.EXT_LIKE), Valid: true}},
		DetLike0:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_0), Valid: true}},
		MlCts1:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_1), Valid: true}},
		MlCtsErr1:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_1), Valid: true}},
		MlRate1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_1), Valid: true}},
		MlRateErr1:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_1), Valid: true}},
		MlFlux1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_1), Valid: true}},
		MlFluxErr1:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_1), Valid: true}},
		MlBkg1:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_1), Valid: true}},
		MlExp1:         repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_1), Valid: true}},
		ApeBkg1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_1), Valid: true}},
		ApeRadius1:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_1), Valid: true}},
		ApePois1:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_1), Valid: true}},
		DetLikeP1:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P1), Valid: true}},
		MlCtsP1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P1), Valid: true}},
		MlCtsErrP1:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P1), Valid: true}},
		MlRateP1:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P1), Valid: true}},
		MlRateErrP1:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P1), Valid: true}},
		MlFluxP1:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P1), Valid: true}},
		MlFluxErrP1:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P1), Valid: true}},
		MlBkgP1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P1), Valid: true}},
		MlExpP1:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P1), Valid: true}},
		ApeBkgP1:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P1), Valid: true}},
		ApeRadiusP1:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P1), Valid: true}},
		ApePoisP1:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P1), Valid: true}},
		DetLikeP2:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P2), Valid: true}},
		MlCtsP2:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P2), Valid: true}},
		MlCtsErrP2:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P2), Valid: true}},
		MlRateP2:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P2), Valid: true}},
		MlRateErrP2:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P2), Valid: true}},
		MlFluxP2:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P2), Valid: true}},
		MlFluxErrP2:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P2), Valid: true}},
		MlBkgP2:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P2), Valid: true}},
		MlExpP2:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P2), Valid: true}},
		ApeBkgP2:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P2), Valid: true}},
		ApeRadiusP2:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P2), Valid: true}},
		ApePoisP2:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P2), Valid: true}},
		DetLikeP3:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P3), Valid: true}},
		MlCtsP3:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P3), Valid: true}},
		MlCtsErrP3:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P3), Valid: true}},
		MlRateP3:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P3), Valid: true}},
		MlRateErrP3:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P3), Valid: true}},
		MlFluxP3:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P3), Valid: true}},
		MlFluxErrP3:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P3), Valid: true}},
		MlBkgP3:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P3), Valid: true}},
		MlExpP3:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P3), Valid: true}},
		ApeBkgP3:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P3), Valid: true}},
		ApeRadiusP3:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P3), Valid: true}},
		ApePoisP3:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P3), Valid: true}},
		DetLikeP4:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P4), Valid: true}},
		MlCtsP4:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P4), Valid: true}},
		MlCtsErrP4:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P4), Valid: true}},
		MlRateP4:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P4), Valid: true}},
		MlRateErrP4:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P4), Valid: true}},
		MlFluxP4:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P4), Valid: true}},
		MlFluxErrP4:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P4), Valid: true}},
		MlBkgP4:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P4), Valid: true}},
		MlExpP4:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P4), Valid: true}},
		ApeBkgP4:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P4), Valid: true}},
		ApeRadiusP4:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P4), Valid: true}},
		ApePoisP4:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P4), Valid: true}},
		DetLikeP5:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P5), Valid: true}},
		MlCtsP5:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P5), Valid: true}},
		MlCtsErrP5:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P5), Valid: true}},
		MlRateP5:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P5), Valid: true}},
		MlRateErrP5:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P5), Valid: true}},
		MlFluxP5:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P5), Valid: true}},
		MlFluxErrP5:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P5), Valid: true}},
		MlBkgP5:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P5), Valid: true}},
		MlExpP5:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P5), Valid: true}},
		ApeBkgP5:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P5), Valid: true}},
		ApeRadiusP5:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P5), Valid: true}},
		ApePoisP5:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P5), Valid: true}},
		DetLikeP6:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.DET_LIKE_P6), Valid: true}},
		MlCtsP6:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_P6), Valid: true}},
		MlCtsErrP6:     repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_CTS_ERR_P6), Valid: true}},
		MlRateP6:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_P6), Valid: true}},
		MlRateErrP6:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_RATE_ERR_P6), Valid: true}},
		MlFluxP6:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_P6), Valid: true}},
		MlFluxErrP6:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_FLUX_ERR_P6), Valid: true}},
		MlBkgP6:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_BKG_P6), Valid: true}},
		MlExpP6:        repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.ML_EXP_P6), Valid: true}},
		ApeBkgP6:       repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_BKG_P6), Valid: true}},
		ApeRadiusP6:    repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_RADIUS_P6), Valid: true}},
		ApePoisP6:      repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: float64(schema.APE_POIS_P6), Valid: true}},
		FlagSpSnr:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_SP_SNR), Valid: true}},
		FlagSpBps:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_SP_BPS), Valid: true}},
		FlagSpScl:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_SP_SCL), Valid: true}},
		FlagSpLga:      repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_SP_LGA), Valid: true}},
		FlagSpGcCons:   repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_SP_GC_CONS), Valid: true}},
		FlagNoRadecErr: repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_NO_RADEC_ERR), Valid: true}},
		FlagNoExtErr:   repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_NO_EXT_ERR), Valid: true}},
		FlagNoCtsErr:   repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_NO_CTS_ERR), Valid: true}},
		FlagOpt:        repository.NullInt64{NullInt64: sql.NullInt64{Int64: int64(schema.FLAG_OPT), Valid: true}},
	}, nil
}
