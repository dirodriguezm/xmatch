package repository

import (
	"context"
	"database/sql"
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

func (schema AllwiseInputSchema) GetCoordinates() (float64, float64) {
	return *schema.Ra, *schema.Dec
}

func (schema AllwiseInputSchema) GetId() string {
	return *schema.Source_id
}

func FillAllwiseMetadata(schema InputSchema) Metadata {
	allwise := &Allwise{
		ID:   *schema.(AllwiseInputSchema).Source_id,
		Cntr: *schema.(AllwiseInputSchema).Cntr,
	}
	if schema.(AllwiseInputSchema).W1mpro != nil {
		allwise.W1mpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W1mpro, Valid: true}
	} else {
		allwise.W1mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W1sigmpro != nil {
		allwise.W1sigmpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W1sigmpro, Valid: true}
	} else {
		allwise.W1sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W2mpro != nil {
		allwise.W2mpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W2mpro, Valid: true}
	} else {
		allwise.W2mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W2sigmpro != nil {
		allwise.W2sigmpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W2sigmpro, Valid: true}
	} else {
		allwise.W2sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W3mpro != nil {
		allwise.W3mpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W3mpro, Valid: true}
	} else {
		allwise.W3mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W3sigmpro != nil {
		allwise.W3sigmpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W3sigmpro, Valid: true}
	} else {
		allwise.W3sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W4mpro != nil {
		allwise.W4mpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W4mpro, Valid: true}
	} else {
		allwise.W4mpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).W4sigmpro != nil {
		allwise.W4sigmpro = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).W4sigmpro, Valid: true}
	} else {
		allwise.W4sigmpro = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).J_m_2mass != nil {
		allwise.JM2mass = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).J_m_2mass, Valid: true}
	} else {
		allwise.JM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).J_msig_2mass != nil {
		allwise.JMsig2mass = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).J_msig_2mass, Valid: true}
	} else {
		allwise.JMsig2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).H_m_2mass != nil {
		allwise.HM2mass = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).H_m_2mass, Valid: true}
	} else {
		allwise.HM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).H_msig_2mass != nil {
		allwise.HMsig2mass = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).H_msig_2mass, Valid: true}
	} else {
		allwise.HMsig2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	if schema.(AllwiseInputSchema).K_m_2mass != nil {
		allwise.KM2mass = sql.NullFloat64{Float64: *schema.(AllwiseInputSchema).K_m_2mass, Valid: true}
	} else {
		allwise.KM2mass = sql.NullFloat64{Float64: -9999.0, Valid: false}
	}

	return allwise
}

func FillAllwiseMastercat(schema InputSchema, ipix int64) Mastercat {
	return Mastercat{
		ID:   *schema.(AllwiseInputSchema).Source_id,
		Ipix: ipix,
		Ra:   *schema.(AllwiseInputSchema).Ra,
		Dec:  *schema.(AllwiseInputSchema).Dec,
		Cat:  "allwise",
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
