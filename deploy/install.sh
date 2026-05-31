#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/opt/sub2api-distributor}"
SERVICE_NAME="sub2api-distributor"

echo "Installing sub2api-distributor to ${INSTALL_DIR}"

mkdir -p "${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}/deploy"

cp /Users/lhl/Desktop/code/sub2api-distributor/Dockerfile "${INSTALL_DIR}/Dockerfile"
cp -R /Users/lhl/Desktop/code/sub2api-distributor/backend "${INSTALL_DIR}/backend"
cp -R /Users/lhl/Desktop/code/sub2api-distributor/frontend "${INSTALL_DIR}/frontend"
cp -R /Users/lhl/Desktop/code/sub2api-distributor/deploy/. "${INSTALL_DIR}/deploy"
cp /Users/lhl/Desktop/code/sub2api-distributor/go.mod "${INSTALL_DIR}/go.mod"
cp /Users/lhl/Desktop/code/sub2api-distributor/go.sum "${INSTALL_DIR}/go.sum"

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
