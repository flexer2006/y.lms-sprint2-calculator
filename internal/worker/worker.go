package worker

import (
	"fmt"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"go.uber.org/zap"
)

// worker представляет собой горутину-вычислителя
func (a *Agent) worker(id int) {
	defer a.wg.Done()

	a.logger.Info(common.LogWorkerStarting, zap.Int(common.FieldWorkerID, id))

	for {
		select {
		case <-a.ctx.Done():
			a.logger.Info(common.LogWorkerStopped, zap.Int(common.FieldWorkerID, id))
			return
		default:
			if err := a.processTask(id); err != nil {
				a.logger.Error(common.LogFailedProcessTask,
					zap.Int(common.FieldWorkerID, id),
					zap.Error(err))
				time.Sleep(time.Second)
			}
		}
	}
}

// processTask обрабатывает одну задачу
func (a *Agent) processTask(workerID int) error {
	// Получаем задачу от оркестратора
	task, err := a.getTask()
	if err != nil {
		return fmt.Errorf(common.ErrFormatWithWrap, common.LogFailedGetTask, err)
	}

	// Если задач нет, ждем немного
	if task == nil {
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	a.logger.Info(common.LogProcessingTask,
		zap.Int(common.FieldWorkerID, workerID),
		zap.String(common.FieldTaskID, task.ID),
		zap.String(common.FieldOperation, task.Operation))

	// Имитируем время выполнения операции
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	// Вычисляем результат
	result := a.Calculate(task)

	// Отправляем результат
	if err := a.sendResult(task.ID, result); err != nil {
		return fmt.Errorf(common.ErrFormatWithWrap, common.LogFailedSendResult, err)
	}

	return nil
}
