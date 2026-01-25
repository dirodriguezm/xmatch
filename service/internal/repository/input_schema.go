package repository

type InputSchema interface {
	GetCoordinates() (float64, float64)
	GetId() string
}
