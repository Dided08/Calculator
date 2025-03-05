package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port     string
	LogLevel string
}

// LoadConfig загружает конфигурацию из файла .env и переменных окружения.
func LoadConfig() (*Config, bool, error) {
	err := godotenv.Load()
	envLoaded := err == nil

	config := &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return config, envLoaded, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию, если переменная не установлена.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
