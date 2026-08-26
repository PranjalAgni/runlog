# FIT Extraction Checklist

## 1. File metadata

- [ ] Original filename
- [ ] SHA-256 checksum
- [ ] File size
- [ ] FIT protocol/profile version
- [ ] Device/manufacturer info, if present
- [ ] File creation timestamp
- [ ] Sport/activity type

## 2. Session — overall run summary

- [ ] Start time
- [ ] Sport
- [ ] Total distance
- [ ] Total elapsed time
- [ ] Total timer time
- [ ] Total calories
- [ ] Total ascent
- [ ] Total descent
- [ ] Average speed, if present
- [ ] Maximum speed, if present
- [ ] Average heart rate, if present
- [ ] Maximum heart rate, if present
- [ ] Average cadence, if present
- [ ] Maximum cadence, if present

## 3. Laps / watch-recorded splits

For every lap:

- [ ] Lap number
- [ ] Start time
- [ ] Total distance
- [ ] Elapsed time
- [ ] Timer time
- [ ] Average speed
- [ ] Maximum speed
- [ ] Average heart rate
- [ ] Maximum heart rate
- [ ] Average cadence
- [ ] Total ascent
- [ ] Total descent

## 4. Records — GPS/time-series data

For every valid record:

- [ ] Timestamp
- [ ] Latitude
- [ ] Longitude
- [ ] Elevation / altitude
- [ ] Speed
- [ ] Heart rate
- [ ] Cadence
- [ ] Distance, **only if actually present**

Convert FIT semicircle coordinates into normal degrees before putting them into your domain model.

## 5. Preserve missing values properly

FIT uses sentinel values such as:

- `255` heart rate → missing
- invalid `uint32` distance → missing / `NaN`
- invalid lat/lng → no GPS point

Your normalized data should represent these as **optional/missing**, not real measurements.

## 6. Do not extract these as stored FIT values

These are better calculated by our application:

- [ ] Average pace
- [ ] Split pace
- [ ] Cumulative GPS-derived distance
- [ ] GPS-derived 1 km splits
- [ ] Route polyline
- [ ] First-half / second-half pace
- [ ] Negative/positive split
- [ ] Pace consistency
- [ ] Personal records
- [ ] Heart-rate zones

So the final transformation should be:

`FIT → FileMetadata + Run + Splits + TrackPoints`

with the original FIT file retained as the immutable source of truth.
