package fit

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/PranjalAgni/runlog/internal/domain"
	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func Parse(filePath string) (domain.Run, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return domain.Run{}, fmt.Errorf("open fit file: %w", err)
	}

	defer f.Close()

	fit, err := decoder.New(f).Decode()
	if err != nil {
		return domain.Run{}, fmt.Errorf("decode fit file: %w", err)
	}

	activity := filedef.NewActivity(fit.Messages...)
	if len(activity.Sessions) == 0 {
		return domain.Run{}, fmt.Errorf("fit file has no sessions")
	}

	session := activity.Sessions[0]

	run := domain.Run{
		Id:             filepath.Base(filePath),
		StartedAt:      session.StartTime,
		DistanceMeters: scaledOrZero(session.TotalDistanceScaled()),
		ElapsedTime:    secondsToDuration(session.TotalElapsedTimeScaled()),
		TimerTime:      secondsToDuration(session.TotalTimerTimeScaled()),
	}

	if session.TotalCalories != 0xFFFF {
		run.Calories = int(session.TotalCalories)
	}

	if session.TotalAscent != 0xFFFF {
		run.TotalAscent = float64(session.TotalAscent)
	}

	if session.TotalDescent != 0xFFFF {
		run.TotalDescent = float64(session.TotalDescent)
	}

	// Turning Laps data to Splits
	run.Split = make([]domain.Split, 0, len(activity.Laps))

	for i, lap := range activity.Laps {
		run.Split = append(run.Split, domain.Split{
			Lap:            i + 1,
			DistanceMeters: scaledOrZero(lap.TotalDistanceScaled()),
			ElapsedTime:    secondsToDuration(lap.TotalElapsedTimeScaled()),
			TimerTime:      secondsToDuration(lap.TotalTimerTimeScaled()),
		})
	}

	// turning Records to TrackPoints
	run.TrackPoint = make([]domain.TrackPoint, 0, len(activity.Records))
	for _, rec := range activity.Records {
		lat := rec.PositionLatDegrees()
		lng := rec.PositionLongDegrees()
		if math.IsNaN(lat) || math.IsNaN(lng) {
			continue
		}

		// enhanced fields supersede the legacy ones; some devices only write one of them
		run.TrackPoint = append(run.TrackPoint, domain.TrackPoint{
			Sequence:  len(run.TrackPoint) + 1,
			Timestamp: rec.Timestamp,
			Latitude:  lat,
			Longitude: lng,
			Elevation: scaledOrNil(rec.EnhancedAltitudeScaled(), rec.AltitudeScaled()),
			Speed:     scaledOrNil(rec.EnhancedSpeedScaled(), rec.SpeedScaled()),
			HeartRate: uint8OrNil(rec.HeartRate),
			Cadence:   uint8OrNil(rec.Cadence),
		})
	}

	return run, nil
}
