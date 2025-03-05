package orchestrator

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
	"time"

    "github.com/gorilla/mux"
    "github.com/Dided08/Calculator/internal/models"
)

// ServerConfig содержит конфигурационные параметры для сервера.
type ServerConfig struct {
    Port int
}

// Server представляет HTTP-сервер для оркестратора.
type Server struct {
    router     *mux.Router
    config     *ServerConfig
    expressions sync.Map // Хранилище выражений
    tasks      sync.Map // Хранилище задач
}

// NewServer создает новый экземпляр сервера.
func NewServer(config *ServerConfig) *Server {
    s := &Server{
        router: mux.NewRouter(),
        config: config,
    }

    s.routes()

    return s
}

// routes настраивает маршрутизацию для сервера.
func (s *Server) routes() {
    s.router.HandleFunc("/api/v1/calculate", s.handleCalculate).Methods("POST")
    s.router.HandleFunc("/api/v1/expressions", s.handleExpressionsList).Methods("GET")
    s.router.HandleFunc("/api/v1/expressions/{id}", s.handleExpressionByID).Methods("GET")
    s.router.HandleFunc("/internal/task", s.handleGetTask).Methods("GET")
    s.router.HandleFunc("/internal/task", s.handlePostTaskResult).Methods("POST")
}

// handleCalculate обрабатывает запрос на добавление нового выражения для вычисления.
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
    var req models.ExpressionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Генерация уникального идентификатора для выражения
    exprID := generateUniqueID()

    // Создание объекта Expression
    expression := models.Expression{
        ID:         exprID,
        Status:     models.StatusPending,
        Expression: req.Expression,
    }

    // Сохранение выражения в хранилище
    s.expressions.Store(exprID, expression)

    // Формирование ответа
    response := models.ExpressionResponse{
        ID: exprID,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

// handleExpressionsList обрабатывает запрос на получение списка всех выражений.
func (s *Server) handleExpressionsList(w http.ResponseWriter, r *http.Request) {
    var expressions []models.Expression
    s.expressions.Range(func(_, value interface{}) bool {
        expressions = append(expressions, value.(models.Expression))
        return true
    })

    response := models.ExpressionsResponse{
        Expressions: expressions,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

// handleExpressionByID обрабатывает запрос на получение конкретного выражения по его идентификатору.
func (s *Server) handleExpressionByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    exprID := vars["id"]

    var expression models.Expression
    if value, ok := s.expressions.Load(exprID); ok {
        expression = value.(models.Expression)
    } else {
        http.NotFound(w, r)
        return
    }

    response := models.ExpressionDetailResponse{
        Expression: expression,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

// handleGetTask обрабатывает запрос на получение задачи для выполнения.
func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    if value, ok := s.tasks.LoadAndDelete(generateUniqueID()); ok {
        task = value.(models.Task)
    } else {
        http.NotFound(w, r)
        return
    }

    response := models.TaskResponse{
        Task: &task,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

// handlePostTaskResult обрабатывает запрос на прием результата выполнения задачи.
func (s *Server) handlePostTaskResult(w http.ResponseWriter, r *http.Request) {
    var req models.TaskResultRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Обновление статуса выражения и сохранение результата
    if value, ok := s.expressions.Load(req.ID); ok {
        expression := value.(models.Expression)
        expression.Status = models.StatusCompleted
        expression.Result = &req.Result
        s.expressions.Store(req.ID, expression)
    } else {
        http.NotFound(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// generateUniqueID генерирует уникальный идентификатор.
func generateUniqueID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Run запускает HTTP-сервер.
func (s *Server) Run() error {
    log.Printf("Starting server on port %d...", s.config.Port)
    return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}