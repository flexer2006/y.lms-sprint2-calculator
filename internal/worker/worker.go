package worker

import (
	"fmt"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"go.uber.org/zap"
)

// worker represents a computation goroutine.
func (a *Agent) worker(id int) {
	defer a.wg.Done()

	a.logger.Info("Starting worker", zap.Int(common.FieldWorkerID, id))

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Info("Worker stopped", zap.Int(common.FieldWorkerID, id))
			return
		default:
			if err := a.processTask(id); err != nil {
				a.logger.Error("Failed to process task",
					zap.Int(common.FieldWorkerID, id),
					zap.Error(err))
				time.Sleep(time.Second)
			}
		}
	}
}

// processTask processes a single task.
func (a *Agent) processTask(workerID int) error {

	task, err := a.getTask()
	if err != nil {
		return fmt.Errorf(common.ErrFormatWithWrap, "failed to get task", err)
	}

	if task == nil {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	a.logger.Info("Processing task",
		zap.Int(common.FieldWorkerID, workerID),
		zap.String(common.FieldTaskID, task.ID),
		zap.String(common.FieldOperation, task.Operation))

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	result := a.Calculate(task)

	if err := a.sendResult(task.ID, result); err != nil {
		return fmt.Errorf(common.ErrFormatWithWrap, common.LogFailedSendResult, err)
	}

	return nil
}
