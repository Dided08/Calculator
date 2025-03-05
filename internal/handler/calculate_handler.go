package handler

import (
    "encoding/json"
    "net/http"

    "github.com/Dided08/Calculator/errors"
    "github.com/Dided08/Calculator/internal/calculator"
)

// CalculateHandler обрабатывает запросы на вычисление арифметических выражений.
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
    // Проверка метода запроса
    if r.Method != http.MethodPost {
        WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.ErrUnsupportedMethod)
        return
    }

    // Чтение тела запроса
    var requestBody struct {
        Expression string `json:"expression"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        WriteErrorResponse(w, http.StatusBadRequest, errors.ErrMalformedJSON)
        return
    }

    // Проверка на наличие поля "expression"
    if requestBody.Expression == "" {
        WriteErrorResponse(w, http.StatusBadRequest, errors.ErrMissingField)
        return
    }

    // Ограничение длины выражения
    if len(requestBody.Expression) > 256 {
        WriteErrorResponse(w, http.StatusBadRequest, errors.ErrTooLongExpression)
        return
    }

    // Вычисление выражения
    result, err := calculator.Calc(requestBody.Expression)
    if err != nil {
        if err.Error() == "division by zero" {
            WriteErrorResponse(w, http.StatusBadRequest, errors.ErrDivisionByZero)
        } else {
            WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidExpression)
        }
        return
    }

    // Отправка успешного ответа
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(struct {
        Result float64 `json:"result"`
    }{Result: result})
}