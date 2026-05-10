package lightcurve_test

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
	"github.com/stretchr/testify/require"
)

type stubExternalClient struct {
	result lc.ClientResult
	called *int
}

func (c stubExternalClient) FetchLightcurve(float64, float64, float64, int) lc.ClientResult {
	if c.called != nil {
		(*c.called)++
	}
	return c.result
}

type stubConesearchService struct {
	results       []conesearch.MetadataResult
	catalogTarget *string
}

func (s *stubConesearchService) FindMetadataByConesearch(_, _, _ float64, _ int, catalog string) ([]conesearch.MetadataResult, error) {
	if s.catalogTarget != nil {
		*s.catalogTarget = catalog
	}
	return s.results, nil
}

func TestGetLightcurve_AppliesZtfDrFilter(t *testing.T) {
	service, err := lc.New(
		[]lc.Source{{
			Catalog: "ztf",
			Client: stubExternalClient{result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{
				ztfdr.Detection{Oid: 1, Hmjd: 1},
				ztfdr.Detection{Oid: 2, Hmjd: 2},
			}}}},
			Filter: ztfdr.Filter,
		}},
		&stubConesearchService{results: []conesearch.MetadataResult{{
			Catalog: "ztf",
			Data:    []conesearch.MetadataExtended{{Metadata: repository.Metadata{ID: "1", Catalog: "ztf"}}},
		}}},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10, "ztf")
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, "1", result.Detections[0].GetObjectId())
}

func TestGetLightcurve_ReturnsClientDataWhenConesearchHasNoMatches(t *testing.T) {
	service, err := lc.New(
		[]lc.Source{{
			Catalog: "ztf",
			Client: stubExternalClient{result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{
				ztfdr.Detection{Oid: 1, Hmjd: 1},
			}}}},
			Filter: ztfdr.Filter,
		}},
		&stubConesearchService{},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10, "ztf")
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, "1", result.Detections[0].GetObjectId())
}

func TestGetLightcurve_ExecutesOnlySelectedCatalogClient(t *testing.T) {
	ztfCalls := 0
	neowiseCalls := 0
	conesearchCatalog := ""

	service, err := lc.New(
		[]lc.Source{
			{
				Catalog: "ztf",
				Client: stubExternalClient{
					called: &ztfCalls,
					result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{ztfdr.Detection{Oid: 1, Hmjd: 1}}}},
				},
				Filter: ztfdr.Filter,
			},
			{
				Catalog: "neowise",
				Client:  stubExternalClient{called: &neowiseCalls},
				Filter:  lc.DummyLightcurveFilter,
			},
		},
		&stubConesearchService{
			catalogTarget: &conesearchCatalog,
			results: []conesearch.MetadataResult{{
				Catalog: "ztf",
				Data:    []conesearch.MetadataExtended{{Metadata: repository.Metadata{ID: "1", Catalog: "ztf"}}},
			}},
		},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10, "ztf")
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, 1, ztfCalls)
	require.Equal(t, 0, neowiseCalls)
	require.Equal(t, "ztf", conesearchCatalog)
}

func TestGetLightcurve_AllwiseAliasUsesNeowiseSource(t *testing.T) {
	ztfCalls := 0
	neowiseCalls := 0
	conesearchCatalog := ""

	service, err := lc.New(
		[]lc.Source{
			{
				Catalog: "ztf",
				Client:  stubExternalClient{called: &ztfCalls},
				Filter:  lc.DummyLightcurveFilter,
			},
			{
				Catalog: "neowise",
				Client: stubExternalClient{
					called: &neowiseCalls,
					result: lc.ClientResult{Lightcurve: lc.Lightcurve{}},
				},
				Filter: lc.DummyLightcurveFilter,
			},
		},
		&stubConesearchService{catalogTarget: &conesearchCatalog},
	)
	require.NoError(t, err)

	_, err = service.GetLightcurve(10, -10, 0.2, 10, "allwise")
	require.NoError(t, err)
	require.Equal(t, 0, ztfCalls)
	require.Equal(t, 1, neowiseCalls)
	require.Equal(t, "allwise", conesearchCatalog)
}

func TestGetLightcurve_AllCatalogDoesNotDuplicateFilteredSource(t *testing.T) {
	service, err := lc.New(
		[]lc.Source{
			{
				Catalog: "ztf",
				Client: stubExternalClient{result: lc.ClientResult{Lightcurve: lc.Lightcurve{Detections: []lc.LightcurveObject{
					ztfdr.Detection{Oid: 1, Hmjd: 1},
				}}}},
				Filter: ztfdr.Filter,
			},
			{
				Catalog: "neowise",
				Client:  stubExternalClient{result: lc.ClientResult{Lightcurve: lc.Lightcurve{}}},
				Filter:  lc.DummyLightcurveFilter,
			},
		},
		&stubConesearchService{results: []conesearch.MetadataResult{{
			Catalog: "ztf",
			Data:    []conesearch.MetadataExtended{{Metadata: repository.Metadata{ID: "1", Catalog: "ztf"}}},
		}}},
	)
	require.NoError(t, err)

	result, err := service.GetLightcurve(10, -10, 0.2, 10, "all")
	require.NoError(t, err)
	require.Len(t, result.Detections, 1)
	require.Equal(t, "1", result.Detections[0].GetObjectId())
}
