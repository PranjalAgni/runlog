# DB Schema for persistence model

1. source_files

- id
- filename
- checksum
- storage_path
- file_size_bytes (??)
- imported_at

2. runs

- id
- source_file_id
- started_at
- distance_meters
- elapsed_time_ms
- timer_time_ms
- calories
- total_ascent_meters
- total_descent_meters
- created_at

3. splits

- id
- run_id
- lap_number
- distance_meters
- elpased_time_ms
- timer_time_ms

UNIQUE(run_id, lap_number)

4. track_points

- id
- run_id
- sequence_number (??)
- timestamp
- lat
- lon
- elevation_meters
- speed_mps
- heart_rate
- cadence

## Types

IDs BIGINT
timestamps TIMESTAMPTZ
distance/elevation DOUBLE PRECISION
lat/lng DOUBLE PRECISION
durations BIGINT (milliseconds)
calories INTEGER
HR/cadence SMALLINT
checksum TEXT
