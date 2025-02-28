// Package server provides HTTP handlers for managing expressions and tasks.
package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// / handleCalculate processes a calculation request and initiates expression processing.
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("Failed to decode request body",
			zap.Error(err))
		s.writeError(w, http.StatusUnprocessableEntity, common.ErrInvalidRequestBody)
		return
	}

	if req.Expression == "" {
		s.logger.Warn("Empty expression received")
		s.writeError(w, http.StatusUnprocessableEntity, common.ErrInvalidRequestBody)
		return
	}

	_, err := s.parseExpression(req.Expression)
	if err != nil {
		s.logger.Error(common.LogFailedParseExpression,
			zap.String(common.FieldExpression, req.Expression),
			zap.Error(err))

		s.writeError(w, http.StatusUnprocessableEntity, err.Error())
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
			zap.String(common.FieldExpression, req.Expression),
			zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, common.ErrFailedProcessExpression)
		return
	}

	s.logger.Info("Expression received for calculation",
		zap.String("id", expr.ID),
		zap.String(common.FieldExpression, expr.Expression))

	go s.processExpression(expr)

	s.writeJSON(w, http.StatusCreated, models.CalculateResponse{ID: expr.ID})
}

// handleListExpressions lists all stored expressions.
func (s *Server) handleListExpressions(w http.ResponseWriter, _ *http.Request) {
	exprPointers := s.storage.ListExpressions()
	expressions := make([]models.Expression, len(exprPointers))
	for i, expr := range exprPointers {
		expressions[i] = *expr
	}
	s.logger.Debug("Listing all expressions",
		zap.Int(common.FieldCount, len(expressions)))
	s.writeJSON(w, http.StatusOK, models.ExpressionsResponse{Expressions: expressions})
}

// handleGetExpression retrieves a specific expression by ID.
func (s *Server) handleGetExpression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expr, err := s.storage.GetExpression(id)
	if err != nil {
		s.logger.Warn(common.LogExpressionRetrieved,
			zap.String("id", id))
		s.writeError(w, http.StatusNotFound, common.ErrExpressionNotFound)
		return
	}

	s.logger.Debug(common.LogExpressionRetrieved,
		zap.String("id", id),
		zap.String(common.FieldStatus, string(expr.Status)))
	s.writeJSON(w, http.StatusOK, models.ExpressionResponse{Expression: *expr})
}

// handleGetTask retrieves the next available task.
func (s *Server) handleGetTask(w http.ResponseWriter, _ *http.Request) {
	task, err := s.storage.GetNextTask()
	if err != nil {
		s.logger.Debug(common.LogNoTasksAvailable)
		s.writeError(w, http.StatusNotFound, common.ErrTaskNotFound)
		return
	}

	s.logger.Debug(common.LogTaskRetrieved,
		zap.String(common.FieldTaskID, task.ID),
		zap.String(common.FieldOperation, task.Operation))
	s.writeJSON(w, http.StatusOK, models.TaskResponse{Task: *task})
}

// handleSubmitTaskResult processes the result of a completed task.
func (s *Server) handleSubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		s.logger.Error(common.LogFailedDecodeTask,
			zap.Error(err))
		s.writeError(w, http.StatusUnprocessableEntity, common.ErrInvalidRequestBody)
		return
	}

	if err := s.storage.UpdateTaskResult(result.ID, result.Result); err != nil {
		s.logger.Error(common.LogFailedUpdateTask,
			zap.String(common.FieldTaskID, result.ID),
			zap.Error(err))
		s.writeError(w, http.StatusNotFound, common.ErrTaskNotFound)
		return
	}

	task, err := s.storage.GetTask(result.ID)
	if err != nil {
		s.logger.Error(common.LogFailedGetTaskResult,
			zap.String(common.FieldTaskID, result.ID),
			zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, common.ErrFailedProcessResult)
		return
	}

	if err := s.storage.UpdateExpressionResult(task.ExpressionID, result.Result); err != nil {
		s.logger.Error(common.LogFailedUpdateExpr,
			zap.String(common.FieldExpressionID, task.ExpressionID),
			zap.String(common.FieldTaskID, task.ID),
			zap.Float64(common.FieldResult, result.Result),
			zap.Error(err))
	}

	s.logger.Info(common.LogTaskProcessed,
		zap.String(common.FieldTaskID, task.ID),
		zap.String(common.FieldExpressionID, task.ExpressionID),
		zap.Float64(common.FieldResult, result.Result))

	w.WriteHeader(http.StatusOK)
}
