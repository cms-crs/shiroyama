#!/bin/bash

set -e

COMPOSE_FILES="-f docker-compose.infrastructure.yml \
               -f docker-compose.migrations.yml \
               -f docker-compose.services.yml"

echo "🔧 Попускаем docker-compose down..."

docker-compose $COMPOSE_FILES down