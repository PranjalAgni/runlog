#!/usr/bin/env bash
# Wipes all data from the local runlog database, keeping the schema intact.
# Usage: ./scripts/reset-db.sh
set -euo pipefail

cd "$(dirname "$0")/.."

if [[ -z "${DATABASE_URL:-}" && -f .env ]]; then
  set -a
  source .env
  set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is not set (and not found in .env)" >&2
  exit 1
fi

# refuse to touch anything that isn't a local database
if [[ ! "$DATABASE_URL" =~ ^postgres(ql)?://([^@/]*@)?(localhost|127\.0\.0\.1|\[::1\])?([:/]|$) ]]; then
  echo "refusing to reset: DATABASE_URL does not point at a local database" >&2
  exit 1
fi

psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -q <<'SQL'
TRUNCATE TABLE
    track_points,
    splits,
    runs,
    source_files;
SQL

echo "reset done: track_points, splits, runs, source_files truncated"
