package neowise

import (
	"strconv"
	"strings"

	"github.com/dirodriguezm/xmatch/service/internal/repository"
	"github.com/dirodriguezm/xmatch/service/internal/search/conesearch"
	lc "github.com/dirodriguezm/xmatch/service/internal/search/lightcurve"
)

func Filter(lightcurve lc.Lightcurve, objects []conesearch.MetadataResult) lc.Lightcurve {
	newLightcurve := lc.Lightcurve{}
	allCntrs := make(map[string]struct{})

	for _, catalog := range objects {
		if !strings.EqualFold(catalog.Catalog, "allwise") {
			continue
		}

		for _, object := range catalog.Data {
			allwise, ok := object.Metadata.(repository.Allwise)
			if !ok {
				continue
			}

			allCntrs[strconv.FormatInt(allwise.Cntr, 10)] = struct{}{}
		}
	}

	for _, detection := range lightcurve.Detections {
		if _, ok := allCntrs[detection.GetObjectId()]; ok {
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
