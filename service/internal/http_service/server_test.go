package httpservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

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
		assert.Equal(t, result, parsedRa)
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
		assert.Equal(t, result, expectedDec)
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
		assert.Equal(t, result, parsedRadius)
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
		assert.Equal(t, expectedResult, result)
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
		assert.Equal(t, expectedResult, result)
	}
}

func TestConesearchValidation(t *testing.T) {
	type Expected struct {
		Status       int
		ErrorMessage string
	}
	testCases := map[string]Expected{
		"/conesearch":                                               {400, "RA can't be empty\n"},
		"/conesearch?ra=1":                                          {400, "Dec can't be empty\n"},
		"/conesearch?ra=1&dec=1":                                    {400, "Radius can't be empty\n"},
		"/conesearch?ra=1&dec=1&radius=1":                           {200, ""},
		"/conesearch?ra=1&dec=1&radius=1&catalog=a":                 {400, "Catalog must be one of [all wise vlass lsdr10]\n"},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise":              {200, ""},
		"/conesearch?ra=1&dec=1&radius=1&catalog=wise&nneighbor=-1": {400, "Nneighbor must be a positive integer\n"},
	}

	router := SetupServer()
	for testPath, expected := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testPath, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, expected.Status, w.Code)
		assert.Equal(t, expected.ErrorMessage, w.Body.String())
	}
}
