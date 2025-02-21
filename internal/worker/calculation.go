package worker

import (
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

// Calculate выполняет вычисление
func (a *Agent) Calculate(task *models.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			a.logger.Error("Division by zero",
				zap.String("task_id", task.ID))
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		a.logger.Error("Unknown operation",
			zap.String("task_id", task.ID),
			zap.String("operation", task.Operation))
		return 0
	}
}
