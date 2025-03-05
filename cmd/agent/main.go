package main

import (
    "context"
    "flag"
    "log"
    "os"
    "time"
	"os/signal"

    "github.com/your-project/internal/agent"
)

func main() {
    // Определение флагов командной строки
    port := flag.Int("port", 8080, "The port the agent will listen on")
    orchestratorURL := flag.String("orchestrator-url", "http://localhost:8000", "The URL of the orchestrator")
    flag.Parse()

    // Конфигурация агента
    agentConfig := agent.Config{
        Port:            *port,
        OrchestratorURL: *orchestratorURL,
    }

    // Создание контекста для контролируемого завершения программы
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Создание нового агента
    agt, err := agent.NewAgent(agentConfig)
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }

    // Запуск агента
    go agt.Start(ctx)

    // Ждем завершения программы
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)

    <-quit
    log.Println("Received interrupt signal, shutting down gracefully...")

    // Остановка агента
    cancel()

    // Ожидание завершения всех запущенных процессов
    timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
    defer timeoutCancel()

    if err := agt.Shutdown(timeoutCtx); err != nil {
        log.Fatalf("Failed to shutdown agent: %v", err)
    }

    log.Println("Agent shutdown completed.")
}