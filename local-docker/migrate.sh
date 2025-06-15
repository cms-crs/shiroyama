#!/bin/bash

set -e

COMPOSE_FILES="-f docker-compose.infrastructure.yml \
               -f docker-compose.migrations.yml"

echo "Запуск docker-compose up..."

docker-compose $COMPOSE_FILES build
docker-compose $COMPOSE_FILES up