package conesearch

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name            string
		ra1, dec1       float64
		ra2, dec2       float64
		expectedArcsecs float64
		tolerance       float64
	}{
		{
			name: "zero separation",
			ra1:  17.3, dec1: 10,
			ra2: 17.3, dec2: 10,
			expectedArcsecs: 0,
			tolerance:       1e-6,
		},
		{
			name: "one degree along the equator",
			ra1:  0, dec1: 0,
			ra2: 1, dec2: 0,
			expectedArcsecs: 3600,
			tolerance:       1e-6,
		},
		{
			name: "quarter circle along the equator",
			ra1:  0, dec1: 0,
			ra2: 90, dec2: 0,
			expectedArcsecs: 90 * 3600,
			tolerance:       1e-6,
		},
		{
			name: "antipodal points clamp to 180 degrees",
			ra1:  17.3, dec1: 10,
			ra2: 197.3, dec2: -10,
			expectedArcsecs: 180 * 3600,
			tolerance:       1e-6,
		},
		{
			name: "near-antipodal points are finite",
			ra1:  17.3, dec1: 10,
			ra2: 197.3, dec2: -9.999999,
			expectedArcsecs: 180 * 3600,
			tolerance:       0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := haversineDistance(tt.ra1, tt.dec1, tt.ra2, tt.dec2)
			require.Falsef(t, math.IsNaN(distance), "distance was NaN")
			require.InDelta(t, tt.expectedArcsecs, distance, tt.tolerance)
		})
	}
}
