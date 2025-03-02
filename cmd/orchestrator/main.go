
package main

import (
    "log"
    "net/http"

    "Calculator/internal/orchestrator"
)

func main() {
    // Инициализация оркестратора
    orch := orchestrator.New()

    // Регистрация обработчиков HTTP-запросов
    http.HandleFunc("/api/v1/calculate", orch.CalculateHandler)
    http.HandleFunc("/api/v1/expressions", orch.ExpressionsHandler)
    http.HandleFunc("/api/v1/expressions/", orch.ExpressionByIDHandler)
    http.HandleFunc("/internal/task", orch.TaskHandler)

    // Запуск HTTP-сервера
    log.Fatal(http.ListenAndServe(":8080", nil))
}