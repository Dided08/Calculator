package calculator

import (
    "strconv"
    "strings"
)

// StrSpace удаляет пробелы из строки и возвращает массив символов.
func StrSpace(expr string) []string {
    var str []string
    expr = strings.ReplaceAll(expr, " ", "")
    for _, char := range expr {
        str = append(str, string(char))
    }
    return str
}

// IsNum проверяет, является ли строка числом.
func IsNum(num string) bool {
    _, err := strconv.ParseFloat(num, 64)
    return err == nil
}

// IsOp проверяет, является ли строка оператором (+, -, *, /).
func IsOp(op string) bool {
    return op == "+" || op == "-" || op == "*" || op == "/"
}

// Priority определяет приоритет оператора.
func Priority(op string) int {
    switch op {
    case "+", "-":
        return 1
    case "*", "/":
        return 2
    default:
        return 0
    }
}

// SetPriority преобразует инфиксное выражение в постфиксно
func SetPriority(str []string) ([]string, error) {
    var output []string
    var ops []string

    for _, char := range str {
        if IsNum(char) {
            output = append(output, char)
        } else if char == "(" {
            ops = append(ops, char)
        } else if char == ")" {
            for len(ops) > 0 && ops[len(ops)-1] != "(" {
                output = append(output, ops[len(ops)-1])
                ops = ops[:len(ops)-1]
            }
            if len(ops) == 0 {
                return nil, ErrInvalidParentheses
            }
            ops = ops[:len(ops)-1]
        } else if IsOp(char) {
            for len(ops) > 0 && Priority(ops[len(ops)-1]) >= Priority(char) {
                output = append(output, ops[len(ops)-1])
                ops = ops[:len(ops)-1]
            }
            ops = append(ops, char)
        } else {
            return nil, ErrInvalidCharacter
        }
    }
    for len(ops) > 0 {
        if ops[len(ops)-1] == "(" {
            return nil, ErrInvalidParentheses
        }
        output = append(output, ops[len(ops)-1])
        ops = ops[:len(ops)-1]
    }

    return output, nil
}

// Calculation выполняет вычисление постфиксного выражения.
func Calculation(output []string) (float64, error) {
    var nums []float64

    for _, char := range output {
        if IsNum(char) {
            num, _ := strconv.ParseFloat(char, 64)
            nums = append(nums, num)
        } else if IsOp(char) {
            if len(nums) < 2 {
                return 0, ErrInvalidExpression
            }
            a := nums[len(nums)-2]
            b := nums[len(nums)-1]
            nums = nums[:len(nums)-2]
            switch char {
            case "+":
                nums = append(nums, a+b)
            case "-":
                nums = append(nums, a-b)
            case "*":
                nums = append(nums, a*b)
            case "/":
                if b == 0 {
                    return 0, ErrDivisionByZero
                }
                nums = append(nums, a/b)
            default:
                return 0, ErrInvalidOperator
            }
        } else {
            return 0, ErrInvalidCharacter
        }
    }
    if len(nums) != 1 {
        return 0, ErrInvalidExpression
    }
    return nums[0], nil
}

// Calculate вычисляет значение арифметического выражения.
func Calculate(expression string) (float64, error) {
    expr := StrSpace(expression)
    postfix, err := SetPriority(expr)
    if err != nil {
        return 0, err
    }
    return Calculation(postfix)
}