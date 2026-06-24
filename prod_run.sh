#!/bin/sh
set -eu

[ -f "${PWD}/.env" ] && source "${PWD}/.env"

: ${PROJECT_NAME=$(basename "$PWD")}
: ${PODMAN_NAME=${PROJECT_NAME}_prod}
: ${IMAGE_NAME=${PROJECT_NAME}_prod:latest}

podman run -d \
    --name "${PODMAN_NAME}" \
    --restart unless-stopped \
    -v "${PODMAN_NAME}_jrdelperu-images:/srv/app/uploads" \
    -e POSTGRESQL_URL="${POSTGRESQL_URL}" \
    -e REL="${REL}" \
    -e SESSION_SECRET="${SESSION_SECRET}" \
    -e SMTP_HOST="${SMTP_HOST}" \
    -e SMTP_PORT="${SMTP_PORT}" \
    -e SMTP_USERNAME="${SMTP_USERNAME}" \
    -e SMTP_PASSWORD="${SMTP_PASSWORD}" \
    -e SMTP_FROM_EMAIL="${SMTP_FROM_EMAIL}" \
    -e SMTP_FROM_NAME="${SMTP_FROM_NAME}" \
    -e SMTP_TO_EMAIL="${SMTP_TO_EMAIL}" \
    -e RECAPTCHA_SITE_KEY="${RECAPTCHA_SITE_KEY}" \
    -e RECAPTCHA_SECRET_KEY="${RECAPTCHA_SECRET_KEY}" \
    -e ADMIN_USERNAME="${ADMIN_USERNAME}" \
    -e ADMIN_EMAIL="${ADMIN_EMAIL}" \
    -e ADMIN_PASSWORD="${ADMIN_PASSWORD}" \
    -e APP_ENV="production" \
    "$@" \
    "$IMAGE_NAME"

