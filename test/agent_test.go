package agent_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "Calculator/internal/agent"
    "Calculator/internal/api"
    "Calculator/pkg/utils/logger"
)

func TestAgent_TaskHandler_GET(t *testing.T) {
    // Создаем новый агент
    mockDB := &mockDatabase{}
    agt := agent.New(mockDB)

    // Подготавливаем запрос
    req, _ := http.NewRequest("GET", "/internal/task", nil)

    // Создаем тестовый сервер
    recorder := httptest.NewRecorder()
    agt.TaskHandler(recorder, req)

    // Проверяем статус-код ответа
    if recorder.Code != http.StatusOK {
        t.Errorf("Ожидался статус-код %d, получил %d", http.StatusOK, recorder.Code)
    }

    // Проверяем, что задача была получена из базы данных
    if mockDB.getTaskCalled == false {
        t.Errorf("getTask не был вызван")
    }
}

func TestAgent_TaskHandler_POST(t *testing.T) {
    // Создаем новый агент
    mockDB := &mockDatabase{}
    agt := agent.New(mockDB)

    // Подготавливаем запрос
    result := api.TaskResult{ID: 1, Result: 42.0}
    resultJSON, _ := json.Marshal(result)
    req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(resultJSON))

    // Создаем тестовый сервер
    recorder := httptest.NewRecorder()
    agt.TaskHandler(recorder, req)

    // Проверяем статус-код ответа
    if recorder.Code != http.StatusOK {
        t.Errorf("Ожидался статус-код %d, получил %d", http.StatusOK, recorder.Code)
    }

    // Проверяем, что результат был сохранен в базе данных
    if !reflect.DeepEqual(result, mockDB.savedResult) {
        t.Errorf("Результат не был сохранен правильно. Ожидаемый: %+v, фактический: %+v", result, mockDB.savedResult)
    }
}

type mockDatabase struct {
    getTaskCalled bool
    savedResult   api.TaskResult
}

func (m *mockDatabase) GetNextTask() (*task.Task, error) {
    m.getTaskCalled = true
    return &task.Task{ID: 1}, nil
}

func (m *mockDatabase) UpdateTaskResult(ID int, result float64) error {
    m.savedResult = api.TaskResult{ID: ID, Result: result}
    return nil
}