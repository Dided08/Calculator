package orchestrator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Dided08/Calculator/internal/models"
)

type node struct {
	op    string
	left  *node
	right *node
	value string
}

var operators = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func ParseExpression(expr string) ([]models.Task, error) {
	tokens := tokenize(expr)
	postfix, err := toPostfix(tokens)
	if err != nil {
		return nil, err
	}
	tree, err := buildTree(postfix)
	if err != nil {
		return nil, err
	}
	tasks := []models.Task{}
	_, err = buildTasks(tree, &tasks, map[string]int{}, 0)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// --- Вспомогательные функции ---

func tokenize(expr string) []string {
	expr = strings.ReplaceAll(expr, "(", " ( ")
	expr = strings.ReplaceAll(expr, ")", " ) ")
	return strings.Fields(expr)
}

func toPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
		case isOperator(token):
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) &&
				operators[token] <= operators[stack[len(stack)-1]] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case token == "(":
			stack = append(stack, token)
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1] != "(" {
				return nil, fmt.Errorf("несовпадающие скобки")
			}
			stack = stack[:len(stack)-1]
		default:
			return nil, fmt.Errorf("недопустимый токен: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("несовпадающие скобки")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output, nil
}

func buildTree(postfix []string) (*node, error) {
	var stack []*node

	for _, token := range postfix {
		if isOperator(token) {
			if len(stack) < 2 {
				return nil, fmt.Errorf("недостаточно аргументов для оператора %s", token)
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, &node{op: token, left: left, right: right})
		} else {
			stack = append(stack, &node{value: token})
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("ошибка построения дерева")
	}
	return stack[0], nil
}

func buildTasks(n *node, tasks *[]models.Task, memo map[string]int, nextID int) (int, error) {
	if n == nil {
		return -1, fmt.Errorf("пустой узел")
	}

	// Если это лист — число
	if n.op == "" {
		if _, ok := memo[n.value]; !ok {
			memo[n.value] = nextID
			*tasks = append(*tasks, models.Task{
				ID:         nextID,
				Operation:  "const",
				Arg1:       n.value,
				Arg2:       "",
				DependsOn:  []int{},
			})
			nextID++
		}
		return memo[n.value], nil
	}

	// Внутренний узел — операция
	leftID, err := buildTasks(n.left, tasks, memo, nextID)
	if err != nil {
		return -1, err
	}
	nextID = maxID(*tasks) + 1
	rightID, err := buildTasks(n.right, tasks, memo, nextID)
	if err != nil {
		return -1, err
	}
	nextID = maxID(*tasks) + 1

	task := models.Task{
		ID:         nextID,
		Operation:  n.op,
		Arg1:       fmt.Sprintf("res%d", leftID),
		Arg2:       fmt.Sprintf("res%d", rightID),
		DependsOn:  []int{leftID, rightID},
	}
	*tasks = append(*tasks, task)
	return nextID, nil
}

func isOperator(s string) bool {
	_, ok := operators[s]
	return ok
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func maxID(tasks []models.Task) int {
	max := -1
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}