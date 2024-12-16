package httpservice

import (
	"fmt"
	"strconv"
)

type ParseError struct {
	ErrValue string
	Field    string
	Reason   string
}

func (e ParseError) Error() string {
	err := "Could not parse field %s with value %s:\n%s"
	err = fmt.Sprintf(err, e.Field, e.ErrValue, e.Reason)
	return err
}

func NewParseError(errValue, field, reason string) error {
	return ParseError{
		ErrValue: errValue,
		Field:    field,
		Reason:   reason,
	}
}

func parseRadius(rad string) (float64, error) {
	radius, err := strconv.ParseFloat(rad, 64)
	if err != nil {
		return -999, NewParseError(rad, "radius", "Could not parse float.")
	}
	return radius, nil
}

func parseRa(ra string) (float64, error) {
	parsedRa, err := strconv.ParseFloat(ra, 64)
	if err != nil {
		return -999, NewParseError(ra, "RA", "Could not parse float.")
	}
	return parsedRa, nil
}

func parseDec(dec string) (float64, error) {
	parsedDec, err := strconv.ParseFloat(dec, 64)
	if err != nil {
		return -999, NewParseError(dec, "Dec", "Could not parse float.")
	}
	return parsedDec, nil
}

func parseNneighbor(nneighbor string) (int, error) {
	parsedNneighbor, err := strconv.Atoi(nneighbor)
	if err != nil {
		return -999, NewParseError(nneighbor, "nneighbor", "Could not parse int.")
	}
	return parsedNneighbor, nil
}
