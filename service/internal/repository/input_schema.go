// Copyright 2024-2025 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import "database/sql"

type InputSchema interface {
	ToMastercat(ipix int64) Mastercat
	ToMetadata() any
	GetCoordinates() (float64, float64)
	GetId() string
}

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

func (a *AllwiseInputSchema) ToMastercat(ipix int64) Mastercat {
	return Mastercat{
		ID:   *a.Source_id,
		Ipix: ipix,
		Ra:   *a.Ra,
		Dec:  *a.Dec,
		Cat:  "allwise",
	}
}

func (a *AllwiseInputSchema) GetCoordinates() (float64, float64) {
	return *a.Ra, *a.Dec
}

func (a *AllwiseInputSchema) ToMetadata() any {
	return Allwise{
		ID:         *a.Source_id,
		W1mpro:     sql.NullFloat64{Float64: defaultVal(a.W1mpro), Valid: a.W1mpro != nil},
		W1sigmpro:  sql.NullFloat64{Float64: defaultVal(a.W1sigmpro), Valid: a.W1sigmpro != nil},
		W2mpro:     sql.NullFloat64{Float64: defaultVal(a.W2mpro), Valid: a.W2mpro != nil},
		W2sigmpro:  sql.NullFloat64{Float64: defaultVal(a.W2sigmpro), Valid: a.W2sigmpro != nil},
		W3mpro:     sql.NullFloat64{Float64: defaultVal(a.W3mpro), Valid: a.W3mpro != nil},
		W3sigmpro:  sql.NullFloat64{Float64: defaultVal(a.W3sigmpro), Valid: a.W3sigmpro != nil},
		W4mpro:     sql.NullFloat64{Float64: defaultVal(a.W4mpro), Valid: a.W4mpro != nil},
		W4sigmpro:  sql.NullFloat64{Float64: defaultVal(a.W4sigmpro), Valid: a.W4sigmpro != nil},
		JM2mass:    sql.NullFloat64{Float64: defaultVal(a.J_m_2mass), Valid: a.J_m_2mass != nil},
		HM2mass:    sql.NullFloat64{Float64: defaultVal(a.H_m_2mass), Valid: a.H_m_2mass != nil},
		KM2mass:    sql.NullFloat64{Float64: defaultVal(a.K_m_2mass), Valid: a.K_m_2mass != nil},
		JMsig2mass: sql.NullFloat64{Float64: defaultVal(a.J_msig_2mass), Valid: a.J_msig_2mass != nil},
		HMsig2mass: sql.NullFloat64{Float64: defaultVal(a.H_msig_2mass), Valid: a.H_msig_2mass != nil},
		KMsig2mass: sql.NullFloat64{Float64: defaultVal(a.K_msig_2mass), Valid: a.K_msig_2mass != nil},
	}
}

func defaultVal(f *float64) float64 {
	if f == nil {
		return -9999.0
	}
	return *f
}

func (a *AllwiseInputSchema) GetId() string {
	return *a.Source_id
}

type VlassInputSchema struct {
	Component_name *string  `parquet:"name=Component_name, type=BYTE_ARRAY"`
	RA             *float64 `parquet:"name=RA, type=DOUBLE"`
	DEC            *float64 `parquet:"name=DEC, type=DOUBLE"`
	ERA            *float64 `parquet:"name=E_RA, type=DOUBLE"`
	EDEC           *float64 `parquet:"name=E_DEC, type=DOUBLE"`
	TotalFlux      *float64 `parquet:"name=Total_flux, type=DOUBLE"`
	ETotalFlux     *float64 `parquet:"name=E_Total_flux, type=DOUBLE"`
}

func (v *VlassInputSchema) GetId() string {
	return *v.Component_name
}

func (v *VlassInputSchema) ToMastercat(ipix int64) Mastercat {
	catalog := "vlass"
	return Mastercat{
		ID:   *v.Component_name,
		Ipix: ipix,
		Ra:   *v.RA,
		Dec:  *v.DEC,
		Cat:  catalog,
	}
}

func (v *VlassInputSchema) GetCoordinates() (float64, float64) {
	return *v.RA, *v.DEC
}

func (v *VlassInputSchema) ToMetadata() any {
	return nil
}

type VlassObjectSchema struct {
	Id    *string  `parquet:"name=id, type=BYTE_ARRAY"`
	Ra    *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec   *float64 `parquet:"name=dec, type=DOUBLE"`
	Era   *float64 `parquet:"name=e_ra, type=DOUBLE"`
	Edec  *float64 `parquet:"name=e_dec, type=DOUBLE"`
	Flux  *float64 `parquet:"name=flux, type=DOUBLE"`
	EFlux *float64 `parquet:"name=e_flux, type=DOUBLE"`
}

func (v *VlassObjectSchema) GetId() string {
	return *v.Id
}

// TODO: CREATE THIS THROUGH A DATABASE MODEL
type Vlass struct {
	Id    string
	Ra    float64
	Dec   float64
	Era   float64
	Edec  float64
	Flux  float64
	EFlux float64
}

func (v *VlassObjectSchema) ToMetadata() any {
	return Vlass{
		Id:    *v.Id,
		Ra:    *v.Ra,
		Dec:   *v.Dec,
		Era:   *v.Era,
		Edec:  *v.Edec,
		Flux:  *v.Flux,
		EFlux: *v.EFlux,
	}
}

func (v *VlassObjectSchema) GetCoordinates() (float64, float64) {
	return *v.Ra, *v.Dec
}

func (v *VlassObjectSchema) ToMastercat(ipix int64) Mastercat {
	catalog := "vlass"
	return Mastercat{
		ID:   *v.Id,
		Ipix: ipix,
		Ra:   *v.Ra,
		Dec:  *v.Dec,
		Cat:  catalog,
	}
}
