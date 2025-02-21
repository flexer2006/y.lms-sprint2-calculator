package worker

import (
	"github.com/flexer2006/y.lms-sprint2-calculator/common"
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
			a.logger.Error(common.ErrDivisionByZero,
				zap.String(common.FieldTaskID, task.ID))
			panic(common.ErrDivisionByZero)
		}
		return task.Arg1 / task.Arg2
	default:
		a.logger.Error(common.ErrUnexpectedToken,
			zap.String(common.FieldTaskID, task.ID),
			zap.String(common.FieldOperation, task.Operation))
		panic(common.ErrUnexpectedToken)
	}
}
