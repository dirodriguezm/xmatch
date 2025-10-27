package neowise

import (
	"net/url"
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/utils"
	"github.com/stretchr/testify/require"
)

func Test_addQueryParameters(t *testing.T) {
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
			url, _ := url.Parse(tt.u)
			got := addQueryParameters(url, tt.params)
			require.Equal(t, tt.want, got.String())
		})
	}
}

func Test_convertToLightcurveObject(t *testing.T) {
	tests := []struct {
		name       string
		detections *utils.VOTable
		want       []lightcurve.LightcurveObject
		wantErr    string
	}{
		{"empty VOTable", &utils.VOTable{}, []lightcurve.LightcurveObject{}, "no tables found in response"},
		{"with data", buildFakeVOTable(
			[]string{"ra", "dec", "clon", "clat", "mjd", "w1mpro", "w1sigmpro", "w2mpro", "w2sigmpro", "cntr", "source_id", "dist", "angle"},
			[][]string{
				// mjd, ra, dec, clon, clat, w1mpro, w1sigmpro, w2mpro, w2sigmpro, cntr, source_id, dist, angle
				{"1.0", "2.0", "3.0", "ignore", "ignore", "4.0", "5.0", "6.0", "7.0", "1", "id1", "ignore", "ignore"},
				{"1.0", "2.0", "3.0", "ignore", "ignore", "4.0", "5.0", "6.0", "7.0", "2", "id2", "ignore", "ignore"},
			},
		), []lightcurve.LightcurveObject{
			// This is the expected parsed lightcurve from the VOTable above.
			Detection{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 1, "id1"},
			Detection{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 2, "id2"},
		}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToLightcurveObject(tt.detections)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func buildFakeVOTable(fields []string, data [][]string) *utils.VOTable {
	tableFields := make([]utils.Field, len(fields))
	for i, field := range fields {
		tableFields[i] = utils.Field{Name: field}
	}

	tableData := utils.Data{
		TableData: utils.TableData{
			Rows: make([]utils.Row, len(data)),
		},
	}

	for i, row := range data {
		tableData.TableData.Rows[i] = utils.Row{Columns: make([]utils.Column, len(row))}
		for j, value := range row {
			tableData.TableData.Rows[i].Columns[j] = utils.Column{Value: value}
		}
	}

	return &utils.VOTable{
		Resource: utils.Resource{
			Tables: []utils.Table{
				{
					Fields: tableFields,
					Data:   tableData,
				},
			},
		},
	}
}
