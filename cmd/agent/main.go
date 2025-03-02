package main

import (
    "log"
    "net/http"

    "Calculator/internal/agent"
)

func main() {
    // Инициализация агента
    agt := agent.New()

    // Регистрация обработчика HTTP-запросов
    http.HandleFunc("/internal/task", agt.TaskHandler)

    // Запуск HTTP-сервера
    log.Fatal(http.ListenAndServe(":8001", nil))
}