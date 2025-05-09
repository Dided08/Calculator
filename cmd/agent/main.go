package main

import (
	"context"
	"log"
	"time"

	"github.com/Dided08/Calculator/internal/agent"
	"github.com/Dided08/Calculator/internal/parser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/Dided08/Calculator/api" // пакет сгенерирован из agent.proto
)

func main() {
	// Настройки подключения
	serverAddr := "localhost:50051" // адрес оркестратора

	// Подключение к gRPC-серверу оркестратора
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC серверу: %v", err)
	}
	defer conn.Close()

	// Создание gRPC клиента
	client := api.NewAgentServiceClient(conn)

	// Инициализация парсера
	exprParser := parser.NewParser()

	// Запуск агента
	agent := agent.NewAgent(client, exprParser)

	// Бесконечная регистрация агента и обработка запросов
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := agent.ProcessTask(ctx); err != nil {
			log.Printf("Ошибка при обработке задания: %v", err)
		}
		cancel()

		// Можно добавить небольшой интервал
		time.Sleep(2 * time.Second)
	}
}