package repository

import (
	"context"
	"database/sql"
)

type AllwiseInputSchema struct {
	Source_id    *string  `parquet:"name=source_id, type=BYTE_ARRAY"`
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

func (schema AllwiseInputSchema) GetCoordinates() (float64, float64) {
	return *schema.Ra, *schema.Dec
}

func (schema AllwiseInputSchema) GetId() string {
	return *schema.Source_id
}

func (schema AllwiseInputSchema) FillMetadata(dst Metadata) {
	dst.(*Allwise).ID = *schema.Source_id
	if schema.W1mpro != nil {
		dst.(*Allwise).W1mpro = sql.NullFloat64{Float64: *schema.W1mpro, Valid: true}
	} else {
		dst.(*Allwise).W1mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W1sigmpro != nil {
		dst.(*Allwise).W1sigmpro = sql.NullFloat64{Float64: *schema.W1sigmpro, Valid: true}
	} else {
		dst.(*Allwise).W1sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W2mpro != nil {
		dst.(*Allwise).W2mpro = sql.NullFloat64{Float64: *schema.W2mpro, Valid: true}
	} else {
		dst.(*Allwise).W2mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W2sigmpro != nil {
		dst.(*Allwise).W2sigmpro = sql.NullFloat64{Float64: *schema.W2sigmpro, Valid: true}
	} else {
		dst.(*Allwise).W2sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W3mpro != nil {
		dst.(*Allwise).W3mpro = sql.NullFloat64{Float64: *schema.W3mpro, Valid: true}
	} else {
		dst.(*Allwise).W3mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W3sigmpro != nil {
		dst.(*Allwise).W3sigmpro = sql.NullFloat64{Float64: *schema.W3sigmpro, Valid: true}
	} else {
		dst.(*Allwise).W3sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W4mpro != nil {
		dst.(*Allwise).W4mpro = sql.NullFloat64{Float64: *schema.W4mpro, Valid: true}
	} else {
		dst.(*Allwise).W4mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.W4sigmpro != nil {
		dst.(*Allwise).W4sigmpro = sql.NullFloat64{Float64: *schema.W4sigmpro, Valid: true}
	} else {
		dst.(*Allwise).W4sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.J_m_2mass != nil {
		dst.(*Allwise).JM2mass = sql.NullFloat64{Float64: *schema.J_m_2mass, Valid: true}
	} else {
		dst.(*Allwise).JM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.J_msig_2mass != nil {
		dst.(*Allwise).JMsig2mass = sql.NullFloat64{Float64: *schema.J_msig_2mass, Valid: true}
	} else {
		dst.(*Allwise).JMsig2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.H_m_2mass != nil {
		dst.(*Allwise).HM2mass = sql.NullFloat64{Float64: *schema.H_m_2mass, Valid: true}
	} else {
		dst.(*Allwise).HM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.H_msig_2mass != nil {
		dst.(*Allwise).HMsig2mass = sql.NullFloat64{Float64: *schema.H_msig_2mass, Valid: true}
	} else {
		dst.(*Allwise).HMsig2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.K_m_2mass != nil {
		dst.(*Allwise).KM2mass = sql.NullFloat64{Float64: *schema.K_m_2mass, Valid: true}
	} else {
		dst.(*Allwise).KM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}
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
