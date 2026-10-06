package fit

import (
	"math"
	"time"
)

func scaledOrZero(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}

	return v
}

// scaledOrNil returns the first non-NaN value, or nil if none is set.
func scaledOrNil(values ...float64) *float64 {
	for _, v := range values {
		if !math.IsNaN(v) {
			return &v
		}
	}

	return nil
}

// uint8OrNil returns nil for the FIT invalid sentinel 0xFF.
func uint8OrNil(v uint8) *uint8 {
	if v == math.MaxUint8 {
		return nil
	}

	return &v
}

func secondsToDuration(seconds float64) time.Duration {
	if math.IsNaN(seconds) {
		return 0
	}

	return time.Duration(seconds * float64(time.Second))
}
