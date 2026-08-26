package domain

import "time"

type TrackPoint struct {
	Timestamp time.Time
	Latitude  float64
	Longitude float64
	Speed     float64
	HeartRate uint8
}
