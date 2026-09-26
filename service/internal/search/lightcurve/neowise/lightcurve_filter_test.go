package neowise

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	t.Run("keeps detections matching allwise cntr", func(t *testing.T) {
		lightcurve := lc.Lightcurve{
			Detections: []lc.LightcurveObject{&Detection{Cntr: 1}},
		}
		objects := []conesearch.MetadataResult{{Catalog: "allwise", Data: []conesearch.MetadataExtended{
			{Metadata: repository.Metadata{ID: "1", Catalog: "allwise", Object: repository.Allwise{ID: "1", Cntr: 1}}},
			{Metadata: repository.Metadata{ID: "2", Catalog: "allwise", Object: repository.Allwise{ID: "2", Cntr: 2}}},
		}}}

		filtered := Filter(lightcurve, objects)

		require.Equal(t, []lc.LightcurveObject{&Detection{Cntr: 1}}, filtered.Detections)
	})

	t.Run("ignores non allwise catalogs", func(t *testing.T) {
		lightcurve := lc.Lightcurve{
			Detections: []lc.LightcurveObject{&Detection{Cntr: 1}},
		}
		objects := []conesearch.MetadataResult{
			{Catalog: "gaia", Data: []conesearch.MetadataExtended{{Metadata: repository.Metadata{ID: "gaia-1", Catalog: "gaia", Object: repository.Gaia{ID: "gaia-1"}}}}},
			{Catalog: "allwise", Data: []conesearch.MetadataExtended{{Metadata: repository.Metadata{ID: "1", Catalog: "allwise", Object: repository.Allwise{ID: "1", Cntr: 1}}}}},
		}

		filtered := Filter(lightcurve, objects)

		require.Equal(t, []lc.LightcurveObject{&Detection{Cntr: 1}}, filtered.Detections)
	})
}
