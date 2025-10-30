package neowise

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
)

type NeowiseClient struct {
	url     string
	headers map[string]string
	catalog string
	columns []string
}

func NewNeowiseClient() *NeowiseClient {
	return &NeowiseClient{
		url:     "https://irsa.ipac.caltech.edu/cgi-bin/Gator/nph-query",
		headers: map[string]string{},
		catalog: "neowiser_p1bs_psd",
		columns: []string{"mjd", "ra", "dec", "w1mpro", "w1sigmpro", "w2mpro", "w2sigmpro", "allwise_cntr", "source_id"},
	}
}

func (client *NeowiseClient) FetchLightcurve(ra, dec, radius float64, nobjects int) lightcurve.ClientResult {
	u, err := url.Parse(client.url)
	if err != nil {
		return lightcurve.ClientResult{
			Error: fmt.Errorf("could not parse url: %s", err),
		}
	}

	u = addQueryParameters(u, map[string]string{
		"catalog":  client.catalog,
		"spatial":  "cone",
		"objstr":   fmt.Sprintf("%f %f", ra, dec),
		"radunits": "deg",
		"radius":   fmt.Sprintf("%f", arcSecToDeg(radius)),
		"outfmt":   "3",
		"selcols":  strings.Join(client.columns, ","),
	})

	resp, err := http.Get(u.String())
	if err != nil {
		return lightcurve.ClientResult{
			Error: fmt.Errorf("could not make request: %s", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return lightcurve.ClientResult{
			Error: fmt.Errorf("could not read response body: %s", err),
		}
	}

	resultTable, err := utils.NewVOTableFromBytes(body)
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not parse response: %s", err)}
	}

	lightcurveObjects, err := convertToLightcurveObject(resultTable)
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not convert response to lightcurve object: %s", err)}
	}

	return lightcurve.ClientResult{
		Error: nil,
		Lightcurve: lightcurve.Lightcurve{
			Detections: lightcurveObjects,
		},
	}
}

func addQueryParameters(u *url.URL, params map[string]string) *url.URL {
	newURL := *u

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}

	newURL.RawQuery = q.Encode()

	return &newURL
}

func arcSecToDeg(value float64) float64 {
	return value / 3600.0
}

func convertToLightcurveObject(detections *utils.VOTable) ([]lightcurve.LightcurveObject, error) {
	if len(detections.Resource.Tables) == 0 {
		return nil, fmt.Errorf("no tables found in response")
	}

	result := make([]lightcurve.LightcurveObject, len(detections.Resource.Tables[0].Data.TableData.Rows))

	for i, row := range detections.Resource.Tables[0].Data.TableData.Rows {
		if len(row.Columns) != len(detections.Resource.Tables[0].Fields) {
			return nil, fmt.Errorf("number of columns in row does not match number of fields")
		}
		detection, err := detectionFromColumns(row.Columns)
		if err != nil {
			return nil, fmt.Errorf("could not convert row to detection: %s", err)
		}
		result[i] = detection
	}

	return result, nil
}

func detectionFromColumns(columns []utils.Column) (lightcurve.LightcurveObject, error) {
	detection, err := setMjd(columns[0].Value, Detection{})
	if err != nil {
		return detection, err
	}
	detection, err = setRa(columns[1].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setDec(columns[2].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setW1mpro(columns[5].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setW1sigmpro(columns[6].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setW2mpro(columns[7].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setW2sigmpro(columns[8].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setCntr(columns[9].Value, detection)
	if err != nil {
		return detection, err
	}
	detection, err = setSourceId(columns[10].Value, detection)
	if err != nil {
		return detection, err
	}

	return detection, nil
}

func setRa(value string, detection Detection) (Detection, error) {
	ra, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse RA: %s", err)
	}
	detection.Ra = ra
	return detection, nil
}

func setDec(value string, detection Detection) (Detection, error) {
	dec, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse DEC: %s", err)
	}
	detection.Dec = dec

	return detection, nil
}

func setMjd(value string, detection Detection) (Detection, error) {
	mjd, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse MJD: %s", err)
	}
	detection.Mjd = mjd
	return detection, nil
}

func setW1mpro(value string, detection Detection) (Detection, error) {
	if value == "" {
		detection.W1mpro = -999
		return detection, nil
	}

	w1mpro, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse W1mpro: %s", err)
	}
	detection.W1mpro = w1mpro
	return detection, nil
}

func setW1sigmpro(value string, detection Detection) (Detection, error) {
	if value == "" {
		detection.W1sigmpro = -999
		return detection, nil
	}

	w1sigmpro, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return detection, fmt.Errorf("could not parse W1sigmpro: %s", err)
	}
	detection.W1sigmpro = float32(w1sigmpro)
	return detection, nil
}

func setW2mpro(value string, detection Detection) (Detection, error) {
	if value == "" {
		detection.W2mpro = -999
		return detection, nil
	}

	w2mpro, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse W2mpro: %s", err)
	}
	detection.W2mpro = w2mpro
	return detection, nil
}

func setW2sigmpro(value string, detection Detection) (Detection, error) {
	if value == "" {
		detection.W2sigmpro = -999
		return detection, nil
	}

	w2sigmpro, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return detection, fmt.Errorf("could not parse W2sigmpro: %s", err)
	}
	detection.W2sigmpro = float32(w2sigmpro)
	return detection, nil
}

func setCntr(value string, detection Detection) (Detection, error) {
	if value == "" {
		detection.Cntr = -999
		return detection, nil
	}
	cntr, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return detection, fmt.Errorf("could not parse CNTR: %s", err)
	}
	detection.Cntr = cntr
	return detection, nil
}

func setSourceId(value string, detection Detection) (Detection, error) {
	detection.Source_id = value
	return detection, nil
}
