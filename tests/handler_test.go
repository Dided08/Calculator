package tests

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Dided08/Calculator/errors"
    "github.com/Dided08/Calculator/internal/handler"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// TestCalculateHandler проверяет обработчик для вычисления арифметических выражений.
func TestCalculateHandler(t *testing.T) {
    tests := []struct {
        name         string
        method       string
        requestBody  string
        expectedCode int
        expectedResp string
    }{
        {
            name:         "Valid expression",
            method:       http.MethodPost,
            requestBody:  `{"expression": "2 + 2"}`,
            expectedCode: http.StatusOK,
            expectedResp: `{"result":4}`,
        },
        {
            name:         "Invalid expression",
            method:       http.MethodPost,
            requestBody:  `{"expression": "2 + "}`,
            expectedCode: http.StatusBadRequest,
            expectedResp: `{"error":"Invalid expression"}`,
        },
        {
            name:         "Empty expression",
            method:       http.MethodPost,
            requestBody:  `{"expression": ""}`,
            expectedCode: http.StatusBadRequest,
            expectedResp: `{"error":"Missing field: expression"}`,
        },
        {
            name:         "Invalid method",
            method:       http.MethodGet,
            requestBody:  "",
            expectedCode: http.StatusMethodNotAllowed,
            expectedResp: `{"error":"Unsupported HTTP method"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Создание запроса
            req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBufferString(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")

            // Создание записи для тестирования
            recorder := httptest.NewRecorder()

            // Вызов обработчика
            handler.CalculateHandler(recorder, req)

            // Проверка кода ответа
            if recorder.Code != tt.expectedCode {
                t.Errorf("Expected status code %d, got %d", tt.expectedCode, recorder.Code)
            }

            // Проверка ответа
            respBody, err := ioutil.ReadAll(recorder.Body)
            if err != nil {
                t.Fatalf("Failed to read response body: %v", err)
            }

            if string(respBody) != tt.expectedResp {
                t.Errorf("Expected response %s, got %s", tt.expectedResp, string(respBody))
            }
        })
    }
}