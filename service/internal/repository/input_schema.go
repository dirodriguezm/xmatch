package repository

import "strings"

type InputSchema interface {
	ToMastercat() Mastercat
	SetField(string, any)
}

type AllwiseInputSchema struct {
	Designation  *string  `parquet:"name=designation, type=BYTE_ARRAY"`
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

type AllwiseMetadata struct {
	Designation  string  `parquet:"name=designation, type=BYTE_ARRAY"`
	W1mpro       float64 `parquet:"name=w1mpro, type=DOUBLE"`
	W1sigmpro    float64 `parquet:"name=w1sigmpro, type=DOUBLE"`
	W2mpro       float64 `parquet:"name=w2mpro, type=DOUBLE"`
	W2sigmpro    float64 `parquet:"name=w2sigmpro, type=DOUBLE"`
	W3mpro       float64 `parquet:"name=w3mpro, type=DOUBLE"`
	W3sigmpro    float64 `parquet:"name=w3sigmpro, type=DOUBLE"`
	W4mpro       float64 `parquet:"name=w4mpro, type=DOUBLE"`
	W4sigmpro    float64 `parquet:"name=w4sigmpro, type=DOUBLE"`
	J_m_2mass    float64 `parquet:"name=j_m_2mass, type=DOUBLE"`
	H_m_2mass    float64 `parquet:"name=h_m_2mass, type=DOUBLE"`
	K_m_2mass    float64 `parquet:"name=k_m_2mass, type=DOUBLE"`
	J_msig_2mass float64 `parquet:"name=j_msig_2mass, type=DOUBLE"`
	H_msig_2mass float64 `parquet:"name=h_msig_2mass, type=DOUBLE"`
	K_msig_2mass float64 `parquet:"name=k_msig_2mass, type=DOUBLE"`
}

func (a *AllwiseInputSchema) ToMastercat() Mastercat {
	return Mastercat{
		ID:   *a.Designation,
		Ipix: 0,
		Ra:   *a.Ra,
		Dec:  *a.Dec,
		Cat:  "allwise",
	}
}

func (a *AllwiseInputSchema) ToMetadata() AllwiseMetadata {
	return AllwiseMetadata{
		Designation:  *a.Designation,
		W1mpro:       *a.W1mpro,
		W1sigmpro:    *a.W1sigmpro,
		W2mpro:       *a.W2mpro,
		W2sigmpro:    *a.W2sigmpro,
		W3mpro:       *a.W3mpro,
		W3sigmpro:    *a.W3sigmpro,
		W4mpro:       *a.W4mpro,
		W4sigmpro:    *a.W4sigmpro,
		J_m_2mass:    *a.J_m_2mass,
		H_m_2mass:    *a.H_m_2mass,
		K_m_2mass:    *a.K_m_2mass,
		J_msig_2mass: *a.J_msig_2mass,
		H_msig_2mass: *a.H_msig_2mass,
		K_msig_2mass: *a.K_msig_2mass,
	}
}

func (a *AllwiseInputSchema) SetField(name string, val any) {
	switch n := strings.ToLower(name); n {
	case "designation":
		a.Designation = val.(*string)
	case "ra":
		a.Ra = val.(*float64)
	case "dec":
		a.Dec = val.(*float64)
	case "w1mpro":
		a.W1mpro = val.(*float64)
	case "w1sigmpro":
		a.W1sigmpro = val.(*float64)
	case "w2mpro":
		a.W2mpro = val.(*float64)
	case "w2sigmpro":
		a.W2sigmpro = val.(*float64)
	case "w3mpro":
		a.W3mpro = val.(*float64)
	case "w3sigmpro":
		a.W3sigmpro = val.(*float64)
	case "w4mpro":
		a.W4mpro = val.(*float64)
	case "w4sigmpro":
		a.W4sigmpro = val.(*float64)
	case "j_m_2mass":
		a.J_m_2mass = val.(*float64)
	case "h_m_2mass":
		a.H_m_2mass = val.(*float64)
	case "k_m_2mass":
		a.K_m_2mass = val.(*float64)
	case "j_msig_2mass":
		a.J_msig_2mass = val.(*float64)
	case "h_msig_2mass":
		a.H_msig_2mass = val.(*float64)
	case "k_msig_2mass":
		a.K_msig_2mass = val.(*float64)
	}
}
