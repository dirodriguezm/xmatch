// Copyright 2024-2026 Diego Rodriguez Mancini
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conesearch

import "math"

// haversineDistance returns the great-circle distance, in arcseconds, between
// two sky positions given in degrees.
func haversineDistance(ra1, dec1, ra2, dec2 float64) float64 {
	ra1Rad := ra1 * math.Pi / 180.0
	dec1Rad := dec1 * math.Pi / 180.0
	ra2Rad := ra2 * math.Pi / 180.0
	dec2Rad := dec2 * math.Pi / 180.0

	deltaRA := ra2Rad - ra1Rad
	deltaDec := dec2Rad - dec1Rad

	a := math.Sin(deltaDec/2)*math.Sin(deltaDec/2) +
		math.Cos(dec1Rad)*math.Cos(dec2Rad)*
			math.Sin(deltaRA/2)*math.Sin(deltaRA/2)

	// Floating point rounding can push a slightly above 1 for (near-)antipodal
	// points, which would make sqrt(1-a) NaN. Clamp it to the valid range.
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * 180.0 / math.Pi * 3600.0
}
