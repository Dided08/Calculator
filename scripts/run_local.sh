#!/bin/bash

set -e

echo "🚀 Запуск локальной среды разработки..."

# Путь до директории проекта
PROJECT_ROOT=$(dirname "$(dirname "$0")")

# Запускаем миграции
"$PROJECT_ROOT/scripts/migrate.sh"

# Запуск оркестратора
echo "🔧 Запуск оркестратора..."
gnome-terminal -- bash -c "cd $PROJECT_ROOT && go run cmd/orchestrator/main.go; exec bash" &

# Небольшая задержка, чтобы GRPC сервер успел запуститься
sleep 2

# Запуск агента
echo "🛠️ Запуск агента..."
gnome-terminal -- bash -c "cd $PROJECT_ROOT && go run cmd/agent/main.go; exec bash" &

echo "✅ Все сервисы запущены!"