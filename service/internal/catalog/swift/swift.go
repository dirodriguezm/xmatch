package swift

import (
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// InputSchema maps the LSXPS Sources CSV files (LSXPS_Sources.csv and its
// clean / ultra-clean variants) from https://www.swift.ac.uk/LSXPS/.
//
// The CSV reader fills fields by position, not by header name, and ignores
// any extra columns. These fields must therefore match the first four
// columns of the file, in order; columns 5 onwards (Err90, AstromType, l,
// b, ...) are ignored.
type InputSchema struct {
	LSXPSID string  `parquet:"name=LSXPS_ID, type=BYTE_ARRAY"` // column 1
	IAUName string  `parquet:"name=IAUName, type=BYTE_ARRAY"`  // column 2
	RA      float64 `parquet:"name=RA, type=DOUBLE"`           // column 3
	Dec     float64 `parquet:"name=Decl, type=DOUBLE"`         // column 4
}

func (r InputSchema) GetCoordinates() (float64, float64) {
	return r.RA, r.Dec
}

// MastercatID uses IAUName (e.g. "LSXPS J123456.7+123456"), which is unique
// and does not collide with other catalogs' ids in mastercat.
func (r InputSchema) MastercatID() string {
	return r.IAUName
}

func (schema InputSchema) ConvertToMetadata() (repository.Swift, error) {
	return repository.Swift{ID: schema.IAUName}, nil
}
