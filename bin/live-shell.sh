#!/bin/sh

exec docker compose -f docker-compose.dev.yml exec -it devrunner sh
