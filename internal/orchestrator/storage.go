package orchestrator

import (
    "sync"

    "github.com/Dided08/Calculator/internal/models"
)

// Storage представляет интерфейс для работы с хранилищем выражений и задач.
type Storage interface {
    SaveExpression(expression *models.Expression) error
    GetExpression(id string) (*models.Expression, error)
    ListExpressions() ([]*models.Expression, error)

    SaveTask(task *models.Task) error
    GetTask(id string) (*models.Task, error)
    ListTasks() ([]*models.Task, error)
    DeleteTask(id string) error
}

// InMemoryStorage реализует хранение выражений и задач в памяти.
type InMemoryStorage struct {
    expressions sync.Map // Карта для хранения выражений
    tasks       sync.Map // Карта для хранения задач
}

// NewInMemoryStorage создает новое хранилище в памяти.
func NewInMemoryStorage() *InMemoryStorage {
    return &InMemoryStorage{}
}

// SaveExpression сохраняет выражение в хранилище.
func (s *InMemoryStorage) SaveExpression(expression *models.Expression) error {
    s.expressions.Store(expression.ID, expression)
    return nil
}

// GetExpression получает выражение по его идентификатору.
func (s *InMemoryStorage) GetExpression(id string) (*models.Expression, error) {
    value, ok := s.expressions.Load(id)
    if !ok {
        return nil, models.ErrNotFound
    }
    expression, ok := value.(*models.Expression)
    if !ok {
        return nil, models.ErrInvalidType
    }
    return expression, nil
}

// ListExpressions возвращает список всех сохраненных выражений.
func (s *InMemoryStorage) ListExpressions() ([]*models.Expression, error) {
    var expressions []*models.Expression
    s.expressions.Range(func(key, value interface{}) bool {
        expression, ok := value.(*models.Expression)
        if !ok {
            return true
        }
        expressions = append(expressions, expression)
        return true
    })
    return expressions, nil
}

// SaveTask сохраняет задачу в хранилище.
func (s *InMemoryStorage) SaveTask(task *models.Task) error {
    s.tasks.Store(task.ID, task)
    return nil
}

// GetTask получает задачу по её идентификатору.
func (s *InMemoryStorage) GetTask(id string) (*models.Task, error) {
    value, ok := s.tasks.Load(id)
    if !ok {
        return nil, models.ErrNotFound
    }
    task, ok := value.(*models.Task)
    if !ok {
        return nil, models.ErrInvalidType
    }
    return task, nil
}

// ListTasks возвращает список всех сохранённых задач.
func (s *InMemoryStorage) ListTasks() ([]*models.Task, error) {
    var tasks []*models.Task
    s.tasks.Range(func(key, value interface{}) bool {
        task, ok := value.(*models.Task)
        if !ok {
            return true
        }
        tasks = append(tasks, task)
        return true
    })
    return tasks, nil
}

// DeleteTask удаляет задачу из хранилища.
func (s *InMemoryStorage) DeleteTask(id string) error {
    s.tasks.Delete(id)
    return nil
}