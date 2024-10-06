package repository

type sqliteRepository struct {
}

func (repository *sqliteRepository) FindObjectIds(pixelList []int64) ([]string, error) {
	return nil, nil
}

func NewSqliteRepository() Repository {
	return &sqliteRepository{}
}
