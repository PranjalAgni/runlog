-- name: CreateRun :one

INSERT INTO runs (
  source_file_id,
  started_at,
  distance_meters,
  elapsed_time_ms,
  timer_time_ms,
  calories,
  total_ascent_meters,
  total_descent_meters
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id;