package router

import (
	"net/http"

	"github.com/Dided08/Calculator/internal/handler"
	"github.com/Dided08/Calculator/internal/middleware"
	"go.uber.org/zap"
)

// NewRouter настраивает маршруты и применяет middleware.
func NewRouter(logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/calculate", handler.CalculateHandler(logger))

	// Применение middleware для логирования запросов
	loggedRouter := middleware.LoggingMiddleware(logger)(mux)

	return loggedRouter
}