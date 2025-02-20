package calculation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var logger *zap.Logger

// InitLogger initializes the logger for the calculation package
func InitLogger(l *zap.Logger) {
	logger = l
}

func EvaluateExpression(expression string) (float64, error) {
	if expression == "" {
		return 0, errors.New("expression is empty")
	}

	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("invalid expression")
	}

	parser := &Parser{tokens: tokens, pos: 0}
	return parser.parse()
}

type Parser struct {
	tokens []string
	pos    int
}

func (p *Parser) parse() (float64, error) {
	result, err := p.parseExpression()
	if err != nil {
		return 0, err
	}
	if p.pos < len(p.tokens) {
		return 0, errors.New("unexpected token")
	}
	return result, nil
}

func (p *Parser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "+" && op != "-" {
			break
		}
		p.pos++

		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		if op == "+" {
			left += right
		} else {
			left -= right
		}
	}

	return left, nil
}

func (p *Parser) parseTerm() (float64, error) {
	left, err := p.parsePower()
	if err != nil {
		return 0, err
	}

	for p.pos < len(p.tokens) {
		op := p.tokens[p.pos]
		if op != "*" && op != "/" && op != "%" {
			break
		}
		p.pos++

		right, err := p.parsePower()
		if err != nil {
			return 0, err
		}

		switch op {
		case "*":
			left *= right
		case "/":
			if right == 0 {
				return 0, errors.New("division by zero")
			}
			left /= right
		case "%":
			if right == 0 {
				return 0, errors.New("modulo by zero")
			}
			if left != float64(int(left)) || right != float64(int(right)) {
				return 0, errors.New("modulo operation requires integer operands")
			}
			left = math.Mod(left, right)
		}
	}

	return left, nil
}

func (p *Parser) parsePower() (float64, error) {
	result, err := p.parseFactor()
	if err != nil {
		return 0, err
	}

	if p.pos < len(p.tokens) && p.tokens[p.pos] == "^" {
		p.pos++

		exponent, err := p.parsePower()
		if err != nil {
			return 0, err
		}
		result = math.Pow(result, exponent)
	}

	return result, nil
}

func (p *Parser) parseFactor() (float64, error) {
	if p.pos >= len(p.tokens) {
		if logger != nil {
			logger.Error("Unexpected end of expression", 
				zap.Strings("tokens", p.tokens),
				zap.Int("position", p.pos))
		}
		return 0, errors.New("unexpected end of expression")
	}

	token := p.tokens[p.pos]
	p.pos++

	switch {
	case token == "(":
		result, err := p.parseExpression()
		if err != nil {
			if logger != nil {
				logger.Error("Failed to parse expression in parentheses",
					zap.Error(err),
					zap.Strings("tokens", p.tokens),
					zap.Int("position", p.pos))
			}
			return 0, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			if logger != nil {
				logger.Error("Missing closing parenthesis",
					zap.Strings("tokens", p.tokens),
					zap.Int("position", p.pos))
			}
			return 0, errors.New("missing closing parenthesis")
		}
		p.pos++
		return result, nil
	case token == "-":
		factor, err := p.parseFactor()
		if err != nil {
			if logger != nil {
				logger.Error("Failed to parse negative factor",
					zap.Error(err),
					zap.Strings("tokens", p.tokens),
					zap.Int("position", p.pos))
			}
			return 0, err
		}
		return -factor, nil
	case isNumber(token):
		num, err := strconv.ParseFloat(token, 64)
		if err != nil {
			if logger != nil {
				logger.Error("Invalid number format",
					zap.String("token", token),
					zap.Error(err))
			}
			return 0, fmt.Errorf("invalid number: %s", token)
		}
		return num, nil
	default:
		if logger != nil {
			logger.Error("Unexpected token",
				zap.String("token", token),
				zap.Strings("tokens", p.tokens),
				zap.Int("position", p.pos))
		}
		return 0, fmt.Errorf("unexpected token: %s", token)
	}
}

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

// Helper functions
func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/", "%", "^":
		return true
	}
	return false
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
