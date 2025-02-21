package calculation

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"go.uber.org/zap"
)

// Parser represents a mathematical expression parser
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
		return 0, errors.New(common.ErrUnexpectedToken)
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
				return 0, errors.New(common.ErrDivisionByZero)
			}
			left /= right
		case "%":
			if right == 0 {
				return 0, errors.New(common.ErrModuloByZero)
			}
			if left != float64(int(left)) || right != float64(int(right)) {
				return 0, errors.New(common.ErrInvalidModulo)
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
			logger.Error(common.LogUnexpectedEndExpr,
				zap.Strings(common.FieldTokens, p.tokens),
				zap.Int(common.FieldPosition, p.pos))
		}
		return 0, errors.New(common.ErrUnexpectedEndExpr)
	}

	token := p.tokens[p.pos]
	p.pos++

	switch {
	case token == "(":
		result, err := p.parseExpression()
		if err != nil {
			if logger != nil {
				logger.Error(common.LogFailedParseParentheses,
					zap.Error(err),
					zap.Strings(common.FieldTokens, p.tokens),
					zap.Int(common.FieldPosition, p.pos))
			}
			return 0, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			if logger != nil {
				logger.Error(common.LogMissingCloseParen,
					zap.Strings(common.FieldTokens, p.tokens),
					zap.Int(common.FieldPosition, p.pos))
			}
			return 0, errors.New(common.ErrMissingCloseParen)
		}
		p.pos++
		return result, nil
	case token == "-":
		factor, err := p.parseFactor()
		if err != nil {
			if logger != nil {
				logger.Error(common.LogFailedParseNegative,
					zap.Error(err),
					zap.Strings(common.FieldTokens, p.tokens),
					zap.Int(common.FieldPosition, p.pos))
			}
			return 0, err
		}
		return -factor, nil
	case isNumber(token):
		num, err := strconv.ParseFloat(token, 64)
		if err != nil {
			if logger != nil {
				logger.Error(common.LogInvalidNumberFormat,
					zap.String(common.FieldToken, token),
					zap.Error(err))
			}
			return 0, fmt.Errorf("invalid number: %s", token)
		}
		return num, nil
	default:
		if logger != nil {
			logger.Error(common.LogUnexpectedToken,
				zap.String(common.FieldToken, token),
				zap.Strings(common.FieldTokens, p.tokens),
				zap.Int(common.FieldPosition, p.pos))
		}
		return 0, fmt.Errorf("unexpected token: %s", token)
	}
}
