package swift

import (
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

// InputSchema describes one row of the Swift/LSXPS input file.
// The CSV reader maps columns by struct field order, so keep this
// order in sync with the input file (id, ra, dec).
type InputSchema struct {
	ID  string  `parquet:"name=id, type=BYTE_ARRAY"`
	RA  float64 `parquet:"name=ra, type=DOUBLE"`
	Dec float64 `parquet:"name=dec, type=DOUBLE"`
}

func (r InputSchema) GetCoordinates() (float64, float64) {
	return r.RA, r.Dec
}

func (r InputSchema) MastercatID() string {
	return r.ID
}

func (schema InputSchema) ConvertToMetadata() (repository.Swift, error) {
	return repository.Swift{ID: schema.ID}, nil
}
