package repository

type VlassInputSchema struct {
	Component_name *string  `parquet:"name=Component_name, type=BYTE_ARRAY"`
	RA             *float64 `parquet:"name=RA, type=DOUBLE"`
	DEC            *float64 `parquet:"name=DEC, type=DOUBLE"`
	ERA            *float64 `parquet:"name=E_RA, type=DOUBLE"`
	EDEC           *float64 `parquet:"name=E_DEC, type=DOUBLE"`
	TotalFlux      *float64 `parquet:"name=Total_flux, type=DOUBLE"`
	ETotalFlux     *float64 `parquet:"name=E_Total_flux, type=DOUBLE"`
}

func (v VlassInputSchema) GetId() string {
	return *v.Component_name
}

func (v VlassInputSchema) GetCoordinates() (float64, float64) {
	return *v.RA, *v.DEC
}

func (v VlassInputSchema) FillMetadata(dst Metadata)                {}
func (v VlassInputSchema) FillMastercat(dst *Mastercat, ipix int64) {}

type VlassObjectSchema struct {
	Id    *string  `parquet:"name=id, type=BYTE_ARRAY"`
	Ra    *float64 `parquet:"name=ra, type=DOUBLE"`
	Dec   *float64 `parquet:"name=dec, type=DOUBLE"`
	Era   *float64 `parquet:"name=e_ra, type=DOUBLE"`
	Edec  *float64 `parquet:"name=e_dec, type=DOUBLE"`
	Flux  *float64 `parquet:"name=flux, type=DOUBLE"`
	EFlux *float64 `parquet:"name=e_flux, type=DOUBLE"`
}

func (v VlassObjectSchema) GetId() string {
	return *v.Id
}

func (v VlassObjectSchema) GetCoordinates() (float64, float64) {
	return *v.Ra, *v.Dec
}

func (v VlassObjectSchema) FillMetadata(dst Metadata)                {}
func (v VlassObjectSchema) FillMastercat(dst *Mastercat, ipix int64) {}
