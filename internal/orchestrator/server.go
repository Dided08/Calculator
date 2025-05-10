package orchestrator

import (
	"encoding/json"
	"fmt"
	"github.com/Dided08/Calculator/internal/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

// Server представляет HTTP-сервер оркестратора
type Server struct {
	storage StorageInterface
	parser  ExpressionParser
	logger  *zap.Logger
}

// StorageInterface описывает поведение хранилища
type StorageInterface interface {
	AddExpression(userID, expr string) (int, error)
	GetExpression(userID string, exprID int) (models.Expression, error)
	GetAllExpressions(userID string) ([]models.Expression, error)
}

// ExpressionParser разбирает выражения на задачи
type ExpressionParser interface {
	ParseExpression(expression string) ([]models.Task, error)
}

// NewServer создает новый сервер оркестратора
func NewServer(storage StorageInterface, parser ExpressionParser, logger *zap.Logger) *Server {
	return &Server{
		storage: storage,
		parser:  parser,
		logger:  logger,
	}
}

// SetupRoutes настраивает маршруты HTTP-сервера
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", s.withUser(s.handleCalculate))
	mux.HandleFunc("/api/v1/expressions", s.withUser(s.handleGetExpressions))
	mux.HandleFunc("/api/v1/expressions/", s.withUser(s.handleGetExpression))

	return mux
}

// Middleware для получения userID из заголовка
func (s *Server) withUser(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "User ID обязателен", http.StatusUnauthorized)
			return
		}
		next(w, r, userID)
	}
}

// handleCalculate обрабатывает POST-запрос на вычисление выражения
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req models.ExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Expression) == "" {
		http.Error(w, "Выражение не может быть пустым", http.StatusBadRequest)
		return
	}

	// Парсинг
	_, err := s.parser.ParseExpression(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка разбора выражения: %v", err), http.StatusBadRequest)
		return
	}

	// Сохранение выражения
	exprID, err := s.storage.AddExpression(userID, req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сохранения выражения: %v", err), http.StatusInternalServerError)
		return
	}

	// GRPC агент получает задачи напрямую из базы через grpc-интерфейс

	resp := models.ExpressionResponse{ID: exprID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// handleGetExpressions возвращает все выражения пользователя
func (s *Server) handleGetExpressions(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	exprs, err := s.storage.GetAllExpressions(userID)
	if err != nil {
		http.Error(w, "Ошибка получения выражений", http.StatusInternalServerError)
		return
	}

	resp := models.ExpressionsResponse{Expressions: exprs}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// handleGetExpression возвращает конкретное выражение по ID
func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	exprID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID выражения", http.StatusBadRequest)
		return
	}

	expr, err := s.storage.GetExpression(userID, exprID)
	if err != nil {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	resp := models.ExpressionDetailResponse{Expression: expr}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}