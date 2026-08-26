-- +goose Up

CREATE TABLE source_files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  filename TEXT NOT NULL,
  checksum TEXT NOT NULL UNIQUE,
  storage_path TEXT NOT NULL,
  file_size_bytes BIGINT NOT NULL,
  imported_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source_file_id UUID NOT NULL UNIQUE REFERENCES source_files(id),
  started_at TIMESTAMPTZ NOT NULL,
  distance_meters DOUBLE PRECISION NOT NULL,
  elapsed_time_ms BIGINT NOT NULL,
  timer_time_ms BIGINT NOT NULL,
  calories INT,
  total_ascent_meters DOUBLE PRECISION,
  total_descent_meters DOUBLE PRECISION,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE splits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
  lap_number INT NOT NULL,
  distance_meters DOUBLE PRECISION NOT NULL,
  elapsed_time_ms BIGINT NOT NULL,
  timer_time_ms BIGINT NOT NULL,

  UNIQUE (run_id, lap_number)
);

CREATE TABLE track_points (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  run_id UUID NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
  sequence_number INT NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  elevation_meters DOUBLE PRECISION,
  speed_mps DOUBLE PRECISION,
  heart_rate SMALLINT,
  cadence SMALLINT,

  UNIQUE (run_id, sequence_number)
);

-- +goose Down

DROP TABLE track_points;
DROP TABLE splits;
DROP TABLE runs;
DROP TABLE source_files;