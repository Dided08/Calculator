package tests

import (
    "testing"

    "github.com/Dided08/Calculator/internal/calculator"
)

// TestCalc проверяет корректность вычислений для различных арифметических выражений.
func TestCalc(t *testing.T) {
    tests := []struct {
        name     string
        expression string
        want      float64
        wantErr   bool
    }{
        {
            name:     "Simple addition",
            expression: "2 + 2",
            want:      4,
            wantErr:   false,
        },
        {
            name:     "Simple subtraction",
            expression: "5 - 3",
            want:      2,
            wantErr:   false,
        },
        {
            name:     "Simple multiplication",
            expression: "3 * 4",
            want:      12,
            wantErr:   false,
        },
        {
            name:     "Simple division",
            expression: "8 / 2",
            want:      4,
            wantErr:   false,
        },
        {
            name:     "Complex expression",
            expression: "(2 + 3) * (4 + 5)",
            want:      35,
            wantErr:   false,
        },
        {
            name:     "Division by zero",
            expression: "5 / 0",
            want:      0,
            wantErr:   true,
        },
        {
            name:     "Invalid expression",
            expression: "2 +",
            want:      0,
            wantErr:   true,
        },
        {
            name:     "Missing closing parenthesis",
            expression: "(2 + 3",
            want:      0,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.Calc(tt.expression)

            if (err != nil) != tt.wantErr {
                t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && got != tt.want {
                t.Errorf("Calc() got = %v, want %v", got, tt.want)
            }
        })
    }
}

// TestAdd проверяет корректность выполнения операции сложения.
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        arg1     string
        arg2     string
        want     float64
        wantErr  bool
    }{
        {
            name:    "Valid numbers",
            arg1:    "2",
            arg2:    "3",
            want:    5,
            wantErr: false,
        },
        {
            name:    "Invalid first argument",
            arg1:    "abc",
            arg2:    "3",
            want:    0,
            wantErr: true,
        },
        {
            name:    "Invalid second argument",
            arg1:    "2",
            arg2:    "def",
            want:    0,
            wantErr: true,
        },
        {
            name:    "Both invalid arguments",
            arg1:    "ghi",
            arg2:    "jkl",
            want:    0,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := calculator.Add(tt.arg1, tt.arg2)

            if (err != nil) != tt.wantErr {
                t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && got != tt.want {
                t.Errorf("Add() got = %v, want %v", got, tt.want)
            }
        })
    }
}