package ztfdr

import (
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
)

func Filter(lightcurve lc.Lightcurve, objects []conesearch.MetadataResult) lc.Lightcurve {
	newLightcurve := lc.Lightcurve{}
	ztfIDs := make(map[string]struct{})

	for _, catalog := range objects {
		for _, object := range catalog.Data {
			if catalog.Catalog != "ztf" && object.GetCatalog() != "ztf" {
				continue
			}
			ztfIDs[object.GetId()] = struct{}{}
		}
	}

	for _, detection := range lightcurve.Detections {
		if _, ok := ztfIDs[detection.GetObjectId()]; ok {
			newLightcurve.Detections = append(newLightcurve.Detections, detection)
		}
	}

	for _, nonDetection := range lightcurve.NonDetections {
		newLightcurve.NonDetections = append(newLightcurve.NonDetections, nonDetection)
	}

	for _, forcedPhotometry := range lightcurve.ForcedPhotometry {
		newLightcurve.ForcedPhotometry = append(newLightcurve.ForcedPhotometry, forcedPhotometry)
	}

	return newLightcurve
}
