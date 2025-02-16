#!/bin/sh

exec docker compose -f docker-compose.dev.yml up --build --force-recreate --remove-orphans devrunner
