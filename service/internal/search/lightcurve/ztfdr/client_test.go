package ztfdr

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/stretchr/testify/require"
)

func TestAddQueryParameters(t *testing.T) {
	tests := []struct {
		name   string
		u      string
		params map[string]string
		want   string
	}{
		{name: "empty params", u: "https://localhost:8080", params: map[string]string{}, want: "https://localhost:8080"},
		{name: "with params", u: "https://localhost:8080", params: map[string]string{"a": "1", "b": "2"}, want: "https://localhost:8080?a=1&b=2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tt.u)
			require.NoError(t, err)

			got := addQueryParameters(parsedURL, tt.params)

			require.Equal(t, tt.want, got.String())
		})
	}
}

func TestConvertToLightcurveObjects(t *testing.T) {
	tests := []struct {
		name     string
		response lightCurveResponse
		want     []lightcurve.LightcurveObject
		wantErr  string
	}{
		{
			name: "with detections",
			response: lightCurveResponse{
				Id:       123,
				FilterId: 1,
				FieldId:  2,
				Nepochs:  2,
				ObjRa:    3.4,
				ObjDec:   5.6,
				Rcid:     7,
				Hmjd:     []float64{8.9, 9.1},
				Mag:      []float64{10.1, 11.2},
				Magerr:   []float64{0.1, 0.2},
			},
			want: []lightcurve.LightcurveObject{
				Detection{Oid: 123, FilterId: 1, FieldId: 2, Nepochs: 2, ObjRa: 3.4, ObjDec: 5.6, Rcid: 7, Hmjd: 8.9, Mag: 10.1, Magerr: 0.1},
				Detection{Oid: 123, FilterId: 1, FieldId: 2, Nepochs: 2, ObjRa: 3.4, ObjDec: 5.6, Rcid: 7, Hmjd: 9.1, Mag: 11.2, Magerr: 0.2},
			},
		},
		{
			name: "mismatched arrays",
			response: lightCurveResponse{
				Hmjd:   []float64{1},
				Mag:    []float64{1, 2},
				Magerr: []float64{1},
			},
			wantErr: "response arrays have different lengths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToLightcurveObjects(tt.response)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFetchLightcurve(t *testing.T) {
	t.Run("returns detections on success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/light_curve/", r.URL.Path)
			require.Equal(t, "1", r.URL.Query().Get("ra"))
			require.Equal(t, "2", r.URL.Query().Get("dec"))
			require.Equal(t, "3", r.URL.Query().Get("radius"))
			_, err := w.Write([]byte(`{"_id":1,"filterid":2,"fieldid":3,"nepochs":2,"objra":4.5,"objdec":6.7,"rcid":8,"hmjd":[9.1,9.2],"mag":[20.1,20.2],"magerr":[0.1,0.2]}`))
			require.NoError(t, err)
		}))
		defer server.Close()

		client := &ZtfDrClient{url: server.URL + "/light_curve/"}

		result := client.FetchLightcurve(1, 2, 3, 0)

		require.NoError(t, result.Error)
		require.Equal(t, lightcurve.Lightcurve{Detections: []lightcurve.LightcurveObject{
			Detection{Oid: 1, FilterId: 2, FieldId: 3, Nepochs: 2, ObjRa: 4.5, ObjDec: 6.7, Rcid: 8, Hmjd: 9.1, Mag: 20.1, Magerr: 0.1},
			Detection{Oid: 1, FilterId: 2, FieldId: 3, Nepochs: 2, ObjRa: 4.5, ObjDec: 6.7, Rcid: 8, Hmjd: 9.2, Mag: 20.2, Magerr: 0.2},
		}}, result.Lightcurve)
	})

	t.Run("returns empty lightcurve on not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := &ZtfDrClient{url: server.URL + "/light_curve/"}

		result := client.FetchLightcurve(1, 2, 3, 0)

		require.NoError(t, result.Error)
		require.Equal(t, lightcurve.Lightcurve{}, result.Lightcurve)
	})
}
