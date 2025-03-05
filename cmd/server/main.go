package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/Dided08/Calculator/config"
    "github.com/Dided08/Calculator/internal/router"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    // Загрузка конфигурации
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Установка уровня логгирования
    var logLevel zapcore.Level
    switch cfg.LogLevel {
    case "debug":
        logLevel = zap.DebugLevel
    case "info":
        logLevel = zap.InfoLevel
    case "warn":
        logLevel = zap.WarnLevel
    case "error":
        logLevel = zap.ErrorLevel
    case "fatal":
        logLevel = zap.FatalLevel
    default:
        logLevel = zap.InfoLevel
    }

    // Инициализация логгера
    loggerCfg := zap.Config{
        Encoding:         "console",
        Level:            zap.NewAtomicLevelAt(logLevel),
        OutputPaths:       []string{"stdout"},
        ErrorOutputPaths:  []string{"stderr"},
        EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
        DisableCaller:     true,
        DisableStacktrace: true,
    }
    logger, err := loggerCfg.Build()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()

    // Установка логгера в глобальную переменную
    zap.ReplaceGlobals(logger)

    // Создание роутера
    router := router.New(cfg)

    // Инициализация HTTP-сервера
    srv := &http.Server{
        Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
        Handler:           router,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       60 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
    }

    // Запуск HTTP-сервера
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Обработка сигналов ОС для аккуратного завершения
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    // Ожидание сигнала завершения
    <-stop

    // Завершение сервера
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Failed to shutdown server: %v", err)
    }

    log.Println("Server shutdown completed.")
}