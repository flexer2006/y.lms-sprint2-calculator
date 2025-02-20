package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	config  *configs.ServerConfig
	storage *storage.Storage
	logger  *logger.Logger
	server  *http.Server
}

func New(cfg *configs.ServerConfig, log *logger.Logger) *Server {
	s := &Server{
		config:  cfg,
		storage: storage.New(),
		logger:  log,
	}

	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/calculate", s.handleCalculate).Methods(http.MethodPost)
	api.HandleFunc("/expressions", s.handleListExpressions).Methods(http.MethodGet)
	api.HandleFunc("/expressions/{id}", s.handleGetExpression).Methods(http.MethodGet)

	internal := router.PathPrefix("/internal").Subrouter()
	internal.HandleFunc("/task", s.handleGetTask).Methods(http.MethodGet)
	internal.HandleFunc("/task", s.handleSubmitTaskResult).Methods(http.MethodPost)

	s.server = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("port", s.config.Port))
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Expression) == "" {
		s.writeError(w, http.StatusUnprocessableEntity, "Expression cannot be empty")
		return
	}

	expr := &models.Expression{
		ID:         uuid.New().String(),
		Expression: req.Expression,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.storage.SaveExpression(expr); err != nil {
		s.logger.Error("Failed to save expression", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to process expression")
		return
	}

	go s.processExpression(expr)

	s.writeJSON(w, http.StatusCreated, models.CalculateResponse{ID: expr.ID})
}

func (s *Server) handleListExpressions(w http.ResponseWriter, r *http.Request) {
	expressions := s.storage.ListExpressions()
	s.writeJSON(w, http.StatusOK, models.ExpressionsResponse{Expressions: expressions})
}

func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expr, err := s.storage.GetExpression(id)
	if err != nil {
		s.writeError(w, http.StatusNotFound, "Expression not found")
		return
	}

	s.writeJSON(w, http.StatusOK, models.ExpressionResponse{Expression: *expr})
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := s.storage.GetNextTask()
	if err != nil {
		s.writeError(w, http.StatusNotFound, "No tasks available")
		return
	}

	s.writeJSON(w, http.StatusOK, models.TaskResponse{Task: *task})
}

func (s *Server) handleSubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		s.writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	if err := s.storage.UpdateTaskResult(result.ID, result.Result); err != nil {
		s.writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	task, err := s.storage.GetTask(result.ID)
	if err != nil {
		s.logger.Error("Failed to get task", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to process result")
		return
	}

	if err := s.storage.UpdateExpressionResult(task.ExpressionID, result.Result); err != nil {
		s.logger.Error("Failed to update expression result",
			zap.String("expressionID", task.ExpressionID),
			zap.Error(err))
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) processExpression(expr *models.Expression) {
	s.logger.Info("Processing expression", zap.String("id", expr.ID), zap.String("expression", expr.Expression))

	if err := s.storage.UpdateExpressionStatus(expr.ID, models.StatusProgress); err != nil {
		s.logger.Error("Failed to update expression status",
			zap.String("id", expr.ID),
			zap.Error(err))
		return
	}

	tokens, err := s.parseExpression(expr.Expression)
	if err != nil {
		s.logger.Error("Failed to parse expression",
			zap.String("id", expr.ID),
			zap.String("expression", expr.Expression),
			zap.Error(err))

		if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
			s.logger.Error("Failed to update expression error",
				zap.String("id", expr.ID),
				zap.Error(updateErr))
		}
		return
	}

	tasks := s.createTasks(expr.ID, tokens)
	if len(tasks) == 0 {
		s.logger.Error("No valid tasks created",
			zap.String("id", expr.ID),
			zap.String("expression", expr.Expression))
		if updateErr := s.storage.UpdateExpressionError(expr.ID, "Failed to create valid tasks"); updateErr != nil {
			s.logger.Error("Failed to update expression error",
				zap.String("id", expr.ID),
				zap.Error(updateErr))
		}
		return
	}

	for _, task := range tasks {
		if err := s.storage.SaveTask(task); err != nil {
			s.logger.Error("Failed to save task",
				zap.String("expressionID", expr.ID),
				zap.String("taskID", task.ID),
				zap.Error(err))
			if updateErr := s.storage.UpdateExpressionError(expr.ID, "Failed to create tasks"); updateErr != nil {
				s.logger.Error("Failed to update expression error",
					zap.String("id", expr.ID),
					zap.Error(updateErr))
			}
			return
		}
	}

	s.logger.Info("Tasks created successfully",
		zap.String("id", expr.ID),
		zap.Int("taskCount", len(tasks)))
}

func (s *Server) parseExpression(expression string) ([]string, error) {

	expression = strings.ReplaceAll(expression, " ", "")

	var tokens []string
	var currentNumber strings.Builder

	for i := 0; i < len(expression); i++ {
		char := expression[i]

		if char == '-' && (i == 0 || isOperator(string(expression[i-1]))) {
			currentNumber.WriteRune('-')
			continue
		}

		if isDigit(char) || char == '.' {
			currentNumber.WriteByte(char)

			if i == len(expression)-1 {
				tokens = append(tokens, currentNumber.String())
			}
			continue
		}

		if currentNumber.Len() > 0 {
			tokens = append(tokens, currentNumber.String())
			currentNumber.Reset()
		}

		if isOperator(string(char)) {
			tokens = append(tokens, string(char))
		}
	}

	if len(tokens) < 3 {
		return nil, fmt.Errorf("invalid expression: too few tokens")
	}

	return tokens, nil
}

func (s *Server) createTasks(exprID string, tokens []string) []*models.Task {
	var tasks []*models.Task

	for i := 0; i < len(tokens); i++ {
		if isOperator(tokens[i]) {
			operationTime := s.getOperationTime(tokens[i])

			task := &models.Task{
				ID:            uuid.New().String(),
				ExpressionID:  exprID,
				Operation:     tokens[i],
				OperationTime: operationTime,
			}

			if i > 0 && i < len(tokens)-1 {
				arg1, err1 := strconv.ParseFloat(tokens[i-1], 64)
				arg2, err2 := strconv.ParseFloat(tokens[i+1], 64)
				if err1 == nil && err2 == nil {
					task.Arg1 = arg1
					task.Arg2 = arg2
					tasks = append(tasks, task)
				}
			}
		}
	}

	return tasks
}

func (s *Server) getOperationTime(op string) int64 {
	switch op {
	case "+":
		return s.config.TimeAdditionMS
	case "-":
		return s.config.TimeSubtractionMS
	case "*":
		return s.config.TimeMultiplyMS
	case "/":
		return s.config.TimeDivisionMS
	default:
		return 100
	}
}

func isOperator(token string) bool {
	switch token {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s.logger.Error("Failed to write JSON response", zap.Error(err))
	}
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		s.logger.Error("Failed to write error response", zap.Error(err))
	}
}

// Добавляем метод для получения HTTP handler
func (s *Server) GetHandler() http.Handler {
	return s.server.Handler
}
