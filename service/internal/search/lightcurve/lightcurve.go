package lightcurve

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
