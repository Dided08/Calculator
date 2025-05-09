package agent

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/Dided08/Calculator/api" // путь к сгенерированному .proto
	"google.golang.org/grpc"
)

type AgentClient struct {
	client pb.CalculatorClient
	conn   *grpc.ClientConn
}

// NewAgentClient подключается к gRPC-серверу оркестратора
func NewAgentClient(address string) (*AgentClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}
	client := pb.NewCalculatorClient(conn)
	return &AgentClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close закрывает соединение
func (a *AgentClient) Close() error {
	return a.conn.Close()
}

// Run запускает цикл получения и выполнения задач
func (a *AgentClient) Run() {
	for {
		task, err := a.getTask()
		if err != nil {
			log.Printf("Нет задач: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Получена задача: %v", task)

		result, err := execute(task.Operation, task.Arg1, task.Arg2)
		if err != nil {
			log.Printf("Ошибка выполнения задачи: %v", err)
			continue
		}

		if err := a.submitResult(task.Id, result); err != nil {
			log.Printf("Ошибка отправки результата: %v", err)
			continue
		}

		log.Printf("Результат задачи %d: %f", task.Id, result)
	}
}

// getTask запрашивает задачу у сервера
func (a *AgentClient) getTask() (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := a.client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Task, nil
}

// submitResult отправляет результат выполнения задачи
func (a *AgentClient) submitResult(taskID int32, result float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.client.SubmitResult(ctx, &pb.ResultRequest{
		Id:     taskID,
		Result: result,
	})
	return err
}

// execute выполняет арифметическую операцию
func execute(op, arg1Str, arg2Str string) (float64, error) {
	arg1, err := strconv.ParseFloat(arg1Str, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать аргумент 1: %w", err)
	}
	arg2, err := strconv.ParseFloat(arg2Str, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать аргумент 2: %w", err)
	}

	switch op {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("неподдерживаемая операция: %s", op)
	}
}