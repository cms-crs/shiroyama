#!/bin/bash

set -e

INFRA_COMPOSE="docker-compose.infrastructure.yml"
SERVICES_COMPOSE="docker-compose.services.yml"

echo "Инфраструктура готова. Запускаем сервисы..."
docker compose -f "$INFRA_COMPOSE" -f "$SERVICES_COMPOSE" up --build --watch