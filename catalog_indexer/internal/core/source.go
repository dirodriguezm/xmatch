package core

import (
	"fmt"
	"io"
	"os"
)

type Source struct {
	Reader      io.Reader
	CatalogName string
	RaCol       string
	DecCol      string
	OidCol      string
}

func NewSource(stype string, url string, catalogName string, raCol string, decCol string, oidCol string) (*Source, error) {
	if !validateSourceType(stype) {
		return nil, fmt.Errorf("Can't create source with type %s.", stype)
	}
	reader, err := sourceReader(url)
	if err != nil {
		return nil, err
	}
	return &Source{
		Reader:      reader,
		CatalogName: catalogName,
		RaCol:       raCol,
		DecCol:      decCol,
		OidCol:      oidCol,
	}, nil
}

func validateSourceType(stype string) bool {
	allowedTypes := []string{"csv"}
	for _, t := range allowedTypes {
		if t == stype {
			return true
		}
	}
	return false
}

func sourceReader(url string) (io.Reader, error) {
	// TODO: parse url here
	return os.Open(url)
}
