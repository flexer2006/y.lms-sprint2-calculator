package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*server.Server, *mux.Router) {
	cfg := &configs.ServerConfig{
		Port:              "8080",
		TimeAdditionMS:    100,
		TimeSubtractionMS: 100,
		TimeMultiplyMS:    200,
		TimeDivisionMS:    200,
	}

	log, err := logger.New(logger.Options{
		Level:       logger.Debug,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: true,
	})
	require.NoError(t, err)

	srv := server.New(cfg, log)

	handler := srv.GetHandler()
	router, ok := handler.(*mux.Router)
	require.True(t, ok, "Handler is not *mux.Router type")

	return srv, router
}

func TestServer_HandleCalculate(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	tests := []struct {
		name           string
		request        models.CalculateRequest
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "valid expression",
			request: models.CalculateRequest{
				Expression: "2 + 2",
			},
			expectedStatus: http.StatusCreated,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp models.CalculateResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.ID)
			},
		},
		{
			name: "empty expression",
			request: models.CalculateRequest{
				Expression: "",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Contains(t, resp["error"], "Expression cannot be empty")
			},
		},
		{
			name: "invalid expression",
			request: models.CalculateRequest{
				Expression: "2 + + 2",
			},
			expectedStatus: http.StatusCreated,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp models.CalculateResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.ID)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			body, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateResp(t, w)
		})
	}
}

func TestServer_HandleGetExpression(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	// Создаем выражение через API
	body, err := json.Marshal(models.CalculateRequest{Expression: "2 + 2"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var calcResp models.CalculateResponse
	err = json.NewDecoder(w.Body).Decode(&calcResp)
	require.NoError(t, err)

	tests := []struct {
		name           string
		expressionID   string
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "existing expression",
			expressionID:   calcResp.ID,
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp models.ExpressionResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Equal(t, "2 + 2", resp.Expression.Expression)
			},
		},
		{
			name:           "non-existent expression",
			expressionID:   "non-existent-id",
			expectedStatus: http.StatusNotFound,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err)
				assert.Contains(t, resp["error"], "Expression not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions/"+tt.expressionID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.validateResp(t, w)
		})
	}
}

func TestServer_HandleListExpressions(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	// Создаем несколько выражений
	expressions := []string{"2 + 2", "3 * 4", "10 - 5"}
	for _, expr := range expressions {
		body, err := json.Marshal(models.CalculateRequest{Expression: expr})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// Тестируем получение списка выражений
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.ExpressionsResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Len(t, resp.Expressions, len(expressions))
}

func TestServer_HandleGetTask(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	// Создаем выражение, которое создаст задачу
	body, err := json.Marshal(models.CalculateRequest{Expression: "2 + 2"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Тестируем получение задачи
	req = httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var resp models.TaskResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Task.ID)
		assert.Equal(t, "+", resp.Task.Operation)
	} else {
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

func TestServer_HandleSubmitTaskResult(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	// Создаем выражение и получаем задачу
	body, err := json.Marshal(models.CalculateRequest{Expression: "2 + 2"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var taskResp models.TaskResponse
	err = json.NewDecoder(w.Body).Decode(&taskResp)
	require.NoError(t, err)

	// Тестируем отправку результата
	result := models.TaskResult{
		ID:     taskResp.Task.ID,
		Result: 4.0,
	}
	body, err = json.Marshal(result)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/internal/task", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServer_Integration(t *testing.T) {
	t.Parallel()
	_, router := setupTestServer(t)

	// 1. Создаем выражение
	calcReq := models.CalculateRequest{Expression: "2 + 2"}
	body, err := json.Marshal(calcReq)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var calcResp models.CalculateResponse
	err = json.NewDecoder(w.Body).Decode(&calcResp)
	require.NoError(t, err)
	exprID := calcResp.ID

	// 2. Получаем задачу
	req = httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var taskResp models.TaskResponse
	err = json.NewDecoder(w.Body).Decode(&taskResp)
	require.NoError(t, err)

	// 3. Отправляем результат
	result := models.TaskResult{
		ID:     taskResp.Task.ID,
		Result: 4.0,
	}
	body, err = json.Marshal(result)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/internal/task", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// 4. Проверяем статус выражения
	time.Sleep(100 * time.Millisecond) // Даем время на обработку
	req = httptest.NewRequest(http.MethodGet, "/api/v1/expressions/"+exprID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var exprResp models.ExpressionResponse
	err = json.NewDecoder(w.Body).Decode(&exprResp)
	require.NoError(t, err)
	assert.Equal(t, models.StatusComplete, exprResp.Expression.Status)
	assert.NotNil(t, exprResp.Expression.Result)
	assert.Equal(t, 4.0, *exprResp.Expression.Result)
}
