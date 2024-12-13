package httpservice

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
)

type ParseError struct {
	Err      error
	ErrValue string
	Field    string
}

func (e ParseError) Error() string {
	err := "Could not parse field %s with value %s:\n%s"
	err = fmt.Sprintf(err, e.Field, e.ErrValue, e.Err.Error())
	return err
}

func NewParseError(e error, errValue string, field string) error {
	return ParseError{
		Err:      e,
		ErrValue: errValue,
		Field:    field,
	}
}

func parseRadius(rad string) (float64, error) {
	radius, err := strconv.ParseFloat(rad, 64)
	if err != nil {
		return -999, NewParseError(err, rad, "radius")
	}
	if err := conesearch.ValidateRadius(radius); err != nil {
		return -999, NewParseError(err, rad, "radius")
	}
	return radius, nil
}

func parseRa(ra string) (float64, error) {
	parsedRa, err := strconv.ParseFloat(ra, 64)
	if err != nil {
		return -999, NewParseError(err, ra, "RA")
	}
	if err := conesearch.ValidateRa(parsedRa); err != nil {
		return -999, NewParseError(err, ra, "RA")
	}
	return parsedRa, nil
}

func parseDec(dec string) (float64, error) {
	parsedDec, err := strconv.ParseFloat(dec, 64)
	if err != nil {
		return -999, NewParseError(err, dec, "Dec")
	}
	if err := conesearch.ValidateDec(parsedDec); err != nil {
		return -999, NewParseError(err, dec, "Dec")
	}
	return parsedDec, nil
}

func parseCatalog(catalog string) (string, error) {
	allowed := []string{"all", "wise", "vlass", "lsdr10"}
	isAllowed := false
	catalog = strings.ToLower(catalog)
	for _, cat := range allowed {
		if catalog == cat {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return "", NewParseError(fmt.Errorf("Only %s catalogs are available. Catalog %s not found.", allowed, catalog), catalog, "catalog")
	}
	return catalog, nil
}

func parseNneighbor(nneighbor string) (int, error) {
	parsedNneighbor, err := strconv.Atoi(nneighbor)
	if err != nil {
		return -999, NewParseError(err, nneighbor, "nneighbor")
	}
	if err := conesearch.ValidateNneighbor(parsedNneighbor); err != nil {
		return -999, NewParseError(err, nneighbor, "nneighbor")
	}
	return parsedNneighbor, nil
}
