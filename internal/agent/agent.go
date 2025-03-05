package agent

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "strconv"
    "sync"
    "time"
	"bytes"

    "github.com/Dided08/Calculator/internal/models"
)

var (
    errNoTasksAvailable = errors.New("no tasks available")
)

// Agent представляет вычислительный агент.
type Agent struct {
    mux        sync.Mutex
    httpClient *http.Client
    baseURL    string
    computePower int
    runningWorkers chan struct{} // Канал для отслеживания активных рабочих процессов
}

// NewAgent создает новый экземпляр агента.
func NewAgent(baseURL string, computePower int) *Agent {
    a := &Agent{
        httpClient: &http.Client{},
        baseURL:    baseURL,
        computePower: computePower,
        runningWorkers: make(chan struct{}, computePower),
    }

    return a
}

// Start запускает агента.
func (a *Agent) Start(ctx context.Context) {
    rand.Seed(time.Now().UnixNano())

    for i := 0; i < a.computePower; i++ {
        go a.worker(ctx)
    }

    select {
    case <-ctx.Done():
        log.Println("Shutting down agent...")
        close(a.runningWorkers)
    }
}

// worker управляет работой одного рабочего процесса агента.
func (a *Agent) worker(ctx context.Context) {
    for {
        select {
        case a.runningWorkers <- struct{}{}: // Заблокируем канал, пока рабочий процесс активен
            task, err := a.getNextTask()
            if err != nil {
                if err == errNoTasksAvailable {
                    log.Println("No tasks available, waiting for 5 seconds...")
                    time.Sleep(5 * time.Second)
                    break
                }
                log.Printf("Error getting next task: %v\n", err)
                break
            }

            result, err := a.executeTask(task)
            if err != nil {
                log.Printf("Error processing task: %v\n", err)
                break
            }

            err = a.submitResult(result)
            if err != nil {
                log.Printf("Error submitting result: %v\n", err)
                break
            }

            log.Printf("Task processed successfully: ID=%d, Result=%.2f\n", result.ID, result.Result)
        case <-ctx.Done(): // Прекращаем работу, если контекст завершился
            log.Println("Worker stopped due to context cancellation.")
            return
        }
        <-a.runningWorkers // Освобождаем канал после завершения работы
    }
}

// getNextTask запрашивает следующую задачу у оркестратора.
func (a *Agent) getNextTask() (*models.Task, error) {
    url := fmt.Sprintf("%s/internal/task", a.baseURL)

    req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("could not create HTTP request: %w", err)
    }

    resp, err := a.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("could not fetch task: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil, errNoTasksAvailable
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("could not read response body: %w", err)
    }

    var task models.Task
    err = json.Unmarshal(body, &task)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshal task: %w", err)
    }

    return &task, nil
}

// executeTask выполняет задачу и возвращает результат.
func (a *Agent) executeTask(task *models.Task) (*models.TaskResultRequest, error) {
    // Имитируем длительное выполнение операции
    time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

    var result float64
    switch task.Operation {
    case models.OperationAdd:
        result = add(task.Arg1, task.Arg2)
    case models.OperationSubtract:
        result = subtract(task.Arg1, task.Arg2)
    case models.OperationMultiply:
        result = multiply(task.Arg1, task.Arg2)
    case models.OperationDivide:
        result = divide(task.Arg1, task.Arg2)
    default:
        return nil, fmt.Errorf("unsupported operation: %s", task.Operation)
    }

    return &models.TaskResultRequest{
        ID:     task.ID,
        Result: result,
    }, nil
}

// submitResult отправляет результат выполнения задачи оркестратору.
func (a *Agent) submitResult(result *models.TaskResultRequest) error {
    url := fmt.Sprintf("%s/internal/task", a.baseURL)

    payload, err := json.Marshal(result)
    if err != nil {
        return fmt.Errorf("could not marshal task result: %w", err)
    }

    req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(payload))
    if err != nil {
        return fmt.Errorf("could not create HTTP request: %w", err)
    }

    resp, err := a.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("could not submit result: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
    }

    return nil
}

// Вспомогательные функции для выполнения арифметических операций
func add(arg1, arg2 string) float64 {
    f1, _ := strconv.ParseFloat(arg1, 64)
    f2, _ := strconv.ParseFloat(arg2, 64)
    return f1 + f2
}

func subtract(arg1, arg2 string) float64 {
    f1, _ := strconv.ParseFloat(arg1, 64)
    f2, _ := strconv.ParseFloat(arg2, 64)
    return f1 - f2
}

func multiply(arg1, arg2 string) float64 {
    f1, _ := strconv.ParseFloat(arg1, 64)
    f2, _ := strconv.ParseFloat(arg2, 64)
    return f1 * f2
}

func divide(arg1, arg2 string) float64 {
    f1, _ := strconv.ParseFloat(arg1, 64)
    f2, _ := strconv.ParseFloat(arg2, 64)
    if f2 == 0 {
        panic("Division by zero!")
    }
    return f1 / f2
}

// parseArgs извлекает аргументы из строки выражения.
func parseArgs(operation string, args []string) (string, string, error) {
    if len(args) != 2 {
        return "", "", fmt.Errorf("incorrect number of arguments for operation '%s'", operation)
    }
    return args[0], args[1], nil
}