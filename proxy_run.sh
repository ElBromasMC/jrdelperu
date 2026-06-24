#!/bin/sh
set -eu

[ -f "${PWD}/.env" ] && . "${PWD}/.env"

./prod_run.sh \
    --network postgres-network \
    --network http_network \
    --label "traefik.enable=true" \
    --label "traefik.docker.network=http_network" \
    --label "traefik.http.middlewares.redirect-non-www-to-www.redirectregex.permanent=true" \
    --label "traefik.http.middlewares.redirect-non-www-to-www.redirectregex.regex=^https://${WEBSERVER_HOSTNAME}/(.*)" \
    --label "traefik.http.middlewares.redirect-non-www-to-www.redirectregex.replacement=https://www.${WEBSERVER_HOSTNAME}/\${1}" \
    --label "traefik.http.routers.webserver.entrypoints=websecure" \
    --label "traefik.http.routers.webserver.rule=Host(\`www.${WEBSERVER_HOSTNAME}\`) || Host(\`${WEBSERVER_HOSTNAME}\`)" \
    --label "traefik.http.routers.webserver.middlewares=redirect-non-www-to-www" \
    --label "traefik.http.routers.webserver.tls=true" \
    --label "traefik.http.routers.webserver.tls.certresolver=letsencrypt" \
    --label "traefik.http.routers.webserver.tls.domains[0].main=www.${WEBSERVER_HOSTNAME}" \
    --label "traefik.http.routers.webserver.tls.domains[0].sans=${WEBSERVER_HOSTNAME}" \
    --label "traefik.http.routers.webserver.tls.options=default" \
    --label "traefik.http.services.webserver.loadbalancer.server.port=8080" \
    --label "traefik.http.services.webserver.loadbalancer.server.scheme=http"

