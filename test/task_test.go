package task_test

import (
    "testing"

    "Calculator/internal/task"
)

func TestParseExpression(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        want     []*task.Token
        wantErr  bool
    }{
        {"Простое выражение", "2+2", []*task.Token{{Type: task.TokenNumber, Value: 2.0}, {Type: task.TokenOperator, Value: "+"}, {Type: task.TokenNumber, Value: 2.0}}, false},
        {"Выражение с дробями", "2.5*3.7", []*task.Token{{Type: task.TokenNumber, Value: 2.5}, {Type: task.TokenOperator, Value: "*"}, {Type: task.TokenNumber, Value: 3.7}}, false},
        {"Неверное выражение", "2+", nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := task.ParseExpression(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseExpression() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseExpression() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestExecuteTask(t *testing.T) {
    tests := []struct {
        name     string
        operation string
        args     []float64
        want     float64
        wantErr  bool
    }{
        {"Сложение", "+", []float64{2, 2}, 4, false},
        {"Умножение", "*", []float64{3, 4}, 12, false},
        {"Деление на ноль", "/", []float64{5, 0}, 0, true},
        {"Неизвестная операция", "$", []float64{1, 2}, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := task.ExecuteTask(tt.operation, tt.args)
            if (err != nil) != tt.wantErr {
                t.Errorf("ExecuteTask() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ExecuteTask() got = %v, want %v", got, tt.want)
            }
        })
    }
}
