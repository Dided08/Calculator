package calculator

import (
    "testing"
)

var tests = []struct {
    name     string
    input    string
    expected float64
    isErr    bool
}{
    {"Simple addition", "2+2", 4, false},
    {"Addition with multiplication", "2+2*2", 6, false},
    {"Multiplication and division", "10*2/5", 4, false},
    {"Parentheses", "(2+2)*2", 8, false},
    {"Negative numbers", "-2+2", 0, false},
    {"Division by zero", "10/0", 0, true}, 
    {"Invalid character", "2&2", 0, true}, 
    {"Missing operand", "2+", 0, true},    
    {"Incorrect parentheses", "(2+2))", 0, true}, 
}

func TestCalculate(t *testing.T) {
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Calculate(tt.input)

            if tt.isErr {
                if err == nil {
                    t.Errorf("Expected an error but got none for input '%s'", tt.input)
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error: %v", err)
                }
                if result != tt.expected {
                    t.Errorf("Got %f, expected %f for input '%s'", result, tt.expected, tt.input)
                }
            }
        })
    }
}