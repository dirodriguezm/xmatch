package lightcurve_test

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
	"github.com/stretchr/testify/require"
)

type stubExternalClient struct {
	result lc.ClientResult
}

func (c stubExternalClient) FetchLightcurve(float64, float64, float64, int) lc.ClientResult {
	return c.result
}

type stubConesearchService struct {
	results []conesearch.MetadataResult
}

func (s stubConesearchService) FindMetadataByConesearch(float64, float64, float64, int, string) ([]conesearch.MetadataResult, error) {
	return s.results, nil
}

type metadataStub struct {
	id      string
	catalog string
}

func (m metadataStub) GetId() string {
	return m.id
}

func (m metadataStub) GetCatalog() string {
	return m.catalog
}

func TestGetLightcurve_AppliesZtfDrFilter(t *testing.T) {
	service, err := lc.New(
		[]lc.ExternalClient{stubExternalClient{
			result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{
				ztfdr.Detection{Oid: 1, Hmjd: 1},
				ztfdr.Detection{Oid: 2, Hmjd: 2},
			}}},
		}},
		[]lc.LightcurveFilter{ztfdr.Filter},
		stubConesearchService{results: []conesearch.MetadataResult{{
			Catalog: "ztf",
			Data:    []conesearch.MetadataExtended{{Metadata: metadataStub{id: "1", catalog: "ztf"}}},
		}}},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10)
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, "1", result.Detections[0].GetObjectId())
}

func TestGetLightcurve_ReturnsClientDataWhenConesearchHasNoMatches(t *testing.T) {
	service, err := lc.New(
		[]lc.ExternalClient{stubExternalClient{
			result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{
				ztfdr.Detection{Oid: 1, Hmjd: 1},
			}}},
		}},
		[]lc.LightcurveFilter{ztfdr.Filter},
		stubConesearchService{},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10)
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, "1", result.Detections[0].GetObjectId())
}
