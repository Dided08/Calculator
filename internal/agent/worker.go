package agent

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/Dided08/Calculator/internal/proto"
	"google.golang.org/grpc"
)

// Worker представляет gRPC-клиента для получения и выполнения задач
type Worker struct {
	client pb.TaskServiceClient
}

// NewWorker создает нового воркера и подключается к gRPC-серверу
func NewWorker(grpcAddr string) (*Worker, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := pb.NewTaskServiceClient(conn)
	return &Worker{client: client}, nil
}

// Start запускает воркера в бесконечном цикле
func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Воркeр остановлен")
			return
		default:
			w.processOneTask()
			time.Sleep(1 * time.Second)
		}
	}
}

// processOneTask получает одну задачу и обрабатывает её
func (w *Worker) processOneTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Получаем задачу
	taskResp, err := w.client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		// Обычно это означает, что нет задач
		log.Println("Нет задач или ошибка при получении задачи:", err)
		return
	}

	// Вычисляем результат
	result, err := calculate(taskResp.Task.Arg1, taskResp.Task.Arg2, taskResp.Task.Op)
	if err != nil {
		log.Printf("Ошибка при вычислении задачи %d: %v\n", taskResp.Task.Id, err)
		return
	}

	// Отправляем результат
	_, err = w.client.SubmitResult(ctx, &pb.TaskResult{
		Id:     taskResp.Task.Id,
		Result: result,
	})
	if err != nil {
		log.Printf("Ошибка при отправке результата задачи %d: %v\n", taskResp.Task.Id, err)
	}
}

// calculate выполняет арифметическую операцию над двумя аргументами
func calculate(arg1Str, arg2Str, op string) (float64, error) {
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