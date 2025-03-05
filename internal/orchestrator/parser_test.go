package orchestrator

import (
    "reflect"
    "testing"

    "github.com/Dided08/Calculator/internal/models"
)

// TestParseExpression проверяет разбор различных выражений
func TestParseExpression(t *testing.T) {
    opTimes := OperationTimes{
        Addition:       1000,
        Subtraction:    1200,
        Multiplication: 1400,
        Division:       1600,
    }
    parser := NewParser(opTimes)

    testCases := []struct {
        name        string
        input       string
        wantTasks   []models.Task
        wantErr     bool
        errorString string
    }{
        {
            name:  "Simple addition",
            input: "2+2",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "2",
                    Arg2:          "2",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
            },
            wantErr: false,
        },
        {
            name:  "Multiplication and division",
            input: "2*(3+4)",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "3",
                    Arg2:          "4",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            2,
                    Arg1:          "2",
                    Arg2:          "res:1",
                    Operation:     models.OperationMultiply,
                    OperationTime: 1400,
                    Dependencies:  []int{1},
                },
            },
            wantErr: false,
        },
        {
            name:  "Complex expression with dependencies",
            input: "(2+3)*(4+5)",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "2",
                    Arg2:          "3",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            2,
                    Arg1:          "4",
                    Arg2:          "5",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            3,
                    Arg1:         "res:1",
                    Arg2:          "res:2",
                    Operation:     models.OperationMultiply,
                    OperationTime: 1400,
                    Dependencies:  []int{1, 2},
                },
            },
            wantErr: false,
        },
        {
            name:  "Invalid expression",
            input: "2++2",
            wantTasks: nil,
            wantErr: true,
            errorString: "некорректный символ: +",
        },
        {
            name:  "Missing closing parenthesis",
            input: "(2+3",
            wantTasks: nil,
            wantErr: true,
            errorString: "ожидалась закрывающая скобка",
        },
    }

    for _, tc := range testCases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            gotTasks, err := parser.ParseExpression(tc.input)

            if !tc.wantErr {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if !reflect.DeepEqual(gotTasks, tc.wantTasks) {
                    t.Errorf("gotTasks = %v, want %v", gotTasks, tc.wantTasks)
                }
            } else {
                if err == nil {
                    t.Errorf("expected error but got none")
                } else if err.Error() != tc.errorString {
                    t.Errorf("got error = %q, want %q", err.Error(), tc.errorString)
                }
            }
        })
    }
}

// TestBuildTasks проверяет создание задач из дерева выражения
func TestBuildTasks(t *testing.T) {
    opTimes := OperationTimes{
        Addition:       1000,
        Subtraction:    1200,
        Multiplication: 1400,
        Division:       1600,
    }
    parser := NewParser(opTimes)

    testCases := []struct {
        name        string
        input       string
        wantTasks   []models.Task
        wantErr     bool
        errorString string
    }{
        {
            name:  "Simple addition",
            input: "2+2",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "2",
                    Arg2:          "2",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
            },
            wantErr: false,
        },
        {
            name:  "Multiplication and division",
            input: "2*(3+4)",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "3",
                    Arg2:          "4",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            2,
                    Arg1:          "2",
                    Arg2:          "res:1",
                    Operation:     models.OperationMultiply,
                    OperationTime: 1400,
                    Dependencies:  []int{1},
                },
            },
            wantErr: false,
        },
        {
            name:  "Complex expression with dependencies",
            input: "(2+3)*(4+5)",
            wantTasks: []models.Task{
                {
                    ID:            1,
                    Arg1:          "2",
                    Arg2:          "3",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            2,
                    Arg1:          "4",
                    Arg2:          "5",
                    Operation:     models.OperationAdd,
                    OperationTime: 1000,
                    Dependencies:  []int{},
                },
                {
                    ID:            3,
                    Arg1:          "res:1",
                    Arg2:          "res:2",
                    Operation:     models.OperationMultiply,
                    OperationTime: 1400,
                    Dependencies:  []int{1, 2},
                },
            },
            wantErr: false,
        },
        {
            name:  "Invalid expression",
            input: "2++2",
            wantTasks: nil,
            wantErr: true,
            errorString: "некорректный символ: +",
        },
        {
            name:  "Missing closing parenthesis",
            input: "(2+3",
            wantTasks: nil,
            wantErr: true,
            errorString: "ожидалась закрывающая скобка",
        },
    }

    for _, tc := range testCases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            gotTasks, err := parser.ParseExpression(tc.input)

            if !tc.wantErr {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if !reflect.DeepEqual(gotTasks, tc.wantTasks) {
                    t.Errorf("gotTasks = %v, want %v", gotTasks, tc.wantTasks)
                }
            } else {
                if err == nil {
                    t.Errorf("expected error but got none")
                } else if err.Error() != tc.errorString {
                    t.Errorf("got error = %q, want %q", err.Error(), tc.errorString)
                }
            }
        })
    }
}