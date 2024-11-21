package source

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/config"
)

type Source struct {
	Reader      []io.Reader
	CatalogName string
	RaCol       string
	DecCol      string
	OidCol      string
	Nside       int
}

func NewSource(cfg *config.SourceConfig) (*Source, error) {
	if !validateSourceType(cfg.Type) {
		return nil, fmt.Errorf("Can't create source with type %s.", cfg.Type)
	}
	reader, err := sourceReader(cfg.Url)
	if err != nil {
		return nil, err
	}
	return &Source{
		Reader:      reader,
		CatalogName: cfg.CatalogName,
		RaCol:       cfg.RaCol,
		DecCol:      cfg.DecCol,
		OidCol:      cfg.OidCol,
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

func sourceReader(url string) ([]io.Reader, error) {
	if strings.HasPrefix(url, "file:") {
		reader, err := os.Open(strings.Split(url, "file:")[1])
		if err != nil {
			return nil, err
		}
		return []io.Reader{reader}, nil
	}
	if strings.HasPrefix(url, "buffer:") {
		return []io.Reader{&bytes.Buffer{}}, nil
	}
	if strings.HasPrefix(url, "files:") {
		entries, err := os.ReadDir(strings.Split(url, "files:")[1])
		if err != nil {
			return nil, err
		}
		readers := []io.Reader{}
		for _, entry := range entries {
			if entry.IsDir() {
				rdrs, err := sourceReader("files:" + entry.Name())
				if err != nil {
					return nil, err
				}
				readers = append(readers, rdrs...)
			}
			file, err := os.Open(entry.Name())
			if err != nil {
				return nil, err
			}
			readers = append(readers, file)
		}
		return readers, nil
	}
	return nil, fmt.Errorf("Could not parse URL: %s", url)
}
