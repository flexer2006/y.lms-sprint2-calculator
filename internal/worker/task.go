package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
)

// getTask получает задачу от оркестратора
func (a *Agent) getTask() (*models.Task, error) {
	resp, err := a.httpClient.Get(fmt.Sprintf("%s/internal/task", a.config.OrchestratorURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Нет доступных задач
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taskResp models.TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}

	return &taskResp.Task, nil
}

// sendResult отправляет результат вычисления оркестратору
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
		fmt.Sprintf("%s/internal/task", a.config.OrchestratorURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}