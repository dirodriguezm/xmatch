package ztfdr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
)

type lightCurveResponse struct {
	Id       int64     `json:"_id"`
	FilterId int       `json:"filterid"`
	FieldId  int       `json:"fieldid"`
	Nepochs  int       `json:"nepochs"`
	ObjRa    float64   `json:"objra"`
	ObjDec   float64   `json:"objdec"`
	Rcid     int       `json:"rcid"`
	Hmjd     []float64 `json:"hmjd"`
	Mag      []float64 `json:"mag"`
	Magerr   []float64 `json:"magerr"`
}

type ZtfDrClient struct {
	url string
}

func NewZtfDrClient() *ZtfDrClient {
	return &ZtfDrClient{
		url: "https://api.alerce.online/ztf/dr/v1/light_curve/",
	}
}

func (client *ZtfDrClient) FetchLightcurve(ra, dec, radius float64, _ int) lightcurve.ClientResult {
	u, err := url.Parse(client.url)
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not parse url: %w", err)}
	}

	u = addQueryParameters(u, map[string]string{
		"ra":     strconv.FormatFloat(ra, 'f', -1, 64),
		"dec":    strconv.FormatFloat(dec, 'f', -1, 64),
		"radius": strconv.FormatFloat(radius, 'f', -1, 64),
	})

	resp, err := http.Get(u.String())
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not make request: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return lightcurve.ClientResult{Lightcurve: lightcurve.Lightcurve{}}
	}
	if resp.StatusCode != http.StatusOK {
		return lightcurve.ClientResult{Error: fmt.Errorf("unexpected status code: %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not read response body: %w", err)}
	}

	var response lightCurveResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not parse response body: %w", err)}
	}

	detections, err := convertToLightcurveObjects(response)
	if err != nil {
		return lightcurve.ClientResult{Error: fmt.Errorf("could not convert response to lightcurve object: %w", err)}
	}

	return lightcurve.ClientResult{
		Lightcurve: lightcurve.Lightcurve{Detections: detections},
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

func convertToLightcurveObjects(response lightCurveResponse) ([]lightcurve.LightcurveObject, error) {
	if len(response.Hmjd) != len(response.Mag) || len(response.Hmjd) != len(response.Magerr) {
		return nil, fmt.Errorf("response arrays have different lengths")
	}

	detections := make([]lightcurve.LightcurveObject, len(response.Hmjd))
	for i := range response.Hmjd {
		detections[i] = Detection{
			Oid:      response.Id,
			FilterId: response.FilterId,
			FieldId:  response.FieldId,
			Rcid:     response.Rcid,
			Nepochs:  response.Nepochs,
			ObjRa:    response.ObjRa,
			ObjDec:   response.ObjDec,
			Hmjd:     response.Hmjd[i],
			Mag:      response.Mag[i],
			Magerr:   response.Magerr[i],
		}
	}

	return detections, nil
}
