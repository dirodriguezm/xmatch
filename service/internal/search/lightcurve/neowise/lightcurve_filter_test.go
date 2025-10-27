package neowise

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	lightcurve := lc.Lightcurve{
		Detections: []lc.LightcurveObject{&Detection{Cntr: 1}},
	}
	objects := []conesearch.MetadataResult{{Catalog: "allwise", Data: []conesearch.MetadataExtended{
		{Metadata: repository.Allwise{ID: "1", Cntr: 1}},
		{Metadata: repository.Allwise{ID: "2", Cntr: 2}},
	}}}

	filtered := Filter(lightcurve, objects)

	require.Equal(t, []lc.LightcurveObject{&Detection{Cntr: 1}}, filtered.Detections)

}
