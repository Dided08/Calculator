package config

import (
    "fmt"
    "os"
    "strconv"
)

// Config содержит параметры конфигурации приложения.
type Config struct {
    ServerPort      int    // Порт, на котором будет запущен сервер
    OrchestratorURL string // Адрес оркестратора
    AgentPort       int    // Порт, на котором будет запущен агент
}

// LoadConfig загружает конфигурацию из переменных окружения.
func LoadConfig() (*Config, error) {
    serverPort, err := loadIntEnv("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("failed to load SERVER_PORT: %w", err)
    }

    orchestratorURL := loadStringEnv("ORCHESTRATOR_URL", "http://localhost:8000")

    agentPort, err := loadIntEnv("AGENT_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("failed to load AGENT_PORT: %w", err)
    }

    return &Config{
        ServerPort:      serverPort,
        OrchestratorURL: orchestratorURL,
        AgentPort:       agentPort,
    }, nil
}

// loadIntEnv загружает целое число из переменной окружения.
func loadIntEnv(name string, defaultValue int) (int, error) {
    val := os.Getenv(name)
    if val == "" {
        return defaultValue, nil
    }

    intVal, err := strconv.Atoi(val)
    if err != nil {
        return 0, fmt.Errorf("failed to convert %s to integer: %w", name, err)
    }

    return intVal, nil
}

// loadStringEnv загружает строку из переменной окружения.
func loadStringEnv(name string, defaultValue string) string {
    val := os.Getenv(name)
    if val == "" {
        return defaultValue
    }

    return val
}