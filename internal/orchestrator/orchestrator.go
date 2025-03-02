package orchestrator

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "Calculator/internal/api"
    "Calculator/internal/task"
    "Calculator/pkg/utils/logger"
)

// Orchestrator управляет всеми операциями и взаимодействует с агентами
type Orchestrator struct {
    db          *sql.DB
    logger      *logger.Logger
    mu          sync.RWMutex
    expressions map[int]*api.Expression
    tasks       map[int]*task.Task
    nextID      int
}

// New создает новый объект Orchestrator
func New(db *sql.DB) *Orchestrator {
    return &Orchestrator{
        db:          db,
        logger:      logger.NewLogger(),
        expressions: make(map[int]*api.Expression),
        tasks:       make(map[int]*task.Task),
    }
}

// CalculateHandler обрабатывает запросы на вычисление арифметических выражений
func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
    var request api.Request
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    exprID := o.nextID
    o.nextID++

    // Парсинг выражения и создание задачи
    parsedExpr, err := task.ParseExpression(request.Expression)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Создание задачи и сохранение в базе данных
    newTask := task.NewTask(parsedExpr, exprID)
    if err := o.db.CreateTask(newTask); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Отправка задачи агентам
    go o.sendTaskToAgents(newTask)

    // Формирование ответа клиенту
    response := api.Response{ID: exprID}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// ExpressionsHandler возвращает список всех выражений
func (o *Orchestrator) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
    // Получение всех выражений из базы данных
    expressions, err := o.db.GetAllExpressions()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Формирование ответа
    response := api.ExpressionsResponse{Expressions: expressions}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// ExpressionByIDHandler возвращает информацию о конкретном выражении
func (o *Orchestrator) ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Извлечение параметра ID из URL
    exprID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Получение выражения по ID из базы данных
    expression, err := o.db.GetExpressionByID(exprID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Формирование ответа
    response := api.ExpressionResponse{Expression: expression}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// TaskHandler обрабатывает запросы от агентов на получение задач и прием результатов
func (o *Orchestrator) TaskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // Получить задачу для агента
        task, err := o.db.GetNextTask()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Отдать задачу агенту
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(task)
    case "POST":
        // Принять результат выполнения задачи
        var result api.TaskResult
        if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Обновить состояние задачи в базе данных
        if err := o.db.UpdateTaskResult(result.ID, result.Result); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Завершить выполнение задачи
        w.WriteHeader(http.StatusOK)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

// sendTaskToAgents отправляет задачу агентам
func (o *Orchestrator) sendTaskToAgents(task *task.Task) {
    // Логика отправки задачи агентам
}