#!/bin/bash

set -e

COMPOSE_FILES="-f docker-compose.infrastructure.yml \
               -f docker-compose.services.yml"

echo "🔧 Запуск docker-compose up..."

docker-compose $COMPOSE_FILES up --watch