// Package calculation provides functions to tokenize mathematical expressions.
package calculation

import (
	"strconv"
	"strings"
)

// tokenize splits an expression string into tokens.
func tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder
	var lastWasNumber bool

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		switch char {
		case ' ', '\t':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
				lastWasNumber = true
			}
			continue
		case '+', '-', '*', '/', '%', '^', '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
				lastWasNumber = true
			}
			if char == '-' {
				if i == 0 || expression[i-1] == '(' || isOperator(string(expression[i-1])) {
					tokens = append(tokens, "-")
					continue
				}
			}
			if lastWasNumber && char == '(' {
				return nil
			}
			tokens = append(tokens, string(char))
			lastWasNumber = false
		default:
			if lastWasNumber && number.Len() == 0 {
				return nil
			}
			if char == '.' {
				if strings.Contains(number.String(), ".") {
					return nil
				}
			}
			if !isDigit(char) && char != '.' {
				return nil
			}
			number.WriteRune(char)
			lastWasNumber = false
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens
}

// isOperator checks if a token is a valid operator.
func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/", "%", "^":
		return true
	}
	return false
}

// isNumber checks if a string represents a valid number.
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isDigit checks if a rune is a digit.
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
