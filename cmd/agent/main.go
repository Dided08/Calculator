package main

import (
	"log"
	"time"

	"github.com/Dided08/Calculator/internal/agent"
	"github.com/Dided08/Calculator/internal/config"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфиг агента
	cfg, err := config.LoadAgentConfig("internal/config/agent.yaml")
	if err != nil {
		log.Fatalf("ошибка загрузки конфига агента: %v", err)
	}

	// Инициализируем логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	client, err := agent.NewGRPCClient(cfg.OrchestratorAddress)
	if err != nil {
		logger.Fatal("не удалось подключиться к оркестратору", zap.Error(err))
	}

	worker := agent.NewWorker(client, logger)

	ticker := time.NewTicker(time.Duration(cfg.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	logger.Info("Агент запущен", zap.String("address", cfg.OrchestratorAddress))

	for {
		select {
		case <-ticker.C:
			if err := worker.ProcessTask(); err != nil {
				logger.Warn("ошибка обработки задачи", zap.Error(err))
			}
		}
	}
}