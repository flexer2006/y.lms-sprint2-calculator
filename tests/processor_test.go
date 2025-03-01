package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"github.com/stretchr/testify/assert"
)

func setupTestServer2(t *testing.T) http.Handler {
	logOpts := logger.DefaultOptions()
	log, err := logger.New(logOpts)
	assert.NoError(t, err)

	config := &configs.ServerConfig{
		Port:              "8080",
		TimeAdditionMS:    1000,
		TimeSubtractionMS: 1000,
		TimeMultiplyMS:    2000,
		TimeDivisionMS:    2000,
	}

	srv := server.New(config, log)
	return srv.GetHandler()
}

func TestExpressionValidation(t *testing.T) {
	handler := setupTestServer2(t)

	tests := []struct {
		name           string
		expression     string
		expectedStatus int
		expectedError  string
	}{

		{"Simple addition", "1+2", http.StatusCreated, ""},
		{"Expression with spaces", "1 + 2", http.StatusCreated, ""},
		{"Multiple operations", "1+2*3", http.StatusCreated, ""},
		{"Decimal numbers", "1.5+2.3", http.StatusCreated, ""},
		{"Single number", "42", http.StatusUnprocessableEntity, "invalid expression: too few tokens"},
		{"Two numbers", "42 53", http.StatusUnprocessableEntity, "invalid expression: too few tokens"},
		{"Trailing operator", "1+2+", http.StatusUnprocessableEntity, "invalid expression: trailing operator"},
		{"Leading operator", "+1+2", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
		{"Invalid character", "1+{2", http.StatusUnprocessableEntity, "invalid expression: unexpected character"},
		{"Curly braces", "{1+*}", http.StatusUnprocessableEntity, "invalid expression: unexpected character"},
		{"Too few tokens with operator", "2*", http.StatusUnprocessableEntity, "invalid expression: too few tokens"},
		{"Invalid structure", "1++2", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
		{"Division by zero", "5/0", http.StatusCreated, ""}, // Division by zero обрабатывается позже

		{"Subtraction", "5-3", http.StatusCreated, ""},
		{"Multiplication", "4*2", http.StatusCreated, ""},
		{"Division", "10/2", http.StatusCreated, ""},
		{"Mixed operations with precedence", "10+5-3*2/4", http.StatusCreated, ""},
		{"Simple parentheses", "(2+3)*4", http.StatusCreated, ""},
		{"Nested parentheses", "((2+3)*4)-5", http.StatusCreated, ""},
		{"Unary minus at beginning", "-2+3", http.StatusCreated, ""},
		{"Unary minus inside parentheses", "(2+(-3))*4", http.StatusCreated, ""},
		{"Complex nested expression", "((1+2)*((3-4)+5))/6", http.StatusCreated, ""},
		{"Flexible spacing", "  3   +   4 *   2  ", http.StatusCreated, ""},
		{"High precision decimals", "1.234567890123456789+9.876543210987654321", http.StatusCreated, ""},
		{"Very large numbers", "999999999999999999+1", http.StatusCreated, ""},
		{"Very small decimals", "0.000000000000000001*1000000000000000000", http.StatusCreated, ""},
		{"Combination of all operators", "1+2-3*4/5", http.StatusCreated, ""},
		{"Unary minus on entire expression", "-(1+2)", http.StatusCreated, ""},
		{"Multiple unary minus inside parentheses", "(-1+(2*(3-(-4))))", http.StatusCreated, ""},
		{"Redundant parentheses", "((1+2))", http.StatusCreated, ""},

		{"Double decimal point", "1.2.3+4", http.StatusUnprocessableEntity, "invalid expression: invalid number format"},
		{"Only operator", "+", http.StatusUnprocessableEntity, "invalid expression: too few tokens"},
		{"Missing operand in parentheses", "(1+)", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
		{"Unmatched opening parenthesis", "(1+2", http.StatusUnprocessableEntity, "invalid expression: unmatched parentheses"},
		{"Extra closing parenthesis", "1+2)", http.StatusUnprocessableEntity, "invalid expression: unmatched parentheses"},
		{"Operator without left operand", "*1+2", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
		{"Consecutive operators", "1+-+2", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
		{"Empty parentheses", "()", http.StatusUnprocessableEntity, "invalid expression: empty expression"},
		{"Multiple unary minus", "--1+2", http.StatusUnprocessableEntity, "invalid expression: invalid structure"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reqBody, _ := json.Marshal(models.CalculateRequest{Expression: tc.expression})
			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			var respBody map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &respBody)
			assert.NoError(t, err)

			if tc.expectedStatus != http.StatusCreated {
				errorMsg, ok := respBody["error"].(string)
				assert.True(t, ok)
				assert.Contains(t, errorMsg, tc.expectedError)
			} else {
				_, hasID := respBody["id"]
				assert.True(t, hasID, "Response should contain expression ID")
			}
		})
	}
}
