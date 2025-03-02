package task

import (
    "fmt"
    "math"
)

// ExecuteTask выполняет операцию над аргументами
func ExecuteTask(operation string, args []float64) (float64, error) {
    switch operation {
    case "+":
        return args[0] + args[1], nil
    case "-":
        return args[0] - args[1], nil
    case "*":
        return args[0] * args[1], nil
    case "/":
        if args[1] == 0 {
            return 0, fmt.Errorf("деление на ноль")
        }
        return args[0] / args[1], nil
    case "^":
        return math.Pow(args[0], args[1]), nil
    default:
        return 0, fmt.Errorf("неизвестная операция: %s", operation)
    }
}