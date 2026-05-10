package repository

import (
	"context"
)

type ErositaInputSchema struct {
	IAUNAME           string  `parquet:"name=IAUNAME, type=BYTE_ARRAY"`
	DETUID            string  `parquet:"name=DETUID, type=BYTE_ARRAY"`
	SKYTILE           int32   `parquet:"name=SKYTILE, type=INT32"`
	ID_SRC            int32   `parquet:"name=ID_SRC, type=INT32"`
	UID               int64   `parquet:"name=UID, type=INT64"`
	UID_HARD          int64   `parquet:"name=UID_Hard, type=INT64"`
	ID_CLUSTER        int32   `parquet:"name=ID_CLUSTER, type=INT32"`
	RA                float64 `parquet:"name=RA, type=DOUBLE"`
	DEC               float64 `parquet:"name=DEC, type=DOUBLE"`
	RA_LOWERR         float32 `parquet:"name=RA_LOWERR, type=FLOAT"`
	RA_UPERR          float32 `parquet:"name=RA_UPERR, type=FLOAT"`
	DEC_LOWERR        float32 `parquet:"name=DEC_LOWERR, type=FLOAT"`
	DEC_UPERR         float32 `parquet:"name=DEC_UPERR, type=FLOAT"`
	POS_ERR           float32 `parquet:"name=POS_ERR, type=FLOAT"`
	MJD               float32 `parquet:"name=MJD, type=FLOAT"`
	MJD_MIN           float32 `parquet:"name=MJD_MIN, type=FLOAT"`
	MJD_MAX           float32 `parquet:"name=MJD_MAX, type=FLOAT"`
	EXT               float32 `parquet:"name=EXT, type=FLOAT"`
	EXT_ERR           float32 `parquet:"name=EXT_ERR, type=FLOAT"`
	EXT_LIKE          float32 `parquet:"name=EXT_LIKE, type=FLOAT"`
	DET_LIKE_0        float32 `parquet:"name=DET_LIKE_0, type=FLOAT"`
	ML_CTS_1          float32 `parquet:"name=ML_CTS_1, type=FLOAT"`
	ML_CTS_ERR_1      float32 `parquet:"name=ML_CTS_ERR_1, type=FLOAT"`
	ML_RATE_1         float32 `parquet:"name=ML_RATE_1, type=FLOAT"`
	ML_RATE_ERR_1     float32 `parquet:"name=ML_RATE_ERR_1, type=FLOAT"`
	ML_FLUX_1         float32 `parquet:"name=ML_FLUX_1, type=FLOAT"`
	ML_FLUX_ERR_1     float32 `parquet:"name=ML_FLUX_ERR_1, type=FLOAT"`
	ML_BKG_1          float32 `parquet:"name=ML_BKG_1, type=FLOAT"`
	ML_EXP_1          float32 `parquet:"name=ML_EXP_1, type=FLOAT"`
	APE_BKG_1         float32 `parquet:"name=APE_BKG_1, type=FLOAT"`
	APE_RADIUS_1      float32 `parquet:"name=APE_RADIUS_1, type=FLOAT"`
	APE_POIS_1        float32 `parquet:"name=APE_POIS_1, type=FLOAT"`
	DET_LIKE_P1       float32 `parquet:"name=DET_LIKE_P1, type=FLOAT"`
	ML_CTS_P1         float32 `parquet:"name=ML_CTS_P1, type=FLOAT"`
	ML_CTS_ERR_P1     float32 `parquet:"name=ML_CTS_ERR_P1, type=FLOAT"`
	ML_RATE_P1        float32 `parquet:"name=ML_RATE_P1, type=FLOAT"`
	ML_RATE_ERR_P1    float32 `parquet:"name=ML_RATE_ERR_P1, type=FLOAT"`
	ML_FLUX_P1        float32 `parquet:"name=ML_FLUX_P1, type=FLOAT"`
	ML_FLUX_ERR_P1    float32 `parquet:"name=ML_FLUX_ERR_P1, type=FLOAT"`
	ML_BKG_P1         float32 `parquet:"name=ML_BKG_P1, type=FLOAT"`
	ML_EXP_P1         float32 `parquet:"name=ML_EXP_P1, type=FLOAT"`
	APE_BKG_P1        float32 `parquet:"name=APE_BKG_P1, type=FLOAT"`
	APE_RADIUS_P1     float32 `parquet:"name=APE_RADIUS_P1, type=FLOAT"`
	APE_POIS_P1       float32 `parquet:"name=APE_POIS_P1, type=FLOAT"`
	DET_LIKE_P2       float32 `parquet:"name=DET_LIKE_P2, type=FLOAT"`
	ML_CTS_P2         float32 `parquet:"name=ML_CTS_P2, type=FLOAT"`
	ML_CTS_ERR_P2     float32 `parquet:"name=ML_CTS_ERR_P2, type=FLOAT"`
	ML_RATE_P2        float32 `parquet:"name=ML_RATE_P2, type=FLOAT"`
	ML_RATE_ERR_P2    float32 `parquet:"name=ML_RATE_ERR_P2, type=FLOAT"`
	ML_FLUX_P2        float32 `parquet:"name=ML_FLUX_P2, type=FLOAT"`
	ML_FLUX_ERR_P2    float32 `parquet:"name=ML_FLUX_ERR_P2, type=FLOAT"`
	ML_BKG_P2         float32 `parquet:"name=ML_BKG_P2, type=FLOAT"`
	ML_EXP_P2         float32 `parquet:"name=ML_EXP_P2, type=FLOAT"`
	APE_BKG_P2        float32 `parquet:"name=APE_BKG_P2, type=FLOAT"`
	APE_RADIUS_P2     float32 `parquet:"name=APE_RADIUS_P2, type=FLOAT"`
	APE_POIS_P2       float32 `parquet:"name=APE_POIS_P2, type=FLOAT"`
	DET_LIKE_P3       float32 `parquet:"name=DET_LIKE_P3, type=FLOAT"`
	ML_CTS_P3         float32 `parquet:"name=ML_CTS_P3, type=FLOAT"`
	ML_CTS_ERR_P3     float32 `parquet:"name=ML_CTS_ERR_P3, type=FLOAT"`
	ML_RATE_P3        float32 `parquet:"name=ML_RATE_P3, type=FLOAT"`
	ML_RATE_ERR_P3    float32 `parquet:"name=ML_RATE_ERR_P3, type=FLOAT"`
	ML_FLUX_P3        float32 `parquet:"name=ML_FLUX_P3, type=FLOAT"`
	ML_FLUX_ERR_P3    float32 `parquet:"name=ML_FLUX_ERR_P3, type=FLOAT"`
	ML_BKG_P3         float32 `parquet:"name=ML_BKG_P3, type=FLOAT"`
	ML_EXP_P3         float32 `parquet:"name=ML_EXP_P3, type=FLOAT"`
	APE_BKG_P3        float32 `parquet:"name=APE_BKG_P3, type=FLOAT"`
	APE_RADIUS_P3     float32 `parquet:"name=APE_RADIUS_P3, type=FLOAT"`
	APE_POIS_P3       float32 `parquet:"name=APE_POIS_P3, type=FLOAT"`
	DET_LIKE_P4       float32 `parquet:"name=DET_LIKE_P4, type=FLOAT"`
	ML_CTS_P4         float32 `parquet:"name=ML_CTS_P4, type=FLOAT"`
	ML_CTS_ERR_P4     float32 `parquet:"name=ML_CTS_ERR_P4, type=FLOAT"`
	ML_RATE_P4        float32 `parquet:"name=ML_RATE_P4, type=FLOAT"`
	ML_RATE_ERR_P4    float32 `parquet:"name=ML_RATE_ERR_P4, type=FLOAT"`
	ML_FLUX_P4        float32 `parquet:"name=ML_FLUX_P4, type=FLOAT"`
	ML_FLUX_ERR_P4    float32 `parquet:"name=ML_FLUX_ERR_P4, type=FLOAT"`
	ML_BKG_P4         float32 `parquet:"name=ML_BKG_P4, type=FLOAT"`
	ML_EXP_P4         float32 `parquet:"name=ML_EXP_P4, type=FLOAT"`
	APE_BKG_P4        float32 `parquet:"name=APE_BKG_P4, type=FLOAT"`
	APE_RADIUS_P4     float32 `parquet:"name=APE_RADIUS_P4, type=FLOAT"`
	APE_POIS_P4       float32 `parquet:"name=APE_POIS_P4, type=FLOAT"`
	DET_LIKE_P5       float32 `parquet:"name=DET_LIKE_P5, type=FLOAT"`
	ML_CTS_P5         float32 `parquet:"name=ML_CTS_P5, type=FLOAT"`
	ML_CTS_ERR_P5     float32 `parquet:"name=ML_CTS_ERR_P5, type=FLOAT"`
	ML_RATE_P5        float32 `parquet:"name=ML_RATE_P5, type=FLOAT"`
	ML_RATE_ERR_P5    float32 `parquet:"name=ML_RATE_ERR_P5, type=FLOAT"`
	ML_FLUX_P5        float32 `parquet:"name=ML_FLUX_P5, type=FLOAT"`
	ML_FLUX_ERR_P5    float32 `parquet:"name=ML_FLUX_ERR_P5, type=FLOAT"`
	ML_BKG_P5         float32 `parquet:"name=ML_BKG_P5, type=FLOAT"`
	ML_EXP_P5         float32 `parquet:"name=ML_EXP_P5, type=FLOAT"`
	APE_BKG_P5        float32 `parquet:"name=APE_BKG_P5, type=FLOAT"`
	APE_RADIUS_P5     float32 `parquet:"name=APE_RADIUS_P5, type=FLOAT"`
	APE_POIS_P5       float32 `parquet:"name=APE_POIS_P5, type=FLOAT"`
	DET_LIKE_P6       float32 `parquet:"name=DET_LIKE_P6, type=FLOAT"`
	ML_CTS_P6         float32 `parquet:"name=ML_CTS_P6, type=FLOAT"`
	ML_CTS_ERR_P6     float32 `parquet:"name=ML_CTS_ERR_P6, type=FLOAT"`
	ML_RATE_P6        float32 `parquet:"name=ML_RATE_P6, type=FLOAT"`
	ML_RATE_ERR_P6    float32 `parquet:"name=ML_RATE_ERR_P6, type=FLOAT"`
	ML_FLUX_P6        float32 `parquet:"name=ML_FLUX_P6, type=FLOAT"`
	ML_FLUX_ERR_P6    float32 `parquet:"name=ML_FLUX_ERR_P6, type=FLOAT"`
	ML_BKG_P6         float32 `parquet:"name=ML_BKG_P6, type=FLOAT"`
	ML_EXP_P6         float32 `parquet:"name=ML_EXP_P6, type=FLOAT"`
	APE_BKG_P6        float32 `parquet:"name=APE_BKG_P6, type=FLOAT"`
	APE_RADIUS_P6     float32 `parquet:"name=APE_RADIUS_P6, type=FLOAT"`
	APE_POIS_P6       float32 `parquet:"name=APE_POIS_P6, type=FLOAT"`
	FLAG_SP_SNR       int32   `parquet:"name=FLAG_SP_SNR, type=INT32"`
	FLAG_SP_BPS       int32   `parquet:"name=FLAG_SP_BPS, type=INT32"`
	FLAG_SP_SCL       int32   `parquet:"name=FLAG_SP_SCL, type=INT32"`
	FLAG_SP_LGA       int32   `parquet:"name=FLAG_SP_LGA, type=INT32"`
	FLAG_SP_GC_CONS   int32   `parquet:"name=FLAG_SP_GC_CONS, type=INT32"`
	FLAG_NO_RADEC_ERR int32   `parquet:"name=FLAG_NO_RADEC_ERR, type=INT32"`
	FLAG_NO_EXT_ERR   int32   `parquet:"name=FLAG_NO_EXT_ERR, type=INT32"`
	FLAG_NO_CTS_ERR   int32   `parquet:"name=FLAG_NO_CTS_ERR, type=INT32"`
	FLAG_OPT          int32   `parquet:"name=FLAG_OPT, type=INT32"`
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
	ra := 0.0
	dec := 0.0
	if m.Ra.Valid {
		ra = m.Ra.Float64
	}
	if m.Dec.Valid {
		dec = m.Dec.Float64
	}
	return ra, dec
}

func (m GetErositaFromPixelsRow) GetCatalog() string {
	return "eROSITA"
}

func (q *Queries) InsertErositaWithoutParams(ctx context.Context, arg Erosita) error {
	return q.InsertErosita(ctx, InsertErositaParams{
		ID:             arg.ID,
		Detuid:         arg.Detuid,
		Skytile:        arg.Skytile,
		IDSrc:          arg.IDSrc,
		Uid:            arg.Uid,
		UidHard:        arg.UidHard,
		IDCluster:      arg.IDCluster,
		Ra:             arg.Ra,
		Dec:            arg.Dec,
		RaLowerr:       arg.RaLowerr,
		RaUperr:        arg.RaUperr,
		DecLowerr:      arg.DecLowerr,
		DecUperr:       arg.DecUperr,
		PosErr:         arg.PosErr,
		Mjd:            arg.Mjd,
		MjdMin:         arg.MjdMin,
		MjdMax:         arg.MjdMax,
		Ext:            arg.Ext,
		ExtErr:         arg.ExtErr,
		ExtLike:        arg.ExtLike,
		DetLike0:       arg.DetLike0,
		MlCts1:         arg.MlCts1,
		MlCtsErr1:      arg.MlCtsErr1,
		MlRate1:        arg.MlRate1,
		MlRateErr1:     arg.MlRateErr1,
		MlFlux1:        arg.MlFlux1,
		MlFluxErr1:     arg.MlFluxErr1,
		MlBkg1:         arg.MlBkg1,
		MlExp1:         arg.MlExp1,
		ApeBkg1:        arg.ApeBkg1,
		ApeRadius1:     arg.ApeRadius1,
		ApePois1:       arg.ApePois1,
		DetLikeP1:      arg.DetLikeP1,
		MlCtsP1:        arg.MlCtsP1,
		MlCtsErrP1:     arg.MlCtsErrP1,
		MlRateP1:       arg.MlRateP1,
		MlRateErrP1:    arg.MlRateErrP1,
		MlFluxP1:       arg.MlFluxP1,
		MlFluxErrP1:    arg.MlFluxErrP1,
		MlBkgP1:        arg.MlBkgP1,
		MlExpP1:        arg.MlExpP1,
		ApeBkgP1:       arg.ApeBkgP1,
		ApeRadiusP1:    arg.ApeRadiusP1,
		ApePoisP1:      arg.ApePoisP1,
		DetLikeP2:      arg.DetLikeP2,
		MlCtsP2:        arg.MlCtsP2,
		MlCtsErrP2:     arg.MlCtsErrP2,
		MlRateP2:       arg.MlRateP2,
		MlRateErrP2:    arg.MlRateErrP2,
		MlFluxP2:       arg.MlFluxP2,
		MlFluxErrP2:    arg.MlFluxErrP2,
		MlBkgP2:        arg.MlBkgP2,
		MlExpP2:        arg.MlExpP2,
		ApeBkgP2:       arg.ApeBkgP2,
		ApeRadiusP2:    arg.ApeRadiusP2,
		ApePoisP2:      arg.ApePoisP2,
		DetLikeP3:      arg.DetLikeP3,
		MlCtsP3:        arg.MlCtsP3,
		MlCtsErrP3:     arg.MlCtsErrP3,
		MlRateP3:       arg.MlRateP3,
		MlRateErrP3:    arg.MlRateErrP3,
		MlFluxP3:       arg.MlFluxP3,
		MlFluxErrP3:    arg.MlFluxErrP3,
		MlBkgP3:        arg.MlBkgP3,
		MlExpP3:        arg.MlExpP3,
		ApeBkgP3:       arg.ApeBkgP3,
		ApeRadiusP3:    arg.ApeRadiusP3,
		ApePoisP3:      arg.ApePoisP3,
		DetLikeP4:      arg.DetLikeP4,
		MlCtsP4:        arg.MlCtsP4,
		MlCtsErrP4:     arg.MlCtsErrP4,
		MlRateP4:       arg.MlRateP4,
		MlRateErrP4:    arg.MlRateErrP4,
		MlFluxP4:       arg.MlFluxP4,
		MlFluxErrP4:    arg.MlFluxErrP4,
		MlBkgP4:        arg.MlBkgP4,
		MlExpP4:        arg.MlExpP4,
		ApeBkgP4:       arg.ApeBkgP4,
		ApeRadiusP4:    arg.ApeRadiusP4,
		ApePoisP4:      arg.ApePoisP4,
		DetLikeP5:      arg.DetLikeP5,
		MlCtsP5:        arg.MlCtsP5,
		MlCtsErrP5:     arg.MlCtsErrP5,
		MlRateP5:       arg.MlRateP5,
		MlRateErrP5:    arg.MlRateErrP5,
		MlFluxP5:       arg.MlFluxP5,
		MlFluxErrP5:    arg.MlFluxErrP5,
		MlBkgP5:        arg.MlBkgP5,
		MlExpP5:        arg.MlExpP5,
		ApeBkgP5:       arg.ApeBkgP5,
		ApeRadiusP5:    arg.ApeRadiusP5,
		ApePoisP5:      arg.ApePoisP5,
		DetLikeP6:      arg.DetLikeP6,
		MlCtsP6:        arg.MlCtsP6,
		MlCtsErrP6:     arg.MlCtsErrP6,
		MlRateP6:       arg.MlRateP6,
		MlRateErrP6:    arg.MlRateErrP6,
		MlFluxP6:       arg.MlFluxP6,
		MlFluxErrP6:    arg.MlFluxErrP6,
		MlBkgP6:        arg.MlBkgP6,
		MlExpP6:        arg.MlExpP6,
		ApeBkgP6:       arg.ApeBkgP6,
		ApeRadiusP6:    arg.ApeRadiusP6,
		ApePoisP6:      arg.ApePoisP6,
		FlagSpSnr:      arg.FlagSpSnr,
		FlagSpBps:      arg.FlagSpBps,
		FlagSpScl:      arg.FlagSpScl,
		FlagSpLga:      arg.FlagSpLga,
		FlagSpGcCons:   arg.FlagSpGcCons,
		FlagNoRadecErr: arg.FlagNoRadecErr,
		FlagNoExtErr:   arg.FlagNoExtErr,
		FlagNoCtsErr:   arg.FlagNoCtsErr,
		FlagOpt:        arg.FlagOpt,
	})
}
