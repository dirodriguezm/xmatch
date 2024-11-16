package indexer

import (
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct {
	mock.Mock
}

func (r *ReaderMock) Catalog() string {
	return "catalog"
}

func (r *ReaderMock) RaCol() string {
	return "ra"
}

func (r *ReaderMock) DecCol() string {
	return "dec"
}

func (r *ReaderMock) ObjectIdCol() string {
	return "id"
}

func (r *ReaderMock) ReadBatch() ([]Row, error) {
	return make([]Row, 0), nil
}

func (r *ReaderMock) Read() ([]Row, error) {
	return make([]Row, 0), nil
}

func (r *ReaderMock) Start() {

}
