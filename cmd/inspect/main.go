package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func findChecksum(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	fmt.Printf("Checksum is %x\n", h.Sum(nil))
}

func convertSemicircleToCoordinate(semicircle int32) {
	degree := float64(semicircle) * 180 / float64(uint64(1)<<31)
	fmt.Printf("Semicircle %d to degree %f\n", semicircle, degree)
}

func inspectActivity(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := decoder.New(f)

	fit, err := dec.Decode()
	if err != nil {
		panic(err)
	}

	fmt.Printf("FileHeader DataSize: %d\n", fit.FileHeader.DataSize)
	fmt.Printf("Messages count: %d\n", len(fit.Messages))

	fmt.Printf("\nMessage: %v\n", fit.Messages[0].Num)
	for _, v := range fit.Messages[0].Fields {
		fmt.Printf("  %s: %v (%v)\n", v.Name, v.Value, v.BaseType)
	}

	activity := filedef.NewActivity(fit.Messages...)
	fmt.Printf("File Type: %s\n", activity.FileId.Type)
	fmt.Printf("Sessions count: %d\n", len(activity.Sessions))
	fmt.Printf("Laps count: %d\n", len(activity.Laps))
	fmt.Printf("Records count: %d\n", len(activity.Records))

	fmt.Printf("Total cals : %d\n ", activity.Sessions[0].TotalCalories)
	fmt.Printf("Sport type : %s\n ", activity.Sessions[0].Sport)
	fmt.Printf("TotalAscent : %v\n ", activity.Sessions[0].TotalAscent)

	fmt.Printf("Total distance raw: %d\n ", activity.Sessions[0].TotalDistance)
	fmt.Printf("Total distance scaled: %.2f m\n", activity.Sessions[0].TotalDistanceScaled())
	fmt.Printf("Average heart rate: %d \n", activity.Sessions[0].AvgHeartRate)

	i := 73
	fmt.Printf("\nSample value of record[%d]:\n", i)
	fmt.Printf("  Distance: %g m\n", activity.Records[i].DistanceScaled())
	convertSemicircleToCoordinate(activity.Records[i].PositionLat)
	convertSemicircleToCoordinate(activity.Records[i].PositionLong)
	fmt.Printf("From the FIT file cordinates in degree %f,%f\n", activity.Records[i].PositionLatDegrees(), activity.Records[i].PositionLongDegrees())
	fmt.Printf("  Speed: %g m/s\n", activity.Records[i].SpeedScaled())
	fmt.Printf("  HeartRate: %d bpm\n", activity.Records[i].HeartRate)
	fmt.Printf("  Cadence: %d rpm\n", activity.Records[i].Cadence)

	for i := 0; i < 1000; i++ {
		dist := activity.Records[i].DistanceScaled()
		if !math.IsNaN(dist) {
			fmt.Println("Distance we got ", dist)
			break
		}
	}

	for i := 0; i < 11; i++ {
		fmt.Printf("Lap %d total distance scaled: %.2f m with time %f and elapsed time %f \n", i+1, activity.Laps[i].TotalDistanceScaled(), activity.Laps[i].TotalTimerTimeScaled(), activity.Laps[i].TotalElapsedTimeScaled())
	}

	// Info needed

	// 1. File metadata: filename, checksum, fileId
	// 2. Session: sport, start/end time, total distance, elasped/timer time, calories, total ascent/descent
	// 3. Laps: lap number, start/end time, elapsed/timer time
	// 4. Records: timestamp, latitude, longitude, altitude/elevation, speed, cadence, heart rate when valid

}

func main() {

	// 1. Create parser for fit files
	// 2. extract needed information
	// 3. log them

	fmt.Println("Hello, World!")

	// accepting file path
	relativeFP := os.Args[1]
	absFileFP, err := filepath.Abs(relativeFP)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(absFileFP)
	findChecksum(absFileFP)
	inspectActivity(absFileFP)
}
