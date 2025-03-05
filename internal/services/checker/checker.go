package checker

import (
	"Distributed_calculator/internal/models"
	"unicode"
)

type Checker interface {
	CheckParenthesis(expression string) bool
	CheckOperations(expression string) bool
	CheckSymbols(expression string) bool
}
type ExpressionChecker struct{}

// CheckParenthesis проверяет корректность расстановки скобок в выражении.
// Возвращает true, если скобки корректны, иначе false.
func (ExpressionChecker) CheckParenthesis(expression string) bool {
	stack := 0
	for _, char := range expression {
		if char == '(' {
			stack++
		} else if char == ')' {
			stack--
		}
		if stack < 0 {
			return false
		}
	}
	return stack == 0
}

// CheckOperations проверяет корректность операций в выражении.
// Возвращает true, если не обнаружены ошибки, иначе false.
func (ExpressionChecker) CheckOperations(expression string) bool {
	if len(expression) == 0 {
		return false // Пустое выражение считается некорректным
	}

	if IsOperator([]rune(expression)[0]) {
		return false
	} else if IsOperator([]rune(expression)[len([]rune(expression))-1]) {
		return false
	}

	wasOperator := false
	for _, char := range expression {
		if IsOperator(char) {
			if wasOperator {
				return false
			}
			wasOperator = true
		} else {
			wasOperator = false
		}
	}
	return true
}

// CheckSymbols проверяет, что в выражении содержатся только допустимые символы
// (числа, операторы, пробелы и скобки). Возвращает true, если выражение корректно.
func (ExpressionChecker) CheckSymbols(expression string) bool {
	for _, char := range expression {
		if !unicode.IsDigit(char) && !IsOperator(char) && char != '(' && char != ')' && char != '.' {
			return false
		}
	}
	return true
}

func IsOperator(char rune) bool {
	switch char {
	case '+', '-', '*', '/':
		return true
	}
	return false
}

func ValidateExpression(checker Checker, expression string) error {
	if !checker.CheckParenthesis(expression) {
		return models.ErrInvalidParenthesis
	}
	if !checker.CheckOperations(expression) {
		return models.ErrInvalidOperations
	}
	if !checker.CheckSymbols(expression) {
		return models.ErrInvalidSymbols
	}
	return nil
}
