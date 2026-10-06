package domain

import "time"

type Split struct {
	Lap            int
	DistanceMeters float64
	ElapsedTime    time.Duration
	TimerTime      time.Duration
}
