package api

import (
    "time"
)

// Request описывает входящий запрос на вычисление выражения
type Request struct {
    Expression string `json:"expression"` // Арифметическое выражение
}

// Response описывает ответ на запрос вычисления выражения
type Response struct {
    ID int `json:"id"` // Уникальный идентификатор выражения
}

// ExpressionsResponse описывает ответ на запрос списка всех выражений
type ExpressionsResponse struct {
    Expressions []*Expression `json:"expressions"`
}

// ExpressionResponse описывает ответ на запрос конкретного выражения
type ExpressionResponse struct {
    Expression *Expression `json:"expression"`
}

// Expression представляет одно арифметическое выражение
type Expression struct {
    ID        int       `json:"id"`         // Уникальный идентификатор выражения
    CreatedAt time.Time `json:"created_at"` // Время создания выражения
    Status    string    `json:"status"`     // Текущий статус выражения
    Result    float64   `json:"result"`     // Результат вычисления (если готов)
    Error     string    `json:"error"`      // Сообщение об ошибке (если произошла ошибка)
}

// Task описывает одну задачу для агента
type Task struct {
    ID           int     `json:"id"`           // Идентификатор задачи
    ExpressionID int     `json:"expression_id"` // Идентификатор связанного выражения
    Operation    string  `json:"operation"`     // Операция, которую нужно выполнить
    Args         []float64 `json:"args"`        // Аргументы для операции
}

// TaskResult описывает результат выполненной задачи
type TaskResult struct {
    ID     int     `json:"id"`     // Идентификатор задачи
    Result float64 `json:"result"` // Результат выполнения задачи
}