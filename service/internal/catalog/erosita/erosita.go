package erosita

import (
	"database/sql"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type InputSchema struct {
	IAUNAME           string
	DETUID            string
	SKYTILE           int32
	ID_SRC            int32
	UID               int64
	UID_HARD          int64
	ID_CLUSTER        int32
	RA                float64
	DEC               float64
	RA_LOWERR         float32
	RA_UPERR          float32
	DEC_LOWERR        float32
	DEC_UPERR         float32
	POS_ERR           float32
	MJD               float32
	MJD_MIN           float32
	MJD_MAX           float32
	EXT               float32
	EXT_ERR           float32
	EXT_LIKE          float32
	DET_LIKE_0        float32
	ML_CTS_1          float32
	ML_CTS_ERR_1      float32
	ML_RATE_1         float32
	ML_RATE_ERR_1     float32
	ML_FLUX_1         float32
	ML_FLUX_ERR_1     float32
	ML_BKG_1          float32
	ML_EXP_1          float32
	APE_BKG_1         float32
	APE_RADIUS_1      float32
	APE_POIS_1        float32
	DET_LIKE_P1       float32
	ML_CTS_P1         float32
	ML_CTS_ERR_P1     float32
	ML_RATE_P1        float32
	ML_RATE_ERR_P1    float32
	ML_FLUX_P1        float32
	ML_FLUX_ERR_P1    float32
	ML_BKG_P1         float32
	ML_EXP_P1         float32
	APE_BKG_P1        float32
	APE_RADIUS_P1     float32
	APE_POIS_P1       float32
	DET_LIKE_P2       float32
	ML_CTS_P2         float32
	ML_CTS_ERR_P2     float32
	ML_RATE_P2        float32
	ML_RATE_ERR_P2    float32
	ML_FLUX_P2        float32
	ML_FLUX_ERR_P2    float32
	ML_BKG_P2         float32
	ML_EXP_P2         float32
	APE_BKG_P2        float32
	APE_RADIUS_P2     float32
	APE_POIS_P2       float32
	DET_LIKE_P3       float32
	ML_CTS_P3         float32
	ML_CTS_ERR_P3     float32
	ML_RATE_P3        float32
	ML_RATE_ERR_P3    float32
	ML_FLUX_P3        float32
	ML_FLUX_ERR_P3    float32
	ML_BKG_P3         float32
	ML_EXP_P3         float32
	APE_BKG_P3        float32
	APE_RADIUS_P3     float32
	APE_POIS_P3       float32
	DET_LIKE_P4       float32
	ML_CTS_P4         float32
	ML_CTS_ERR_P4     float32
	ML_RATE_P4        float32
	ML_RATE_ERR_P4    float32
	ML_FLUX_P4        float32
	ML_FLUX_ERR_P4    float32
	ML_BKG_P4         float32
	ML_EXP_P4         float32
	APE_BKG_P4        float32
	APE_RADIUS_P4     float32
	APE_POIS_P4       float32
	DET_LIKE_P5       float32
	ML_CTS_P5         float32
	ML_CTS_ERR_P5     float32
	ML_RATE_P5        float32
	ML_RATE_ERR_P5    float32
	ML_FLUX_P5        float32
	ML_FLUX_ERR_P5    float32
	ML_BKG_P5         float32
	ML_EXP_P5         float32
	APE_BKG_P5        float32
	APE_RADIUS_P5     float32
	APE_POIS_P5       float32
	DET_LIKE_P6       float32
	ML_CTS_P6         float32
	ML_CTS_ERR_P6     float32
	ML_RATE_P6        float32
	ML_RATE_ERR_P6    float32
	ML_FLUX_P6        float32
	ML_FLUX_ERR_P6    float32
	ML_BKG_P6         float32
	ML_EXP_P6         float32
	APE_BKG_P6        float32
	APE_RADIUS_P6     float32
	APE_POIS_P6       float32
	FLAG_SP_SNR       int32
	FLAG_SP_BPS       int32
	FLAG_SP_SCL       int32
	FLAG_SP_LGA       int32
	FLAG_SP_GC_CONS   int32
	FLAG_NO_RADEC_ERR int32
	FLAG_NO_EXT_ERR   int32
	FLAG_NO_CTS_ERR   int32
	FLAG_OPT          int32
}

func (r InputSchema) GetCoordinates() (float64, float64) {
	return r.RA, r.DEC
}

func (r InputSchema) MastercatID() string {
	return r.IAUNAME
}

func (schema InputSchema) ConvertToMetadata() (repository.Erosita, error) {
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
