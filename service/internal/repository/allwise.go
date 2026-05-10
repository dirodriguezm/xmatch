package repository

import (
	"context"
)

type AllwiseInputSchema struct {
	Source_id    *string  `parquet:"name=source_id, type=BYTE_ARRAY"`
	Cntr         *int64   `parquet:"name=cntr, type=INT64"`
	Ra           *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec          *float64 `parquet:"name=dec, type=DOUBLE"`
	W1mpro       *float64 `parquet:"name=w1mpro, type=DOUBLE"`
	W1sigmpro    *float64 `parquet:"name=w1sigmpro, type=DOUBLE"`
	W2mpro       *float64 `parquet:"name=w2mpro, type=DOUBLE"`
	W2sigmpro    *float64 `parquet:"name=w2sigmpro, type=DOUBLE"`
	W3mpro       *float64 `parquet:"name=w3mpro, type=DOUBLE"`
	W3sigmpro    *float64 `parquet:"name=w3sigmpro, type=DOUBLE"`
	W4mpro       *float64 `parquet:"name=w4mpro, type=DOUBLE"`
	W4sigmpro    *float64 `parquet:"name=w4sigmpro, type=DOUBLE"`
	J_m_2mass    *float64 `parquet:"name=j_m_2mass, type=DOUBLE"`
	H_m_2mass    *float64 `parquet:"name=h_m_2mass, type=DOUBLE"`
	K_m_2mass    *float64 `parquet:"name=k_m_2mass, type=DOUBLE"`
	J_msig_2mass *float64 `parquet:"name=j_msig_2mass, type=DOUBLE"`
	H_msig_2mass *float64 `parquet:"name=h_msig_2mass, type=DOUBLE"`
	K_msig_2mass *float64 `parquet:"name=k_msig_2mass, type=DOUBLE"`
}

func (a Allwise) GetId() string {
	return a.ID
}

func (a Allwise) GetCatalog() string {
	return "AllWISE"
}

func (m InsertAllwiseParams) GetId() string {
	return m.ID
}

func (m GetAllwiseFromPixelsRow) GetId() string {
	return m.ID
}

func (m GetAllwiseFromPixelsRow) GetCoordinates() (float64, float64) {
	return m.Ra, m.Dec
}

func (m GetAllwiseFromPixelsRow) GetCatalog() string {
	return "AllWISE"
}

func (q *Queries) InsertAllwiseWithoutParams(ctx context.Context, arg Allwise) error {
	_, err := q.db.ExecContext(ctx, insertAllwise,
		arg.ID,
		arg.Cntr,
		arg.W1mpro,
		arg.W1sigmpro,
		arg.W2mpro,
		arg.W2sigmpro,
		arg.W3mpro,
		arg.W3sigmpro,
		arg.W4mpro,
		arg.W4sigmpro,
		arg.JM2mass,
		arg.JMsig2mass,
		arg.HM2mass,
		arg.HMsig2mass,
		arg.KM2mass,
		arg.KMsig2mass,
	)
	return err
}
