package repository

import "database/sql"

type AllwiseMetadata struct {
	Source_id  *string  `parquet:"name=source_id, type=BYTE_ARRAY"`
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

func (m AllwiseMetadata) ToAllwise() Allwise {
	return Allwise{
		ID:         *m.Source_id,
		W1mpro:     sql.NullFloat64{Float64: *m.W1mpro, Valid: m.W1mpro != nil},
		W1sigmpro:  sql.NullFloat64{Float64: *m.W1sigmpro, Valid: m.W1sigmpro != nil},
		W2mpro:     sql.NullFloat64{Float64: *m.W2mpro, Valid: m.W2mpro != nil},
		W2sigmpro:  sql.NullFloat64{Float64: *m.W2sigmpro, Valid: m.W2sigmpro != nil},
		W3mpro:     sql.NullFloat64{Float64: *m.W3mpro, Valid: m.W3mpro != nil},
		W3sigmpro:  sql.NullFloat64{Float64: *m.W3sigmpro, Valid: m.W3sigmpro != nil},
		W4mpro:     sql.NullFloat64{Float64: *m.W4mpro, Valid: m.W4mpro != nil},
		W4sigmpro:  sql.NullFloat64{Float64: *m.W4sigmpro, Valid: m.W4sigmpro != nil},
		JM2mass:    sql.NullFloat64{Float64: *m.J_m_2mass, Valid: m.J_m_2mass != nil},
		JMsig2mass: sql.NullFloat64{Float64: *m.J_msig_2mass, Valid: m.J_msig_2mass != nil},
		HM2mass:    sql.NullFloat64{Float64: *m.H_m_2mass, Valid: m.H_m_2mass != nil},
		HMsig2mass: sql.NullFloat64{Float64: *m.H_msig_2mass, Valid: m.H_msig_2mass != nil},
		KM2mass:    sql.NullFloat64{Float64: *m.K_m_2mass, Valid: m.K_m_2mass != nil},
		KMsig2mass: sql.NullFloat64{Float64: *m.K_msig_2mass, Valid: m.K_msig_2mass != nil},
	}
}

func (m AllwiseMetadata) ToInsertParams() InsertAllwiseParams {
	return InsertAllwiseParams{
		ID:         *m.Source_id,
		W1mpro:     sql.NullFloat64{Float64: nilOr0(m.W1mpro), Valid: m.W1mpro != nil},
		W1sigmpro:  sql.NullFloat64{Float64: nilOr0(m.W1sigmpro), Valid: m.W1sigmpro != nil},
		W2mpro:     sql.NullFloat64{Float64: nilOr0(m.W2mpro), Valid: m.W2mpro != nil},
		W2sigmpro:  sql.NullFloat64{Float64: nilOr0(m.W2sigmpro), Valid: m.W2sigmpro != nil},
		W3mpro:     sql.NullFloat64{Float64: nilOr0(m.W3mpro), Valid: m.W3mpro != nil},
		W3sigmpro:  sql.NullFloat64{Float64: nilOr0(m.W3sigmpro), Valid: m.W3sigmpro != nil},
		W4mpro:     sql.NullFloat64{Float64: nilOr0(m.W4mpro), Valid: m.W4mpro != nil},
		W4sigmpro:  sql.NullFloat64{Float64: nilOr0(m.W4sigmpro), Valid: m.W4sigmpro != nil},
		JM2mass:    sql.NullFloat64{Float64: nilOr0(m.J_m_2mass), Valid: m.J_m_2mass != nil},
		JMsig2mass: sql.NullFloat64{Float64: nilOr0(m.J_msig_2mass), Valid: m.J_msig_2mass != nil},
		HM2mass:    sql.NullFloat64{Float64: nilOr0(m.H_m_2mass), Valid: m.H_m_2mass != nil},
		HMsig2mass: sql.NullFloat64{Float64: nilOr0(m.H_msig_2mass), Valid: m.H_msig_2mass != nil},
		KM2mass:    sql.NullFloat64{Float64: nilOr0(m.K_m_2mass), Valid: m.K_m_2mass != nil},
		KMsig2mass: sql.NullFloat64{Float64: nilOr0(m.K_msig_2mass), Valid: m.K_msig_2mass != nil},
	}
}

func nilOr0(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
