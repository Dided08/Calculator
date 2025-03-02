package orchestrator_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"

    "Calculator/internal/orchestrator"
    "Calculator/internal/task"
    "Calculator/pkg/utils/logger"
)

func TestOrchestrator_CalculateHandler(t *testing.T) {
    // Создаем новый оркестр
    mockDB := &mockDatabase{}
    orch := orchestrator.New(mockDB)

    // Подготавливаем запрос
    request := orchestrator.Request{Expression: "2+2"}
    requestJSON, _ := json.Marshal(request)
    req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(requestJSON))

    // Создаем тестовый сервер
    recorder := httptest.NewRecorder()
    orch.CalculateHandler(recorder, req)

    // Проверяем статус-код ответа
    if recorder.Code != http.StatusCreated {
        t.Errorf("Ожидался статус-код %d, получил %d", http.StatusCreated, recorder.Code)
    }

    // Проверяем содержимое ответа
    var response orchestrator.Response
    if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
        t.Fatalf("Ошибка декодирования ответа: %v", err)
    }

    // Проверяем наличие идентификатора
    if response.ID == 0 {
        t.Errorf("Идентификатор должен быть ненулевым")
    }

    // Проверяем, что задача была создана в базе данных
    expectedTask := task.Task{
        ID:         response.ID,
        Expression: request.Expression,
        Status:     task.Pending,
    }
    if !reflect.DeepEqual(expectedTask, mockDB.createdTask) {
        t.Errorf("Задача не была создана правильно. Ожидаемая: %+v, фактическая: %+v", expectedTask, mockDB.createdTask)
    }
}

type mockDatabase struct {
    createdTask task.Task
}

func (m *mockDatabase) CreateTask(task *task.Task) error {
    m.createdTask = *task
    return nil
}