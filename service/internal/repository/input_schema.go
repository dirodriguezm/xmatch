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

import (
	"strconv"
	"strings"
)

type InputSchema interface {
	ToMastercat(ipix int64) ParquetMastercat
	ToMetadata() any
	GetCoordinates() (float64, float64)
	SetField(string, any)
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

func (a *AllwiseInputSchema) ToMastercat(ipix int64) ParquetMastercat {
	catalog := "allwise"
	return ParquetMastercat{
		ID:   a.Source_id,
		Ipix: &ipix,
		Ra:   a.Ra,
		Dec:  a.Dec,
		Cat:  &catalog,
	}
}

func (a *AllwiseInputSchema) GetCoordinates() (float64, float64) {
	return *a.Ra, *a.Dec
}

func (a *AllwiseInputSchema) ToMetadata() any {
	return AllwiseMetadata{
		Source_id:    a.Source_id,
		W1mpro:       a.W1mpro,
		W1sigmpro:    a.W1sigmpro,
		W2mpro:       a.W2mpro,
		W2sigmpro:    a.W2sigmpro,
		W3mpro:       a.W3mpro,
		W3sigmpro:    a.W3sigmpro,
		W4mpro:       a.W4mpro,
		W4sigmpro:    a.W4sigmpro,
		J_m_2mass:    a.J_m_2mass,
		H_m_2mass:    a.H_m_2mass,
		K_m_2mass:    a.K_m_2mass,
		J_msig_2mass: a.J_msig_2mass,
		H_msig_2mass: a.H_msig_2mass,
		K_msig_2mass: a.K_msig_2mass,
	}
}

func (a *AllwiseInputSchema) SetField(name string, val any) {
	switch n := strings.ToLower(name); n {
	case "source_id":
		a.Source_id = val.(*string)
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

func (v *VlassInputSchema) ToMastercat(ipix int64) ParquetMastercat {
	catalog := "vlass"
	return ParquetMastercat{
		ID:   v.Component_name,
		Ipix: &ipix,
		Ra:   v.RA,
		Dec:  v.DEC,
		Cat:  &catalog,
	}
}

func (v *VlassInputSchema) GetCoordinates() (float64, float64) {
	return *v.RA, *v.DEC
}

func (v *VlassInputSchema) ToMetadata() any {
	return VlassMetadata{}
}

func (v *VlassInputSchema) SetField(name string, val any) {
	switch n := strings.ToLower(name); n {
	case "component_name":
		parsed := val.(string)
		v.Component_name = &parsed
	case "ra":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.RA = &parsed
	case "dec":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.DEC = &parsed
	case "era":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.ERA = &parsed
	case "edec":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.EDEC = &parsed
	case "totalflux":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.TotalFlux = &parsed
	case "etotalflux":
		parsed, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			panic(err)
		}
		v.ETotalFlux = &parsed
	}
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

func (v *VlassObjectSchema) ToMetadata() any {
	return VlassMetadata{
		Id:    v.Id,
		Ra:    v.Ra,
		Dec:   v.Dec,
		Era:   v.Era,
		Edec:  v.Edec,
		Flux:  v.Flux,
		EFlux: v.EFlux,
	}
}

func (v *VlassObjectSchema) GetCoordinates() (float64, float64) {
	return *v.Ra, *v.Dec
}

func (v *VlassObjectSchema) ToMastercat(ipix int64) ParquetMastercat {
	catalog := "vlass"
	return ParquetMastercat{
		ID:   v.Id,
		Ipix: &ipix,
		Ra:   v.Ra,
		Dec:  v.Dec,
		Cat:  &catalog,
	}
}

func (v *VlassObjectSchema) SetField(name string, val any) {
	switch n := strings.ToLower(name); n {
	case "id":
		v.Id = val.(*string)
	case "ra":
		v.Ra = val.(*float64)
	case "dec":
		v.Dec = val.(*float64)
	case "era":
		v.Era = val.(*float64)
	case "edec":
		v.Edec = val.(*float64)
	case "flux":
		v.Flux = val.(*float64)
	case "eflux":
		v.EFlux = val.(*float64)
	}
}
