package db

import (
	"database/sql"
	"fmt"
	"os"
)

// RunMigrations запускает миграции, читая schema.sql и применяя его к БД.
func RunMigrations(db *sql.DB, schemaPath string) error {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла миграции: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("ошибка выполнения миграции: %w", err)
	}

	return nil
}