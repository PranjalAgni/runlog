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

func secondsToDuration(seconds float64) time.Duration {
	if math.IsNaN(seconds) {
		return 0
	}

	return time.Duration(seconds * float64(time.Second))
}
