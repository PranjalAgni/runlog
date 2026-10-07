package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PranjalAgni/runlog/internal/fit"
	db "github.com/PranjalAgni/runlog/internal/store/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	run, err := fit.Parse("data/raw/zepp.fit")
	if err != nil {
		log.Fatal(err)
	}
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		fmt.Println("Our DB is connected.....")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	// source file, run and splits are written atomically
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// no-op once the tx is committed
	defer tx.Rollback(ctx)

	queries := db.New(pool).WithTx(tx)

	sourceFileId, err := queries.CreateSourceFile(ctx, db.CreateSourceFileParams{
		Filename:      "zepp.fit",
		Checksum:      "abc123checksum",
		StoragePath:   "data/raw/zepp.fit",
		FileSizeBytes: 102400,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("created soruce file: ", sourceFileId)

	params := db.CreateRunParams{
		SourceFileID: sourceFileId,
		StartedAt: pgtype.Timestamptz{
			Time:  run.StartedAt,
			Valid: true,
		},
		DistanceMeters: run.DistanceMeters,
		ElapsedTimeMs:  run.ElapsedTime.Milliseconds(),
		TimerTimeMs:    run.TimerTime.Milliseconds(),

		Calories: pgtype.Int4{
			Int32: int32(run.Calories),
			Valid: true,
		},

		TotalAscentMeters: pgtype.Float8{
			Float64: run.TotalAscent,
			Valid:   true,
		},

		TotalDescentMeters: pgtype.Float8{
			Float64: run.TotalDescent,
			Valid:   true,
		},
	}

	runID, err := queries.CreateRun(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created a run: ", runID)

	for _, split := range run.Split {
		splitID, err := queries.CreateSplit(ctx, db.CreateSplitParams{
			RunID:          runID,
			LapNumber:      int32(split.Lap),
			DistanceMeters: split.DistanceMeters,
			ElapsedTimeMs:  split.ElapsedTime.Milliseconds(),
			TimerTimeMs:    split.TimerTime.Milliseconds(),
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("created split %d: %v\n", split.Lap, splitID)
	}

	// create a slice where we will append all the trackpoints
	// once the list is ready we can send that in one go using pool.CopyFrom
	// COPY FROM copies data from a file to a table (appending the data to whatever is in the table already
	trackPoints := make([]db.CreateTrackPointsParams, 0, len(run.TrackPoint))
	for _, tp := range run.TrackPoint {
		trackPoints = append(trackPoints, db.CreateTrackPointsParams{
			RunID:          runID,
			SequenceNumber: int32(tp.Sequence),
			RecordedAt: pgtype.Timestamptz{
				Time:  tp.Timestamp,
				Valid: true,
			},
			Latitude:        tp.Latitude,
			Longitude:       tp.Longitude,
			ElevationMeters: float8OrNull(tp.Elevation),
			SpeedMps:        float8OrNull(tp.Speed),
			HeartRate:       int2OrNull(tp.HeartRate),
			Cadence:         int2OrNull(tp.Cadence),
		})
	}

	trackPointsInsertCount, err := queries.CreateTrackPoints(ctx, trackPoints)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created %d track points\n", trackPointsInsertCount)

	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}
}

func float8OrNull(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{}
	}

	return pgtype.Float8{Float64: *v, Valid: true}
}

func int2OrNull(v *uint8) pgtype.Int2 {
	if v == nil {
		return pgtype.Int2{}
	}

	return pgtype.Int2{Int16: int16(*v), Valid: true}
}
