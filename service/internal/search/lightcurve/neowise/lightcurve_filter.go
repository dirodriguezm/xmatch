package neowise

import (
	"slices"
	"strconv"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
)

func Filter(lightcurve lc.Lightcurve, objects []conesearch.MetadataResult) lc.Lightcurve {
	newLightcurve := lc.Lightcurve{}

	allCntrs := make([]string, 0)
	for _, catalog := range objects {
		for _, object := range catalog.Data {
			allCntrs = append(allCntrs, strconv.FormatInt(object.Metadata.(repository.Allwise).Cntr, 10))
		}
	}

	for _, detection := range lightcurve.Detections {
		if slices.Contains(allCntrs, detection.GetObjectId()) {
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
