package repository

import "database/sql"

type AllwiseMetadata struct {
	Source_id    *string  `parquet:"name=source_id, type=BYTE_ARRAY"`
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

func (m Allwise) ToAllwiseMetadata() AllwiseMetadata {
	defaultVal := -9999.0
	w1mpro := m.W1mpro.Float64
	if !m.W1mpro.Valid {
		w1mpro = defaultVal
	}
	w1sigmpro := m.W1sigmpro.Float64
	if !m.W1sigmpro.Valid {
		w1sigmpro = defaultVal
	}
	w2mpro := m.W2mpro.Float64
	if !m.W2mpro.Valid {
		w2mpro = defaultVal
	}
	w2sigmpro := m.W2sigmpro.Float64
	if !m.W2sigmpro.Valid {
		w2sigmpro = defaultVal
	}
	w3mpro := m.W3mpro.Float64
	if !m.W3mpro.Valid {
		w3mpro = defaultVal
	}
	w3sigmpro := m.W3sigmpro.Float64
	if !m.W3sigmpro.Valid {
		w3sigmpro = defaultVal
	}
	w4mpro := m.W4mpro.Float64
	if !m.W4mpro.Valid {
		w4mpro = defaultVal
	}
	w4sigmpro := m.W4sigmpro.Float64
	if !m.W4sigmpro.Valid {
		w4sigmpro = defaultVal
	}
	jM2mass := m.JM2mass.Float64
	if !m.JM2mass.Valid {
		jM2mass = defaultVal
	}
	jMsig2mass := m.JMsig2mass.Float64
	if !m.JMsig2mass.Valid {
		jMsig2mass = defaultVal
	}
	hM2mass := m.HM2mass.Float64
	if !m.HM2mass.Valid {
		hM2mass = defaultVal
	}
	hMsig2mass := m.HMsig2mass.Float64
	if !m.HMsig2mass.Valid {
		hMsig2mass = defaultVal
	}
	kM2mass := m.KM2mass.Float64
	if !m.KM2mass.Valid {
		kM2mass = defaultVal
	}
	kMsig2mass := m.KMsig2mass.Float64
	if !m.KMsig2mass.Valid {
		kMsig2mass = defaultVal
	}

	return AllwiseMetadata{
		Source_id:    &m.ID,
		W1mpro:       &w1mpro,
		W1sigmpro:    &w1sigmpro,
		W2mpro:       &w2mpro,
		W2sigmpro:    &w2sigmpro,
		W3mpro:       &w3mpro,
		W3sigmpro:    &w3sigmpro,
		W4mpro:       &w4mpro,
		W4sigmpro:    &w4sigmpro,
		J_m_2mass:    &jM2mass,
		H_m_2mass:    &hM2mass,
		K_m_2mass:    &kM2mass,
		J_msig_2mass: &jMsig2mass,
		H_msig_2mass: &hMsig2mass,
		K_msig_2mass: &kMsig2mass,
	}
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
