package orchestrator

import (
	"context"
	"log"
	"net"

	pb "github.com/Dided08/Calculator/proto"
	"github.com/Dided08/Calculator/internal/orchestrator"
	"github.com/Dided08/Calculator/internal/storage"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedCalculatorServer
	Orchestrator *orchestrator.Orchestrator
	Storage      *storage.Storage
}

func NewGRPCServer(orch *orchestrator.Orchestrator, st *storage.Storage) *GRPCServer {
	return &GRPCServer{
		Orchestrator: orch,
		Storage:      st,
	}
}

func (s *GRPCServer) EvaluateExpression(ctx context.Context, req *pb.ExpressionRequest) (*pb.ExpressionResponse, error) {
	userID := req.GetUserId()
	expr := req.GetExpression()

	// Сохраняем выражение
	err := s.Storage.SaveExpression(userID, expr)
	if err != nil {
		return nil, err
	}

	// Распарсить выражение в задачи
	tasks, err := orchestrator.ParseExpression(expr)
	if err != nil {
		return nil, err
	}

	// Выполнить задачи
	result, err := s.Orchestrator.ExecuteTasks(ctx, tasks)
	if err != nil {
		return nil, err
	}

	return &pb.ExpressionResponse{
		Result: result,
	}, nil
}

func StartGRPCServer(addr string, server *GRPCServer) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, server)
	log.Printf("gRPC server started on %s", addr)
	return grpcServer.Serve(listener)
}