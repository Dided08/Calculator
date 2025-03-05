package main

import (
    "context"
    "flag"
    "log"
    "os"
    "time"
	"os/signal"

    "github.com/Dided08/Calculator/internal/orchestrator"
)

func main() {
    // Определение флагов командной строки
    port := flag.Int("port", 8000, "The port the orchestrator will listen on")
    flag.Parse()

    // Конфигурация оркестратора
    orchConfig := orchestrator.Config{
        Port: *port,
    }

    // Создание контекста для контролируемого завершения программы
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Создание нового оркестратора
    orc, err := orchestrator.NewOrchestrator(orchConfig)
    if err != nil {
        log.Fatalf("Failed to create orchestrator: %v", err)
    }

    // Запуск оркестратора
    go orc.Start(ctx)

    // Ждем завершения программы
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)

    <-quit
    log.Println("Received interrupt signal, shutting down gracefully...")

    // Остановка оркестратора
    cancel()

    // Ожидание завершения всех запущенных процессов
    timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
    defer timeoutCancel()

    if err := orc.Shutdown(timeoutCtx); err != nil {
        log.Fatalf("Failed to shutdown orchestrator: %v", err)
    }

    log.Println("Orchestrator shutdown completed.")
}