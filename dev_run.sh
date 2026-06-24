#!/bin/sh
set -eu

[ -f "${PWD}/.env" ] && source "${PWD}/.env"

: ${PROJECT_NAME=$(basename "$PWD")}
: ${PODMAN_NAME=${PROJECT_NAME}_dev}
: ${IMAGE_NAME=${PROJECT_NAME}_dev:latest}
: ${SSH_PORT=2222}

mkdir -p "${PWD}/shared"

podman pod create \
    --name "${PODMAN_NAME}" \
    --userns=keep-id:uid=1000,gid=1000 \
    --replace \
    -p "127.0.0.1:8080:8080" \
    -p "127.0.0.1:${SSH_PORT}:2222"

podman run --rm -d \
    --pod "${PODMAN_NAME}" \
    --name "${PODMAN_NAME}-db" \
    -v "${PODMAN_NAME}-db-data:/var/lib/postgresql" \
    -e POSTGRES_USER="${POSTGRES_USER}" \
    -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
    -e POSTGRES_DB="${POSTGRES_DB}" \
    docker.io/postgres:18-alpine

podman run --rm -d \
    --pod "${PODMAN_NAME}" \
    --name "${PODMAN_NAME}" \
    --passwd=false \
    -v "${NIX_STORE}:/nix/store:ro" \
    -v "${PWD}/shared:/home/runner" \
    -v "${PWD}/src:/srv/app" \
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
    -e APP_ENV="development" \
    -e SSH_PUBKEY="${SSH_PUBKEY}" \
    "$IMAGE_NAME"

trap "podman pod stop ${PODMAN_NAME} && podman pod rm ${PODMAN_NAME}" EXIT

eval "$DEV_CMD"
