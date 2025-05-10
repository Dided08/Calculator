package agent

import (
	"fmt"
	"strconv"
)

// Execute принимает два аргумента и оператор, выполняет вычисление и возвращает результат
func Execute(arg1Str, arg2Str, operator string) (float64, error) {
	arg1, err := strconv.ParseFloat(arg1Str, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать аргумент 1: %w", err)
	}

	arg2, err := strconv.ParseFloat(arg2Str, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать аргумент 2: %w", err)
	}

	switch operator {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("неподдерживаемая операция: %s", operator)
	}
}