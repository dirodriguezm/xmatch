package repository

type InputSchema interface {
	FillMastercat(ipix int64) Mastercat
	FillMetadata() Metadata
	GetCoordinates() (float64, float64)
	GetId() string
}
