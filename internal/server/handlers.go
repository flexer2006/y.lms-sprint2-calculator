// Package server предоставляет HTTP-обработчики для управления выражениями и задачами.
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

// / handleCalculate обрабатывает запрос на вычисление и запускает обработку выражений.
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

	go func() {
		if err := s.processExpression(expr); err != nil {
			s.logger.Error("Failed to process expression",
				zap.String("id", expr.ID),
				zap.String(common.FieldExpression, expr.Expression),
				zap.Error(err))

			if updateErr := s.storage.UpdateExpressionError(expr.ID, err.Error()); updateErr != nil {
				s.logger.Error("Failed to update expression error status",
					zap.String("id", expr.ID),
					zap.Error(updateErr))
			}
		}
	}()

	s.writeJSON(w, http.StatusCreated, models.CalculateResponse{ID: expr.ID})
}

// handleListExpressions список всех сохраненных выражений.
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

// handleGetExpression извлекает конкретное выражение по идентификатору.
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

// handleGetTask извлекает следующую доступную задачу.
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

// handleSubmitTaskResult обрабатывает результат выполненного задания.
func (s *Server) handleSubmitTaskResult(w http.ResponseWriter, r *http.Request) {
	var result models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		s.logger.Error(common.LogFailedDecodeTask, zap.Error(err))
		s.writeError(w, http.StatusUnprocessableEntity, common.ErrInvalidRequestBody)
		return
	}

	if err := s.storage.UpdateTaskResult(result.ID, result.Result); err != nil {
		s.logger.Error(common.LogFailedUpdateTask, zap.String(common.FieldTaskID, result.ID), zap.Error(err))
		s.writeError(w, http.StatusNotFound, common.ErrTaskNotFound)
		return
	}

	task, err := s.storage.GetTask(result.ID)
	if err != nil {
		s.logger.Error(common.LogFailedGetTaskResult, zap.String(common.FieldTaskID, result.ID), zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, common.ErrFailedProcessResult)
		return
	}

	dependentTasks := s.storage.GetTasksByDependency(result.ID)
	for _, depTask := range dependentTasks {
		depTaskCopy := *depTask
		allDepsMet := true
		for _, depID := range depTask.DependsOnTaskIDs {
			depResult, err := s.storage.GetTaskResult(depID)
			if err != nil {
				allDepsMet = false
				break
			}
			if depTask.Arg1 == 0 {
				depTaskCopy.Arg1 = depResult
			} else if depTask.Arg2 == 0 {
				depTaskCopy.Arg2 = depResult
			}
		}
		if allDepsMet && depTaskCopy.Arg1 != 0 && depTaskCopy.Arg2 != 0 {
			if err := s.storage.SaveTask(&depTaskCopy); err != nil {
				s.logger.Error("Failed to update dependent task",
					zap.String(common.FieldTaskID, depTaskCopy.ID),
					zap.String(common.FieldExpressionID, depTaskCopy.ExpressionID),
					zap.Error(err))

				// Optionally update the parent expression with an error status
				if updateErr := s.storage.UpdateExpressionError(task.ExpressionID,
					"Failed to update dependent task: "+err.Error()); updateErr != nil {
					s.logger.Error("Failed to update expression error status",
						zap.String(common.FieldExpressionID, task.ExpressionID),
						zap.Error(updateErr))
				}
			}
		}
	}

	allTasks := s.storage.GetTasksByExpressionID(task.ExpressionID)
	allCompleted := true
	for _, t := range allTasks {
		if _, err := s.storage.GetTaskResult(t.ID); err != nil {
			allCompleted = false
			break
		}
	}
	if allCompleted {
		if err := s.storage.UpdateExpressionResult(task.ExpressionID, result.Result); err != nil {
			s.logger.Error(common.LogFailedUpdateExpr, zap.String(common.FieldExpressionID, task.ExpressionID), zap.Error(err))
		}
	}

	s.logger.Info(common.LogTaskProcessed,
		zap.String(common.FieldTaskID, task.ID),
		zap.String(common.FieldExpressionID, task.ExpressionID),
		zap.Float64(common.FieldResult, result.Result))

	w.WriteHeader(http.StatusOK)
}
