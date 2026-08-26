package domain

import "time"

type Run struct {
	Id             string
	StartedAt      time.Time
	DistanceMeters float64
	ElapsedTime    time.Duration
	TimerTime      time.Duration
	Calories       int
	TotalAscent    float64
	TotalDescent   float64

	Split      []Split
	TrackPoint []TrackPoint
}
