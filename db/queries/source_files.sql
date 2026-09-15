-- name: CreateSourceFile :one
INSERT INTO source_files (
  filename,
  checksum, 
  storage_path, 
  file_size_bytes
)
VALUES ($1, $2, $3, $4)
RETURNING id;