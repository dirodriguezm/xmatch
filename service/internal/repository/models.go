// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

type Catalog struct {
	Name  string
	Nside int64
}

type Mastercat struct {
	ID   string
	Ipix int64
	Ra   float64
	Dec  float64
	Cat  string
}
