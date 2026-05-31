#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-${SCRIPT_DIR}/docker-compose.yml}"
ENV_FILE="${ENV_FILE:-${SCRIPT_DIR}/.env}"
SERVICE_NAME="${SERVICE_NAME:-sub2api-distributor}"

if [ ! -f "${ENV_FILE}" ]; then
  echo "Missing env file: ${ENV_FILE}" >&2
  echo "Create it first from .env.example" >&2
  exit 1
fi

if [ ! -f "${COMPOSE_FILE}" ]; then
  echo "Missing compose file: ${COMPOSE_FILE}" >&2
  exit 1
fi

echo "Deploying ${SERVICE_NAME}"
docker compose --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" up -d --build

echo "Waiting for service health"
sleep 3

"${SCRIPT_DIR}/check.sh"
