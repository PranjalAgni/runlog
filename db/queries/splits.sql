-- name: CreateSplit :one
INSERT INTO splits (
    run_id,
    lap_number,
    distance_meters,
    elapsed_time_ms,
    timer_time_ms
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;