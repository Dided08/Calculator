package orchestrator

import (
	"fmt"
	"github.com/Dided08/Calculator/internal/models"
	"strings"
)

// Parser обрабатывает арифметические выражения и преобразует их в задачи
type Parser struct {
}

// NewParser создает новый экземпляр парсера
func NewParser() *Parser {
	return &Parser{}
}

// Parse разбивает выражение на задачи, поддерживаются +, -, *, /
func (p *Parser) Parse(expr string) ([]*models.Task, error) {
	tokens := strings.Fields(expr)
	if len(tokens)%2 == 0 {
		return nil, fmt.Errorf("неправильный формат выражения")
	}

	var tasks []*models.Task
	var tempVarCounter int
	var prevResult string

	for i := 0; i < len(tokens); i += 2 {
		if i == 0 {
			// Первый операнд
			prevResult = tokens[i]
			continue
		}

		operator := tokens[i-1]
		arg1 := prevResult
		arg2 := tokens[i]
		resName := fmt.Sprintf("res:%d", tempVarCounter)

		task := &models.Task{
			Operation: operator,
			Arg1:      arg1,
			Arg2:      arg2,
			ResultVar: resName,
		}

		tasks = append(tasks, task)
		prevResult = resName
		tempVarCounter++
	}

	// Установка зависимостей
	for i := 1; i < len(tasks); i++ {
		tasks[i].Dependencies = append(tasks[i].Dependencies, i-1)
	}

	return tasks, nil
}