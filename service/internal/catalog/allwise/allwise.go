package allwise

import (
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
	allwise := repository.Allwise{
		W1mpro:     schema.W1mpro,
		W1sigmpro:  schema.W1sigmpro,
		W2mpro:     schema.W2mpro,
		W2sigmpro:  schema.W2sigmpro,
		W3mpro:     schema.W3mpro,
		W3sigmpro:  schema.W3sigmpro,
		W4mpro:     schema.W4mpro,
		W4sigmpro:  schema.W4sigmpro,
		JM2mass:    schema.J_m_2mass,
		JMsig2mass: schema.J_msig_2mass,
		HM2mass:    schema.H_m_2mass,
		HMsig2mass: schema.H_msig_2mass,
		KM2mass:    schema.K_m_2mass,
		KMsig2mass: schema.K_msig_2mass,
	}
	if schema.Source_id != nil {
		allwise.ID = *schema.Source_id
	}
	if schema.Cntr != nil {
		allwise.Cntr = *schema.Cntr
	}
	return allwise, nil
}
