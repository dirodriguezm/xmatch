package core

type Row map[string]any

type Reader interface {
	Read() ([]Row, error)
	ReadBatch() ([]Row, error)
}
