package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

// getTask retrieves a task from the orchestrator.
func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.httpClient.Get(fmt.Sprintf(common.PathInternalTask, a.config.OrchestratorURL))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error(common.ErrFailedCloseRespBody, zap.Error(err))
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(common.ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return &taskResp.Task, nil
}

// sendResult sends the calculation result to the orchestrator.
func (a *Agent) sendResult(taskID string, result float64) error {
	taskResult := models.TaskResult{
		ID:     taskID,
		Result: result,
	}

	body, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Post(
		fmt.Sprintf(common.PathInternalTask, a.config.OrchestratorURL),
		common.ContentTypeJSON,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error(common.ErrFailedCloseRespBody, zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(common.ErrUnexpectedStatusCode, resp.StatusCode)
	}

	return nil
}
