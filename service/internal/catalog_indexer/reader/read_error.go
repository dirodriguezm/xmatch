package reader

import (
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/source"
)

type ReadError struct {
	CurrentReader int
	OriginalError error
	Source        *source.Source
	Message       string
}

func (e ReadError) Error() string {
	err := "Error while reading from %s. %s\n%s"
	return fmt.Sprintf(
		err,
		e.Source.Reader[e.CurrentReader].Url,
		e.Message,
		e.OriginalError,
	)
}

func NewReadError(currentReader int, err error, src *source.Source, msg string) error {
	return ReadError{
		CurrentReader: currentReader,
		OriginalError: err,
		Source:        src,
		Message:       msg,
	}
}
