package repository

func (schema AllwiseInputSchema) FillMastercat(dst *Mastercat, ipix int64) {
	dst.ID = *schema.Source_id
	dst.Ra = *schema.Ra
	dst.Dec = *schema.Dec
	dst.Cat = "allwise"
	dst.Ipix = ipix
}
