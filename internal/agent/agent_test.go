package agent

import (
	"github.com/Dided08/Calculator/internal/models"
	"testing"
	"time"
)

func TestExecuteTask(t *testing.T) {
	// Создаем агента
	a := NewAgent("http://localhost:8080", 1)

	tests := []struct {
		name      string
		task      models.Task
		want      float64
		wantError bool
	}{
		{
			name: "Сложение",
			task: models.Task{
				ID:            1,
				Arg1:          "2",
				Arg2:          "3",
				Operation:     models.OperationAdd,
				OperationTime: 1, 
			},
			want:      5,
			wantError: false,
		},
		{
			name: "Вычитание",
			task: models.Task{
				ID:            2,
				Arg1:          "5",
				Arg2:          "3",
				Operation:     models.OperationSubtract,
				OperationTime: 1,
			},
			want:      2,
			wantError: false,
		},
		{
			name: "Умножение",
			task: models.Task{
				ID:            3,
				Arg1:          "2",
				Arg2:          "3",
				Operation:     models.OperationMultiply,
				OperationTime: 1,
			},
			want:      6,
			wantError: false,
		},
		{
			name: "Деление",
			task: models.Task{
				ID:            4,
				Arg1:          "6",
				Arg2:          "3",
				Operation:     models.OperationDivide,
				OperationTime: 1,
			},
			want:      2,
			wantError: false,
		},
		{
			name: "Деление на ноль",
			task: models.Task{
				ID:            5,
				Arg1:          "6",
				Arg2:          "0",
				Operation:     models.OperationDivide,
				OperationTime: 1,
			},
			want:      0,
			wantError: true,
		},
		{
			name: "Некорректный аргумент",
			task: models.Task{
				ID:            6,
				Arg1:          "abc",
				Arg2:          "3",
				Operation:     models.OperationAdd,
				OperationTime: 1,
			},
			want:      0, 
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()

			got, err := a.executeTask(&tt.task)

			elapsed := time.Since(start)
			if elapsed < time.Duration(tt.task.OperationTime)*time.Millisecond {
				t.Errorf("executeTask() took %v, want at least %v", elapsed, time.Duration(tt.task.OperationTime)*time.Millisecond)
			}

			// Проверяем наличие ошибки
			if (err != nil) != tt.wantError {
				t.Errorf("executeTask() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				return
			}

			if got != tt.want {
				t.Errorf("executeTask() = %v, want %v", got, tt.want)
			}
		})
	}
}