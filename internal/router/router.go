package router

import (
	"net/http"

	"github.com/Dided08/Calculator/internal/middleware"
	"github.com/Dided08/Calculator/internal/orchestrator"
	"go.uber.org/zap"
)

// NewRouter создает и настраивает HTTP-роутер оркестратора
func NewRouter(logger *zap.Logger, server *orchestrator.Server) http.Handler {
	mux := http.NewServeMux()

	// Маршруты для пользователей
	mux.HandleFunc("/api/v1/calculate", server.HandleCalculate)
	mux.HandleFunc("/api/v1/expressions", server.HandleGetExpressions)
	mux.HandleFunc("/api/v1/expressions/", server.HandleGetExpression)

	// Применяем middleware логирования
	return middleware.LoggingMiddleware(logger)(mux)
}