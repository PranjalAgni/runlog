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

	queries := db.New(pool)

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

	fmt.Println("created a run: ", runID)
}
