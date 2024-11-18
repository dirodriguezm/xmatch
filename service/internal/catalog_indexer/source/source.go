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
	Reader      io.Reader
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

func sourceReader(url string) (io.Reader, error) {
	if strings.HasPrefix(url, "file:") {
		return os.Open(url)
	}
	if strings.HasPrefix(url, "buffer:") {
		return &bytes.Buffer{}, nil
	}
	return nil, fmt.Errorf("Could not parse URL: %s", url)
}
