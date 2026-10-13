package allwise

import (
	"database/sql"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type InputSchema struct {
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

func (r InputSchema) GetCoordinates() (float64, float64) {
	ra := 0.0
	dec := 0.0
	if r.Ra != nil {
		ra = *r.Ra
	}
	if r.Dec != nil {
		dec = *r.Dec
	}
	return ra, dec
}

func (r InputSchema) MastercatID() string {
	if r.Source_id != nil {
		return *r.Source_id
	}
	return ""
}

func (schema InputSchema) ConvertToMetadata() (repository.Allwise, error) {
	allwise := repository.Allwise{}
	if schema.Source_id != nil {
		allwise.ID = *schema.Source_id
	}
	if schema.Cntr != nil {
		allwise.Cntr = *schema.Cntr
	}
	if schema.W1mpro != nil {
		allwise.W1mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W1mpro, Valid: true}}
	} else {
		allwise.W1mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W1sigmpro != nil {
		allwise.W1sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W1sigmpro, Valid: true}}
	} else {
		allwise.W1sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2mpro != nil {
		allwise.W2mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W2mpro, Valid: true}}
	} else {
		allwise.W2mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W2sigmpro != nil {
		allwise.W2sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W2sigmpro, Valid: true}}
	} else {
		allwise.W2sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3mpro != nil {
		allwise.W3mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W3mpro, Valid: true}}
	} else {
		allwise.W3mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W3sigmpro != nil {
		allwise.W3sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W3sigmpro, Valid: true}}
	} else {
		allwise.W3sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4mpro != nil {
		allwise.W4mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W4mpro, Valid: true}}
	} else {
		allwise.W4mpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.W4sigmpro != nil {
		allwise.W4sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.W4sigmpro, Valid: true}}
	} else {
		allwise.W4sigmpro = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_m_2mass != nil {
		allwise.JM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.J_m_2mass, Valid: true}}
	} else {
		allwise.JM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.J_msig_2mass != nil {
		allwise.JMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.J_msig_2mass, Valid: true}}
	} else {
		allwise.JMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_m_2mass != nil {
		allwise.HM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.H_m_2mass, Valid: true}}
	} else {
		allwise.HM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.H_msig_2mass != nil {
		allwise.HMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.H_msig_2mass, Valid: true}}
	} else {
		allwise.HMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_m_2mass != nil {
		allwise.KM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.K_m_2mass, Valid: true}}
	} else {
		allwise.KM2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	if schema.K_msig_2mass != nil {
		allwise.KMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: *schema.K_msig_2mass, Valid: true}}
	} else {
		allwise.KMsig2mass = repository.NullFloat64{NullFloat64: sql.NullFloat64{Float64: -9999.0, Valid: false}}
	}
	return allwise, nil
}
