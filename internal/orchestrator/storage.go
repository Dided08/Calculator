package orchestrator

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Dided08/Calculator/internal/models"
)

type Storage struct {
	DB *sql.DB
}

// NewStorage создает хранилище с подключением к SQLite
func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

// AddExpression сохраняет выражение в базу данных
func (s *Storage) AddExpression(expr string, userID int) (int, error) {
	res, err := s.DB.Exec(`INSERT INTO expressions (raw_expr, status, user_id) VALUES (?, ?, ?)`, expr, models.StatusPending, userID)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки выражения: %w", err)
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// AddTasks сохраняет задачи для выражения
func (s *Storage) AddTasks(exprID int, tasks []models.Task) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, task := range tasks {
		isReady := len(task.Dependencies) == 0
		deps := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(task.Dependencies)), ","), "[]")
		_, err := tx.Exec(`
			INSERT INTO tasks (expression_id, operation, arg1, arg2, dependencies, is_ready)
			VALUES (?, ?, ?, ?, ?, ?)`,
			exprID, task.Operation, task.Arg1, task.Arg2, deps, isReady,
		)
		if err != nil {
			return fmt.Errorf("ошибка вставки задачи: %w", err)
		}
	}

	_, err = tx.Exec(`UPDATE expressions SET status = ? WHERE id = ?`, models.StatusProcessing, exprID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса выражения: %w", err)
	}

	return tx.Commit()
}

// GetAllExpressions возвращает все выражения
func (s *Storage) GetAllExpressions(userID int) ([]models.Expression, error) {
	rows, err := s.DB.Query(`SELECT id, raw_expr, status, result FROM expressions WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []models.Expression
	for rows.Next() {
		var expr models.Expression
		var result sql.NullString
		if err := rows.Scan(&expr.ID, &expr.RawExpr, &expr.Status, &result); err != nil {
			return nil, err
		}
		if result.Valid {
			expr.Result = &result.String
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}

// GetExpression возвращает одно выражение по ID
func (s *Storage) GetExpression(id int, userID int) (models.Expression, error) {
	var expr models.Expression
	var result sql.NullString
	err := s.DB.QueryRow(`SELECT id, raw_expr, status, result FROM expressions WHERE id = ? AND user_id = ?`, id, userID).
		Scan(&expr.ID, &expr.RawExpr, &expr.Status, &result)
	if err == sql.ErrNoRows {
		return expr, errors.New("выражение не найдено")
	} else if err != nil {
		return expr, err
	}

	if result.Valid {
		expr.Result = &result.String
	}
	return expr, nil
}