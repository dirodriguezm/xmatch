package repository

type InputSchema interface {
	FillMastercat(dst *Mastercat, ipix int64)
	FillMetadata(dst Metadata)
	GetCoordinates() (float64, float64)
	GetId() string
}
