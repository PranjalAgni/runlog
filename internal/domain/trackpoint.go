package domain

import "time"

// Optional fields are nil when the device did not record them.
type TrackPoint struct {
	Sequence  int
	Timestamp time.Time
	Latitude  float64
	Longitude float64
	Elevation *float64
	Speed     *float64
	HeartRate *uint8
	Cadence   *uint8
}
