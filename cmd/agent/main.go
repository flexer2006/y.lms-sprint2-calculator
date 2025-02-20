package main

import (
	"net/http"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"go.uber.org/zap"
)

func main() {
	logger.Info("Starting Agent...")

	// Start the agent's worker routines
	startWorkers()

	logger.Info("Agent is running")
	select {} // Keep the agent running
}

func startWorkers() {
	// Logic to start worker goroutines
	// Each worker will request tasks from the orchestrator and process them
}

func requestTask() {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		logger.Error("Failed to request task", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// Process the task
}

func submitResult(taskID int, result float64) {
	// Logic to submit the result back to the orchestrator
}
