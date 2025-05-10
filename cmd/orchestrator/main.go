package main

import (
	"log"
	"net"
	"net/http"

	"github.com/Dided08/Calculator/db"
	"github.com/Dided08/Calculator/internal/config"
	"github.com/Dided08/Calculator/internal/middleware"
	"github.com/Dided08/Calculator/internal/orchestrator"
	"github.com/Dided08/Calculator/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadOrchestratorConfig("internal/config/orchestrator.yaml")
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}

	// Инициализируем логгер
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Подключение к базе данных
	dbConn, err := db.NewConnection(cfg.DatabasePath)
	if err != nil {
		logger.Fatal("ошибка подключения к базе данных", zap.Error(err))
	}
	defer dbConn.Close()

	// Применение схемы
	if err := db.Migrate(dbConn); err != nil {
		logger.Fatal("ошибка миграции базы данных", zap.Error(err))
	}

	// Инициализация хранилища и парсера
	storage := orchestrator.NewStorage(dbConn)
	parser := orchestrator.NewParser()

	// Запуск GRPC сервера
	go func() {
		lis, err := net.Listen("tcp", cfg.GRPCAddress)
		if err != nil {
			logger.Fatal("не удалось запустить GRPC listener", zap.Error(err))
		}

		grpcServer := grpc.NewServer()
		proto.RegisterTaskServiceServer(grpcServer, orchestrator.NewGRPCServer(storage))

		logger.Info("GRPC сервер запущен", zap.String("address", cfg.GRPCAddress))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("ошибка запуска GRPC сервера", zap.Error(err))
		}
	}()

	// Запуск HTTP сервера
	server := orchestrator.NewServer(storage, parser)
	mux := server.SetupRoutes()

	loggedHandler := middleware.LoggingMiddleware(logger)(mux)

	logger.Info("HTTP сервер запущен", zap.String("address", cfg.HTTPAddress))
	if err := http.ListenAndServe(cfg.HTTPAddress, loggedHandler); err != nil {
		logger.Fatal("ошибка запуска HTTP сервера", zap.Error(err))
	}
}