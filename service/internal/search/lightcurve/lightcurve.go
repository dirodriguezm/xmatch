package lightcurve

// Lightcurve represents astronomical lightcurve data with detections, non-detections, and forced photometry
//
// swagger:model Lightcurve
type Lightcurve struct {
	Detections       []LightcurveObject `json:"detections"`
	NonDetections    []LightcurveObject `json:"non_detections"`
	ForcedPhotometry []LightcurveObject `json:"forced_photometry"`
}

type LightcurveObject interface {
	GetId() string
	GetObjectId() string
	GetBrightness() float64
	GetBrightnessError() float32
	GetMjd() float64
}
