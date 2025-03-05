package parser

import (
	"Distributed_calculator/internal/services/checker"
	"strings"
	"unicode"
)

func priority(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func ParseToPostfix(expression string) []string {
	var numbers []string
	var operators []rune
	var currentNum strings.Builder

	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			currentNum.WriteRune(char)
		} else {
			if currentNum.Len() > 0 {
				numbers = append(numbers, currentNum.String())
				currentNum.Reset()
			}
			if char == '(' {
				operators = append(operators, '(')
			} else if char == ')' {
				for len(operators) > 0 && operators[len(operators)-1] != '(' {
					numbers = append(numbers, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				if len(operators) > 0 && operators[len(operators)-1] == '(' {
					operators = operators[:len(operators)-1]
				}
			} else if checker.IsOperator(char) {
				op := char
				for len(operators) > 0 && priority(operators[len(operators)-1]) >= priority(op) {
					numbers = append(numbers, string(operators[len(operators)-1]))
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, op)
			}
		}
	}
	if currentNum.Len() > 0 {
		numbers = append(numbers, currentNum.String())
	}

	for len(operators) > 0 {
		numbers = append(numbers, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return numbers
}
