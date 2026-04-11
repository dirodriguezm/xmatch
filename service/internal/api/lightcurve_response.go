package api

import (
	"encoding/json"
	"fmt"

	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/neowise"
	"github.com/dirodriguezm/xmatch/service/internal/search/lightcurve/ztfdr"
)

// LightcurveResponse represents the public JSON contract for the lightcurve endpoint.
//
// swagger:model LightcurveResponse
type LightcurveResponse struct {
	Detections       []LightcurveEntry `json:"detections"`
	NonDetections    []LightcurveEntry `json:"non_detections"`
	ForcedPhotometry []LightcurveEntry `json:"forced_photometry"`
}

// LightcurveEntry represents a catalog-aware lightcurve measurement.
//
// swagger:model LightcurveEntry
type LightcurveEntry struct {
	Catalog  string         `json:"catalog"`
	ID       string         `json:"id"`
	ObjectID string         `json:"object_id"`
	Mjd      float64        `json:"mjd"`
	Mag      float64        `json:"mag"`
	Magerr   float32        `json:"magerr"`
	Data     map[string]any `json:"data" swaggertype:"object"`
}

func newLightcurveResponse(lightcurveData lightcurve.Lightcurve) (LightcurveResponse, error) {
	detections, err := newLightcurveEntries(lightcurveData.Detections)
	if err != nil {
		return LightcurveResponse{}, err
	}

	nonDetections, err := newLightcurveEntries(lightcurveData.NonDetections)
	if err != nil {
		return LightcurveResponse{}, err
	}

	forcedPhotometry, err := newLightcurveEntries(lightcurveData.ForcedPhotometry)
	if err != nil {
		return LightcurveResponse{}, err
	}

	return LightcurveResponse{
		Detections:       detections,
		NonDetections:    nonDetections,
		ForcedPhotometry: forcedPhotometry,
	}, nil
}

func newLightcurveEntries(objects []lightcurve.LightcurveObject) ([]LightcurveEntry, error) {
	entries := make([]LightcurveEntry, len(objects))

	for i, object := range objects {
		entry, err := newLightcurveEntry(object)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}

	return entries, nil
}

func newLightcurveEntry(object lightcurve.LightcurveObject) (LightcurveEntry, error) {
	switch detection := object.(type) {
	case ztfdr.Detection:
		return newZtfEntry(detection)
	case *ztfdr.Detection:
		if detection == nil {
			return LightcurveEntry{}, fmt.Errorf("unsupported nil lightcurve object")
		}
		return newZtfEntry(*detection)
	case neowise.Detection:
		return newNeowiseEntry(detection)
	case *neowise.Detection:
		if detection == nil {
			return LightcurveEntry{}, fmt.Errorf("unsupported nil lightcurve object")
		}
		return newNeowiseEntry(*detection)
	default:
		return LightcurveEntry{}, fmt.Errorf("unsupported lightcurve object type %T", object)
	}
}

func newZtfEntry(detection ztfdr.Detection) (LightcurveEntry, error) {
	data, err := dataFromObject(detection)
	if err != nil {
		return LightcurveEntry{}, err
	}

	return LightcurveEntry{
		Catalog:  "ztf",
		ID:       detection.GetId(),
		ObjectID: detection.GetObjectId(),
		Mjd:      detection.GetMjd(),
		Mag:      detection.GetBrightness(),
		Magerr:   detection.GetBrightnessError(),
		Data:     data,
	}, nil
}

func newNeowiseEntry(detection neowise.Detection) (LightcurveEntry, error) {
	data, err := dataFromObject(detection)
	if err != nil {
		return LightcurveEntry{}, err
	}

	return LightcurveEntry{
		Catalog:  "neowise",
		ID:       detection.GetId(),
		ObjectID: detection.GetObjectId(),
		Mjd:      detection.GetMjd(),
		Mag:      detection.GetBrightness(),
		Magerr:   detection.GetBrightnessError(),
		Data:     data,
	}, nil
}

func dataFromObject(object any) (map[string]any, error) {
	payload, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("marshal lightcurve object data: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("unmarshal lightcurve object data: %w", err)
	}

	return data, nil
}
