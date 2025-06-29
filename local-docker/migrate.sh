#!/bin/bash

set -e

COMPOSE_FILES="-f docker-compose.infrastructure.yml \
               -f docker-compose.migrations.yml"

echo "Запуск docker сервисов для миграций..."

docker-compose $COMPOSE_FILES build postgres user-migrate team-migrate
docker-compose $COMPOSE_FILES up postgres user-migrate team-migrate