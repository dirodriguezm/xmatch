package repository

import "database/sql"

func (q *Queries) GetDbInstance() *sql.DB {
	return q.db.(*sql.DB)
}
