#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/opt/sub2api-distributor}"
SERVICE_NAME="sub2api-distributor"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
PROJECT_ROOT_REAL="$(cd "${PROJECT_ROOT}" && pwd -P)"
INSTALL_DIR_REAL="$(mkdir -p "${INSTALL_DIR}" && cd "${INSTALL_DIR}" && pwd -P)"

echo "Installing sub2api-distributor to ${INSTALL_DIR}"

mkdir -p "${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}/deploy"

if [ "${PROJECT_ROOT_REAL}" != "${INSTALL_DIR_REAL}" ]; then
  cp "${PROJECT_ROOT}/Dockerfile" "${INSTALL_DIR}/Dockerfile"
  cp -R "${PROJECT_ROOT}/backend" "${INSTALL_DIR}/backend"
  cp -R "${PROJECT_ROOT}/frontend" "${INSTALL_DIR}/frontend"
  cp -R "${PROJECT_ROOT}/deploy/." "${INSTALL_DIR}/deploy"
  cp "${PROJECT_ROOT}/go.mod" "${INSTALL_DIR}/go.mod"
  cp "${PROJECT_ROOT}/go.sum" "${INSTALL_DIR}/go.sum"
else
  echo "Project root already matches install dir, skipping file copy"
fi

if [ ! -f "${INSTALL_DIR}/deploy/.env" ]; then
  cp "${INSTALL_DIR}/deploy/.env.example" "${INSTALL_DIR}/deploy/.env"
  echo "Created ${INSTALL_DIR}/deploy/.env from template"
fi

cp "${INSTALL_DIR}/deploy/${SERVICE_NAME}.service" "/etc/systemd/system/${SERVICE_NAME}.service"
systemctl daemon-reload

cat <<EOF

Install complete.

Next steps:
1. Edit ${INSTALL_DIR}/deploy/.env
2. Run: systemctl enable --now ${SERVICE_NAME}
3. Check: systemctl status ${SERVICE_NAME}
4. Logs: docker compose -f ${INSTALL_DIR}/deploy/docker-compose.yml logs -f sub2api-distributor

EOF
