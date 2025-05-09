package main

import (
	"net"
	"net/http"

	"github.com/Dided08/Calculator/internal/db"
	"github.com/Dided08/Calculator/internal/models"
	"github.com/Dided08/Calculator/internal/orchestrator"
	"github.com/Dided08/Calculator/internal/parser"
	"github.com/Dided08/Calculator/internal/router"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"github.com/Dided08/Calculator/api" // пакет сгенерирован из agent.proto
)

func main() {
	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Инициализация БД
	dbConn, err := db.InitDB("calculator.db")
	if err != nil {
		sugar.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Инициализация компонентов
	storage := orchestrator.NewStorage(dbConn)
	exprParser := parser.NewParser()
	auth := orchestrator.NewAuth("my-secret-key") // можно вынести в конфиг

	// Запуск HTTP API
	go func() {
		sugar.Infof("HTTP API запущен на :8080")
		httpRouter := router.NewRouter(logger, storage, exprParser, auth)
		if err := http.ListenAndServe(":8080", httpRouter); err != nil {
			sugar.Fatalf("Ошибка HTTP сервера: %v", err)
		}
	}()

	// Запуск gRPC-сервера для агентов
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		sugar.Fatalf("Не удалось запустить gRPC-сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterAgentServiceServer(grpcServer, orchestrator.NewGRPCServer(storage, exprParser))

	sugar.Infof("gRPC сервер запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		sugar.Fatalf("Ошибка gRPC сервера: %v", err)
	}
}