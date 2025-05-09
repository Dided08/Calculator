package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dided08/Calculator/internal/models"
	"github.com/Dided08/Calculator/internal/storage"
	"github.com/Dided08/Calculator/internal/token"
)

type Server struct {
	Storage *storage.Storage
	Parser  *Parser
	Token   *token.Manager
}

func NewServer(storage *storage.Storage, parser *Parser, token *token.Manager) *Server {
	return &Server{
		Storage: storage,
		Parser:  parser,
		Token:   token,
	}
}

// CalculateHandler — создание нового выражения
func (s *Server) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(models.UserIDContextKey).(int)
	if !ok {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req models.ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if req.Expression == "" {
		http.Error(w, "Выражение не может быть пустым", http.StatusBadRequest)
		return
	}

	// Парсинг выражения в задачи
	tasks, err := s.Parser.Parse(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка парсинга: %v", err), http.StatusUnprocessableEntity)
		return
	}

	// Сохраняем выражение и задачи
	exprID, err := s.Storage.CreateExpression(userID, req.Expression, tasks)
	if err != nil {
		http.Error(w, "Ошибка сохранения выражения", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.ExpressionResponse{ID: exprID})
}

// GetExpressionsHandler — список выражений пользователя
func (s *Server) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(models.UserIDContextKey).(int)
	if !ok {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	exprs, err := s.Storage.GetExpressionsByUser(userID)
	if err != nil {
		http.Error(w, "Ошибка получения выражений", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.ExpressionsResponse{Expressions: exprs})
}

// GetExpressionHandler — получить одно выражение по ID
func (s *Server) GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(models.UserIDContextKey).(int)
	if !ok {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Некорректный путь", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	expr, err := s.Storage.GetExpressionByID(userID, id)
	if err != nil {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(models.ExpressionDetailResponse{Expression: expr})
}