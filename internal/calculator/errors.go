package calculator

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrInvalidParentheses = errors.New("invalid parentheses")
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidCharacter = errors.New("invalid character")
	ErrDivisionByZero    = errors.New("division by zero")
)