package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Failed to decode request body",
			zap.Error(err))
		s.writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	if req.Expression == "" {
		s.logger.Warn("Empty expression received")
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
		s.logger.Error("Failed to save expression",
			zap.String("expression", req.Expression),
			zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to process expression")
		return
	}

	s.logger.Info("Expression received for calculation",
		zap.String("id", expr.ID),
		zap.String("expression", expr.Expression))

	go s.processExpression(expr)

	s.writeJSON(w, http.StatusCreated, models.CalculateResponse{ID: expr.ID})
}

func (s *Server) handleListExpressions(w http.ResponseWriter, r *http.Request) {
	exprPointers := s.storage.ListExpressions()
	expressions := make([]models.Expression, len(exprPointers))
	for i, expr := range exprPointers {
		expressions[i] = *expr
	}
	s.logger.Debug("Listing all expressions",
		zap.Int("count", len(expressions)))
	s.writeJSON(w, http.StatusOK, models.ExpressionsResponse{Expressions: expressions})
}

func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expr, err := s.storage.GetExpression(id)
	if err != nil {
		s.logger.Warn("Expression not found",
			zap.String("id", id))
		s.writeError(w, http.StatusNotFound, "Expression not found")
		return
	}

	s.logger.Debug("Expression retrieved",
		zap.String("id", id),
		zap.String("status", string(expr.Status)))
	s.writeJSON(w, http.StatusOK, models.ExpressionResponse{Expression: *expr})
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := s.storage.GetNextTask()
	if err != nil {
		s.logger.Debug("No tasks available")
		s.writeError(w, http.StatusNotFound, "No tasks available")
		return
	}

	s.logger.Debug("Task retrieved",
		zap.String("taskID", task.ID),
		zap.String("operation", task.Operation))
	s.writeJSON(w, http.StatusOK, models.TaskResponse{Task: *task})
}

func (s *Server) handleSubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		s.logger.Error("Failed to decode task result",
			zap.Error(err))
		s.writeError(w, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	if err := s.storage.UpdateTaskResult(result.ID, result.Result); err != nil {
		s.logger.Error("Failed to update task result",
			zap.String("taskID", result.ID),
			zap.Error(err))
		s.writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	task, err := s.storage.GetTask(result.ID)
	if err != nil {
		s.logger.Error("Failed to get task after updating result",
			zap.String("taskID", result.ID),
			zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to process result")
		return
	}

	if err := s.storage.UpdateExpressionResult(task.ExpressionID, result.Result); err != nil {
		s.logger.Error("Failed to update expression result",
			zap.String("expressionID", task.ExpressionID),
			zap.String("taskID", task.ID),
			zap.Float64("result", result.Result),
			zap.Error(err))
	}

	s.logger.Info("Task result processed successfully",
		zap.String("taskID", task.ID),
		zap.String("expressionID", task.ExpressionID),
		zap.Float64("result", result.Result))

	w.WriteHeader(http.StatusOK)
}