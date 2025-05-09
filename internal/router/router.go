package router

import (
	"net/http"

	"github.com/Dided08/Calculator/internal/handler"
	"github.com/Dided08/Calculator/internal/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	exprHandler *handler.ExpressionHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Auth routes (без JWT)
	mux.HandleFunc("/api/v1/register", method("POST", authHandler.RegisterHandler))
	mux.HandleFunc("/api/v1/login", method("POST", authHandler.LoginHandler))

	// Expression routes (c JWT)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/v1/calculate", method("POST", exprHandler.CreateExpressionHandler))
	protectedMux.HandleFunc("/api/v1/expressions", method("GET", exprHandler.ListExpressionsHandler))
	protectedMux.HandleFunc("/api/v1/expressions/", method("GET", exprHandler.GetExpressionHandler)) // выражение по id

	// Оборачиваем защищённые маршруты в middleware
	return middleware.JWTMiddleware(merge(mux, protectedMux))
}

// Обёртка, проверяющая метод запроса
func method(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

// merge объединяет открытые и защищённые маршруты
func merge(open, protected *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// определим, защищён ли маршрут
		if isProtected(r.URL.Path) {
			protected.ServeHTTP(w, r)
		} else {
			open.ServeHTTP(w, r)
		}
	})
}

// список защищённых путей
func isProtected(path string) bool {
	return path == "/api/v1/calculate" ||
		path == "/api/v1/expressions" ||
		len(path) > len("/api/v1/expressions/") && path[:len("/api/v1/expressions/")] == "/api/v1/expressions/"
}