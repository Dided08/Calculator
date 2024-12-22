package calculator

import (
    "errors"
    "strconv"
    "strings"
)

func StrSpace(expr string) []string {
    var str []string
    expr = strings.ReplaceAll(expr, " ", "")
    for _, char := range expr {
        str = append(str, string(char))
    }
    return str
}

func IsNum(num string) bool {
    _, err := strconv.ParseFloat(num, 64)
    if err == nil {
        return true
    }
    return false
}

func IsOp(op string) bool {
    return op == "+" || op == "-" || op == "*" || op == "/"
}

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
                return nil, errors.New("Invalid parenthesis")
            }
            ops = ops[:len(ops)-1]
        } else if IsOp(char) {
            for len(ops) > 0 && Priority(ops[len(ops)-1]) >= Priority(char) {
                output = append(output, ops[len(ops)-1])
                ops = ops[:len(ops)-1]
            }
            ops = append(ops, char)
        } else {
            return nil, errors.New("Invalid character")
        }
    }
    for len(ops) > 0 {
        if ops[len(ops)-1] == "(" {
            return nil, errors.New("Invalid parenthesis")
        }
        output = append(output, ops[len(ops)-1])
        ops = ops[:len(ops)-1]
    }

    return output, nil
}

func Calculation(output []string) (float64, error) {
    var nums []float64

    for _, char := range output {
        if IsNum(char) {
            num, _ := strconv.ParseFloat(char, 64)
            nums = append(nums, num)
        } else if IsOp(char) {
            if len(nums) < 2 {
                return 0, errors.New("Invalid expression")
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
                    return 0, errors.New("Division by zero")
                }
                nums = append(nums, a/b)
            default:
                return 0, errors.New("Invalid operator")
            }
        } else {
            return 0, errors.New("Invalid character")
        }
    }
    if len(nums) != 1 {
        return 0, errors.New("Invalid expression")
    }
    return nums[0], nil
}

func Calculate(expression string) (float64, error) {
    expr := StrSpace(expression)
    postfix, err := SetPriority(expr)
    if err != nil {
        return 0, err
    }
    return Calculation(postfix)
}