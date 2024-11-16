package conesearch

import (
	"errors"
)

func ValidateRadius(rad float64) error {
	if rad <= 0 {
		return errors.New("Radius can't be lower than 0\n")
	}
	return nil
}

func ValidateRa(ra float64) error {
	if ra < 0 {
		return errors.New("RA can't be lower than 0\n")
	}
	if ra > 360 {
		return errors.New("RA can't be greater than 360\n")
	}
	return nil
}

func ValidateDec(dec float64) error {
	if dec < -90 {
		return errors.New("Dec can't be lower than -90\n")
	}
	if dec > 90 {
		return errors.New("Dec can't be greater than 90\n")
	}
	return nil
}

func ValidateNneighbor(nneighbor int) error {
	if nneighbor <= 0 {
		return errors.New("Nneighbor must be a positive integer\n")
	}
	return nil
}
