package ztfdr

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	lightcurve := lc.Lightcurve{
		Detections: []lc.LightcurveObject{
			Detection{Oid: 1},
			Detection{Oid: 2},
		},
		NonDetections:    []lc.LightcurveObject{Detection{Oid: 3}},
		ForcedPhotometry: []lc.LightcurveObject{Detection{Oid: 4}},
	}
	objects := []conesearch.MetadataResult{{Catalog: "ztf", Data: []conesearch.MetadataExtended{
		{Metadata: repository.Metadata{ID: "1", Catalog: "ztf"}},
		{Metadata: repository.Metadata{ID: "3", Catalog: "ztf"}},
	}}}

	filtered := Filter(lightcurve, objects)

	require.Equal(t, []lc.LightcurveObject{Detection{Oid: 1}}, filtered.Detections)
	require.Equal(t, lightcurve.NonDetections, filtered.NonDetections)
	require.Equal(t, lightcurve.ForcedPhotometry, filtered.ForcedPhotometry)
}
