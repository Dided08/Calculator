package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dided08/Calculator/internal/models"
	"github.com/Dided08/Calculator/internal/token"
)

type ExpressionHandler struct {
	Storage ExpressionStorage
	Parser  ExpressionParser
}

type ExpressionStorage interface {
	AddExpression(userID int, expr string) (int, error)
	AddTasks(exprID int, tasks []models.Task) error
	GetAllExpressions(userID int) ([]models.Expression, error)
	GetExpression(userID, exprID int) (models.Expression, error)
}

type ExpressionParser interface {
	ParseExpression(expr string) ([]models.Task, error)
}

// CreateExpressionHandler обрабатывает POST /api/v1/calculate
func (h *ExpressionHandler) CreateExpressionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := token.ExtractUserID(r.Context())
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	var req models.ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusUnprocessableEntity)
		return
	}

	if strings.TrimSpace(req.Expression) == "" {
		http.Error(w, "Выражение не может быть пустым", http.StatusUnprocessableEntity)
		return
	}

	tasks, err := h.Parser.ParseExpression(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка разбора: %v", err), http.StatusUnprocessableEntity)
		return
	}

	exprID, err := h.Storage.AddExpression(userID, req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сохранения выражения: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.Storage.AddTasks(exprID, tasks); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сохранения задач: %v", err), http.StatusInternalServerError)
		return
	}

	resp := models.ExpressionResponse{ID: exprID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// ListExpressionsHandler обрабатывает GET /api/v1/expressions
func (h *ExpressionHandler) ListExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := token.ExtractUserID(r.Context())
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	exprs, err := h.Storage.GetAllExpressions(userID)
	if err != nil {
		http.Error(w, "Ошибка получения выражений", http.StatusInternalServerError)
		return
	}

	resp := models.ExpressionsResponse{Expressions: exprs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetExpressionHandler обрабатывает GET /api/v1/expressions/{id}
func (h *ExpressionHandler) GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := token.ExtractUserID(r.Context())
	if err != nil {
		http.Error(w, "Неавторизован", http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	exprID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	expr, err := h.Storage.GetExpression(userID, exprID)
	if err != nil {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	resp := models.ExpressionDetailResponse{Expression: expr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}