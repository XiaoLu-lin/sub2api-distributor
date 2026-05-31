#!/bin/sh
set -eu

if [ "${RUN_MIGRATIONS_ON_STARTUP:-true}" = "true" ]; then
  if [ -z "${DATABASE_DSN:-}" ]; then
    echo "DATABASE_DSN is required when RUN_MIGRATIONS_ON_STARTUP=true" >&2
    exit 1
  fi

  for file in /app/migrations/*.sql; do
    if [ -f "$file" ]; then
      echo "Applying migration: $file"
      psql "$DATABASE_DSN" -v ON_ERROR_STOP=1 -f "$file" >/dev/null
    fi
  done
fi

exec /app/sub2api-distributor
