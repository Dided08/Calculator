package task

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

// ParseExpression разбивает выражение на токены и готовит его для выполнения
func ParseExpression(expression string) ([]*Token, error) {
    tokens := []*Token{}

    // Удаление пробелов
    expression = strings.ReplaceAll(expression, " ", "")

    // Регулярные выражения для поиска чисел и операторов
    numRegex := regexp.MustCompile(`\d+(?:\.\d+)?(?:e[-+]?\d+)?`)
    opRegex := regexp.MustCompile(`[\+\-\*/$$]`)

    // Найти числа и операторы
    pos := 0
    for pos < len(expression) {
        if match := numRegex.FindStringIndex(expression[pos:]); match != nil {
            start, end := pos+match[0], pos+match[1]
            number, err := strconv.ParseFloat(expression[start:end], 64)
            if err != nil {
                return nil, fmt.Errorf("ошибка парсинга числа: %w", err)
            }
            tokens = append(tokens, &Token{Type: TokenNumber, Value: number})
            pos += match[1]
        } else if match := opRegex.FindStringIndex(expression[pos:]); match != nil {
            start, end := pos+match[0], pos+match[1]
            op := expression[start:end]
            tokens = append(tokens, &Token{Type: TokenOperator, Value: op})
            pos += match[1]
        } else {
            return nil, fmt.Errorf("неизвестный символ в позиции %d: %c", pos, expression[pos])
        }
    }

    return tokens, nil
}

// Token представляет один элемент арифметического выражения
type Token struct {
    Type  TokenType
    Value interface{} // Число или оператор
}

// Типы токенов
type TokenType int

const (
    TokenNumber TokenType = iota
    TokenOperator
)