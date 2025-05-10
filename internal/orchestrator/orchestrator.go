package orchestrator

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dided08/Calculator/internal/middleware"
	"github.com/Dided08/Calculator/internal/orchestrator/grpc_server.go"
	"github.com/Dided08/Calculator/internal/orchestrator/parser.go"
	"github.com/Dided08/Calculator/internal/orchestrator/storage.go"
	"github.com/Dided08/Calculator/internal/orchestrator/server.go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Orchestrator объединяет запуск HTTP и gRPC серверов
type Orchestrator struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	logger     *zap.Logger
}

// New создает новый экземпляр Orchestrator
func New(httpAddr, grpcAddr string, dbPath string) *Orchestrator {
	// Логгер
	logger, _ := zap.NewProduction()

	// Компоненты
	dbStorage, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("не удалось инициализировать БД: %v", err)
	}
	exprParser := parser.NewParser()

	// HTTP сервер
	httpHandler := server.NewServer(dbStorage, exprParser)
	httpRouter := middleware.LoggingMiddleware(logger)(httpHandler.SetupRoutes())
	httpSrv := &http.Server{
		Addr:    httpAddr,
		Handler: httpRouter,
	}

	// gRPC сервер
	grpcSrv := grpc.NewServer()
	grpcTaskServer := grpcserver.NewGRPCServer(dbStorage)
	grpcserver.Register(grpcSrv, grpcTaskServer)

	return &Orchestrator{
		httpServer: httpSrv,
		grpcServer: grpcSrv,
		logger:     logger,
	}
}

// Start запускает HTTP и gRPC серверы параллельно
func (o *Orchestrator) Start(grpcAddr string) {
	var wg sync.WaitGroup

	// gRPC запуск
	wg.Add(1)
	go func() {
		defer wg.Done()
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			o.logger.Fatal("не удалось запустить gRPC сервер", zap.Error(err))
		}
		o.logger.Info("gRPC сервер запущен", zap.String("addr", grpcAddr))
		if err := o.grpcServer.Serve(listener); err != nil {
			o.logger.Fatal("gRPC сервер завершился с ошибкой", zap.Error(err))
		}
	}()

	// HTTP запуск
	wg.Add(1)
	go func() {
		defer wg.Done()
		o.logger.Info("HTTP сервер запущен", zap.String("addr", o.httpServer.Addr))
		if err := o.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			o.logger.Fatal("HTTP сервер завершился с ошибкой", zap.Error(err))
		}
	}()

	// Обработка завершения по Ctrl+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	o.logger.Info("Остановка серверов...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := o.httpServer.Shutdown(ctx); err != nil {
		o.logger.Error("Ошибка при остановке HTTP сервера", zap.Error(err))
	}
	o.grpcServer.GracefulStop()

	wg.Wait()
	o.logger.Info("Все серверы остановлены")
}