package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware добавляет логирование всех входящих запросов и ответов.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Логирование запроса
        log.Printf("[%s] %s %s", r.Method, r.URL.Path, start.Format("2006-01-02 15:04:05"))

        // Оборачиваем ResponseWriter для захвата статуса и размера тела ответа
        ww := &responseWriterWrapper{w: w}

        next.ServeHTTP(ww, r)

        // Логирование ответа
        log.Printf("[%s] %s %d %s (%dms)", r.Method, r.URL.Path, ww.status, http.StatusText(ww.status), time.Since(start)/time.Millisecond)
    })
}

// responseWriterWrapper оборачивает стандартный http.ResponseWriter для захвата статуса и размера тела ответа.
type responseWriterWrapper struct {
    w      http.ResponseWriter
    status int
    size   int
}

func (rw *responseWriterWrapper) Header() http.Header {
    return rw.w.Header()
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
    size, err := rw.w.Write(b)
    rw.size += size
    return size, err
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
    rw.status = statusCode
    rw.w.WriteHeader(statusCode)
}