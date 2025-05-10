package agent

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Dided08/Calculator/internal/proto"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.TaskServiceClient
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %w", err)
	}

	client := pb.NewTaskServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}, nil
}

func (g *GRPCClient) Close() error {
	return g.conn.Close()
}

// GetTask запрашивает задачу у оркестратора
func (g *GRPCClient) GetTask(ctx context.Context) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.GetTask(ctx, &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задачи: %w", err)
	}
	return resp, nil
}

// SubmitResult отправляет результат вычисления задачи обратно оркестратору
func (g *GRPCClient) SubmitResult(ctx context.Context, taskID int32, result float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := g.client.SubmitResult(ctx, &pb.TaskResult{
		Id:     taskID,
		Result: result,
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке результата: %w", err)
	}
	return nil
}