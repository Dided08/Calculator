package orchestrator

import (
	"context"
	"fmt"
	pb "github.com/Dided08/Calculator/proto"
	"google.golang.org/grpc"
	"net"
)

// GRPCServer представляет gRPC-сервер оркестратора
type GRPCServer struct {
	pb.UnimplementedTaskServiceServer
	storage *Storage
}

// NewGRPCServer создает новый gRPC-сервер
func NewGRPCServer(storage *Storage) *GRPCServer {
	return &GRPCServer{
		storage: storage,
	}
}

// Start запускает gRPC-сервер на указанном порту
func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("не удалось запустить gRPC-сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, s)

	fmt.Printf("gRPC-сервер запущен на %s\n", port)
	return grpcServer.Serve(lis)
}

// GetTask отдает агенту следующую готовую задачу
func (s *GRPCServer) GetTask(ctx context.Context, _ *pb.Empty) (*pb.TaskResponse, error) {
	task, err := s.storage.GetReadyTask()
	if err != nil {
		return nil, fmt.Errorf("нет доступных задач: %v", err)
	}

	return &pb.TaskResponse{
		Task: &pb.Task{
			Id:           int32(task.ID),
			ExpressionId: int32(task.ExpressionID),
			Operation:    task.Operation,
			Arg1:         task.Arg1,
			Arg2:         task.Arg2,
		},
	}, nil
}

// SubmitResult принимает результат выполнения задачи от агента
func (s *GRPCServer) SubmitResult(ctx context.Context, req *pb.ResultRequest) (*pb.Empty, error) {
	err := s.storage.UpdateTaskResult(int(req.Id), req.Result)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления результата: %v", err)
	}
	return &pb.Empty{}, nil
}