package repository

import "xmatch/service/internal/core"

type SqliteRepository struct {
}

func (repository *SqliteRepository) FindObjects(pixelList []int64) ([]core.MastercatObject, error) {
	return nil, nil
}
