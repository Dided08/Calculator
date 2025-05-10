package config

import (
	"os"
)

type Config struct {
	DatabasePath string // Путь к SQLite базе данных
	GRPCPort     string // Порт, на котором запускается GRPC-сервер
	HTTPPort     string // Порт, на котором запускается HTTP-сервер
}

// Load загружает конфигурацию из переменных окружения с значениями по умолчанию
func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DB_PATH", "calculator.db"),
		GRPCPort:     getEnv("GRPC_PORT", "50051"),
		HTTPPort:     getEnv("HTTP_PORT", "8080"),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}