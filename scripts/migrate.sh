#!/bin/bash

set -e

echo "🔄 Выполнение миграции базы данных..."

# Путь до директории проекта
PROJECT_ROOT=$(dirname "$(dirname "$0")")

# Выполнение миграции
go run "$PROJECT_ROOT/db/migration.go"

echo "✅ Миграция завершена успешно."