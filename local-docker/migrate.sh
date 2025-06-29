#!/bin/bash

set -e

COMPOSE_FILES="-f docker-compose.infrastructure.yml \
               -f docker-compose.migrations.yml \
               -f docker-compose.services.yml"

if [ $# -eq 0 ]; then
  SERVICES="user-migrate team-migrate board-migrate"
else
  SERVICES="$@"
fi

echo "Запуск docker-сервисов для миграций: $SERVICES"

docker-compose $COMPOSE_FILES build $SERVICES
docker-compose $COMPOSE_FILES up postgres $SERVICES
