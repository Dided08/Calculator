#!/bin/bash

DB_PATH="./data/orchestrator/database.db"
SCHEMA_PATH="./schema.sql"

# Создаём папку, если не существует
mkdir -p "$(dirname "$DB_PATH")"

# Проверка, существует ли БД
if [ ! -f "$DB_PATH" ]; then
  echo "🔧 Инициализация базы данных SQLite..."
  sqlite3 "$DB_PATH" < "$SCHEMA_PATH"
  echo "✅ База данных создана по схеме: $SCHEMA_PATH"
else
  echo "ℹ️ База данных уже существует: $DB_PATH"
fi