package repository

type Repository interface {
	FindObjectIds(pixelList []int64) ([]string, error)
}
