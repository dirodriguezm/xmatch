package httpservice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRaValidation(t *testing.T) {
	testCases := map[string]float64{
		"":    -999,
		"aaa": -999,
		"-1":  -999,
		"361": -999,
		"0":   0,
		"1":   1,
		"360": 360,
		"5.5": 5.5,
	}

	for strRa, parsedRa := range testCases {
		result, _ := validateRa(strRa)
		require.Equal(t, result, parsedRa)
	}
}

func TestDecValidation(t *testing.T) {
	testCases := map[string]float64{
		"":    -999,
		"aaa": -999,
		"-91": -999,
		"91":  -999,
		"-90": -90,
		"90":  90,
		"5.5": 5.5,
	}

	for strDec, expectedDec := range testCases {
		result, _ := validateDec(strDec)
		require.Equal(t, result, expectedDec)
	}
}

func TestRadiusValidation(t *testing.T) {
	testCases := map[string]float64{
		"":    -999,
		"aaa": -999,
		"-1":  -999,
		"0":   -999,
		"1":   1,
		"5.5": 5.5,
	}

	for strRadius, parsedRadius := range testCases {
		result, _ := validateRadius(strRadius)
		require.Equal(t, result, parsedRadius)
	}
}

func TestCatalogValidation(t *testing.T) {
	testCases := map[string]string{
		"all":            "all",
		"wise":           "wise",
		"vlass":          "vlass",
		"lsdr10":         "lsdr10",
		"ALL":            "all",
		"WISE":           "wise",
		"VLASS":          "vlass",
		"LSDR10":         "lsdr10",
		"something else": "",
	}
	for testCase, expectedResult := range testCases {
		result, _ := validateCatalog(testCase)
		require.Equal(t, expectedResult, result)
	}
}

func TestNneighborValidation(t *testing.T) {
	testCases := map[string]int{
		"":    -999,
		"aaa": -999,
		"-1":  -999,
		"0":   -999,
		"1":   1,
	}
	for testCase, expectedResult := range testCases {
		result, _ := validateNneighbor(testCase)
		require.Equal(t, expectedResult, result)
	}
}
