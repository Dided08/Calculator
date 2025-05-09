package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

// InitDB инициализирует SQLite базу данных и применяет схему.
func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %w", err)
	}

	if err := applySchema(db); err != nil {
		return nil, fmt.Errorf("ошибка применения схемы: %w", err)
	}

	return db, nil
}

func applySchema(db *sql.DB) error {
	content, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("не удалось прочитать schema.sql: %w", err)
	}

	// SQLite не поддерживает exec нескольких выражений в одной строке, поэтому нужно разделить
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("ошибка выполнения запроса: %v\n%q", err, query)
		}
	}

	log.Println("Схема БД успешно применена.")
	return nil
}