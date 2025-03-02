package agent

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "Calculator/internal/api"
    "Calculator/internal/task"
    "Calculator/pkg/utils/logger"
)

// Agent обрабатывает задачи, полученные от оркестратора
type Agent struct {
    db          *sql.DB
    logger      *logger.Logger
    mu          sync.RWMutex
    currentTask *task.Task
    orchestratorURL string
}

// New создает новый объект Agent
func New(db *sql.DB, orchestratorURL string) *Agent {
    return &Agent{
        db:          db,
        logger:      logger.NewLogger(),
        orchestratorURL: orchestratorURL,
    }
}

// TaskHandler обрабатывает запросы на получение задач и отправку результатов
func (a *Agent) TaskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // Получить задачу от оркестратора
        task, err := a.db.GetNextTask()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Установить текущую задачу
        a.currentTask = task

        // Вернуть задачу агенту
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
        if err := a.db.UpdateTaskResult(result.ID, result.Result); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Завершить выполнение задачи
        w.WriteHeader(http.StatusOK)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

// Run запускает агента
func (a *Agent) Run() {
    for {
        // Получаем новую задачу от оркестратора
        task, err := a.fetchTask()
        if err != nil {
            a.logger.Error(fmt.Sprintf("Ошибка получения задачи: %v", err))
            time.Sleep(5 * time.Second) // Ждем 5 секунд перед повторной попыткой
            continue
        }

        // Выполняем задачу
        result, err := task.Execute()
        if err != nil {
            a.logger.Error(fmt.Sprintf("Ошибка выполнения задачи: %v", err))
            continue
        }

        // Отправляем результат обратно оркестратору
        if err := a.submitResult(task.ID, result); err != nil {
            a.logger.Error(fmt.Sprintf("Ошибка отправки результата: %v", err))
            continue
        }

        // Задержка перед запросом новой задачи
        time.Sleep(1 * time.Second)
    }
}

// fetchTask получает новую задачу от оркестратора
func (a *Agent) fetchTask() (*task.Task, error) {
    url := fmt.Sprintf("%s/internal/task", a.orchestratorURL)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("не удалось получить задачу: статус-код %d", resp.StatusCode)
    }

    var task task.Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil, err
    }

    return &task, nil
}

// submitResult отправляет результат выполнения задачи обратно оркестратору
func (a *Agent) submitResult(taskID int, result float64) error {
    url := fmt.Sprintf("%s/internal/task", a.orchestratorURL)
    payload := api.TaskResult{
        ID:     taskID,
        Result: result,
    }

    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("не удалось отправить результат: статус-код %d", resp.StatusCode)
    }

    return nil
}