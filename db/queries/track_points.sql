-- name: CreateTrackPoints :copyfrom
INSERT INTO track_points (
    run_id,
    sequence_number,
    recorded_at,
    latitude,
    longitude,
    elevation_meters,
    speed_mps,
    heart_rate,
    cadence
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
