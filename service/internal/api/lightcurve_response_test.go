package api

import (
	"testing"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/neowise"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
	"github.com/stretchr/testify/require"
)

type unsupportedLightcurveObject struct{}

func (unsupportedLightcurveObject) GetId() string {
	return "unsupported"
}

func (unsupportedLightcurveObject) GetObjectId() string {
	return "unsupported"
}

func (unsupportedLightcurveObject) GetBrightness() float64 {
	return 0
}

func (unsupportedLightcurveObject) GetBrightnessError() float32 {
	return 0
}

func (unsupportedLightcurveObject) GetMjd() float64 {
	return 0
}

func TestNewLightcurveResponse(t *testing.T) {
	response, err := newLightcurveResponse(lightcurve.Lightcurve{
		Detections: []lightcurve.LightcurveObject{
			ztfdr.Detection{Oid: 42, FilterId: 1, Hmjd: 60234.1, Mag: 18.2, Magerr: 0.05},
			neowise.Detection{Mjd: 60321.4, W1mpro: 15.8, W1sigmpro: 0.12, W2mpro: 15.1, W2sigmpro: 0.2, Cntr: 77, Source_id: "neo-77"},
		},
		ForcedPhotometry: []lightcurve.LightcurveObject{
			&ztfdr.Detection{Oid: 7, Hmjd: 60235.5, Mag: 19.4, Magerr: 0.2},
		},
	})
	require.NoError(t, err)

	require.Len(t, response.Detections, 2)
	require.Empty(t, response.NonDetections)
	require.Len(t, response.ForcedPhotometry, 1)

	ztfEntry := response.Detections[0]
	require.Equal(t, "ztf", ztfEntry.Catalog)
	require.Equal(t, "42", ztfEntry.ObjectID)
	require.Equal(t, 60234.1, ztfEntry.Mjd)
	require.Equal(t, 18.2, ztfEntry.Mag)
	require.Equal(t, float32(0.05), ztfEntry.Magerr)
	require.Equal(t, float64(42), ztfEntry.Data["oid"])
	require.Equal(t, 60234.1, ztfEntry.Data["hmjd"])
	require.Equal(t, 18.2, ztfEntry.Data["mag"])

	neowiseEntry := response.Detections[1]
	require.Equal(t, "neowise", neowiseEntry.Catalog)
	require.Equal(t, "77", neowiseEntry.ObjectID)
	require.Equal(t, 60321.4, neowiseEntry.Mjd)
	require.Equal(t, 15.8, neowiseEntry.Mag)
	require.Equal(t, float32(0.12), neowiseEntry.Magerr)
	require.Equal(t, 15.8, neowiseEntry.Data["w1mpro"])
	require.Equal(t, 15.1, neowiseEntry.Data["w2mpro"])
	require.Equal(t, "neo-77", neowiseEntry.Data["source_id"])

	forcedEntry := response.ForcedPhotometry[0]
	require.Equal(t, "ztf", forcedEntry.Catalog)
	require.Equal(t, "7", forcedEntry.ObjectID)
}

func TestNewLightcurveResponse_UnsupportedObject(t *testing.T) {
	_, err := newLightcurveResponse(lightcurve.Lightcurve{
		Detections: []lightcurve.LightcurveObject{unsupportedLightcurveObject{}},
	})
	require.EqualError(t, err, "unsupported lightcurve object type api.unsupportedLightcurveObject")
}
