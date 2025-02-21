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
				assert.Contains(t, resp["error"], "Invalid request body")
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

	_, router := setupTestServer(t)

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

	_, router := setupTestServer(t)

	expressions := []string{"2 + 2", "3 * 4", "10 - 5"}
	for _, expr := range expressions {
		body, err := json.Marshal(models.CalculateRequest{Expression: expr})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

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

	body, err := json.Marshal(models.CalculateRequest{Expression: "2 + 2"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

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

	_, router := setupTestServer(t)

	body, err := json.Marshal(models.CalculateRequest{Expression: "2 + 2"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var taskResp models.TaskResponse
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("timeout waiting for task to be available")
		case <-ticker.C:
			req = httptest.NewRequest(http.MethodGet, "/internal/task", nil)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				err = json.NewDecoder(w.Body).Decode(&taskResp)
				require.NoError(t, err)
				goto taskFound
			}
		}
	}

taskFound:

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
}

func TestServer_Integration(t *testing.T) {
	_, router := setupTestServer(t)

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

	time.Sleep(100 * time.Millisecond)

	req = httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var taskResp models.TaskResponse
	err = json.NewDecoder(w.Body).Decode(&taskResp)
	require.NoError(t, err)

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

	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("timeout waiting for expression to complete")
		case <-ticker.C:
			req = httptest.NewRequest(http.MethodGet, "/api/v1/expressions/"+exprID, nil)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				continue
			}

			var exprResp models.ExpressionResponse
			err = json.NewDecoder(w.Body).Decode(&exprResp)
			require.NoError(t, err)

			if exprResp.Expression.Status == models.StatusComplete {
				assert.NotNil(t, exprResp.Expression.Result)
				assert.Equal(t, 4.0, *exprResp.Expression.Result)
				return
			}
		}
	}
}
