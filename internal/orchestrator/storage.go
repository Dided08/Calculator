package orchestrator

import (
	"database/sql"
	"fmt"
	"github.com/Dided08/Calculator/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// Storage реализует StorageInterface с использованием SQLite
type Storage struct {
	db *sql.DB
}

// NewStorage создаёт новое соединение с базой данных и инициализирует таблицы
func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %w", err)
	}

	// Инициализация таблицы
	createTable := `
	CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		expression TEXT NOT NULL
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return &Storage{db: db}, nil
}

// AddExpression сохраняет выражение и возвращает его ID
func (s *Storage) AddExpression(userID, expr string) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO expressions (user_id, expression) VALUES (?, ?)",
		userID, expr,
	)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки выражения: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}
	return int(id), nil
}

// GetExpression возвращает одно выражение по ID и userID
func (s *Storage) GetExpression(userID string, exprID int) (models.Expression, error) {
	row := s.db.QueryRow(
		"SELECT id, expression FROM expressions WHERE id = ? AND user_id = ?",
		exprID, userID,
	)

	var expr models.Expression
	err := row.Scan(&expr.ID, &expr.Expression)
	if err == sql.ErrNoRows {
		return expr, fmt.Errorf("выражение не найдено")
	} else if err != nil {
		return expr, fmt.Errorf("ошибка чтения выражения: %w", err)
	}
	return expr, nil
}

// GetAllExpressions возвращает все выражения пользователя
func (s *Storage) GetAllExpressions(userID string) ([]models.Expression, error) {
	rows, err := s.db.Query(
		"SELECT id, expression FROM expressions WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса выражений: %w", err)
	}
	defer rows.Close()

	var expressions []models.Expression
	for rows.Next() {
		var expr models.Expression
		if err := rows.Scan(&expr.ID, &expr.Expression); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}