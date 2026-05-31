#!/usr/bin/env bash
set -euo pipefail

HEALTH_URL="${HEALTH_URL:-http://127.0.0.1:8091/health}"
EXPECTED="${EXPECTED:-{\"status\":\"ok\"}}"

echo "Checking health: ${HEALTH_URL}"
BODY="$(curl -fsS "${HEALTH_URL}")"
echo "${BODY}"

if [ "${BODY}" != "${EXPECTED}" ]; then
  echo "Unexpected health response" >&2
  exit 1
fi

echo "Health check passed"
