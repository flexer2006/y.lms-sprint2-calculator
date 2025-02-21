package test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestStorage_SaveExpression(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	tests := []struct {
		name    string
		expr    *models.Expression
		wantErr bool
	}{
		{
			name: "valid expression",
			expr: &models.Expression{
				ID:         "test-id-1",
				Expression: "2+2",
				Status:     models.StatusPending,
			},
			wantErr: false,
		},
		{
			name: "empty id",
			expr: &models.Expression{
				Expression: "2+2",
				Status:     models.StatusPending,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := store.SaveExpression(tt.expr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify saved expression
			saved, err := store.GetExpression(tt.expr.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.expr.Expression, saved.Expression)
			assert.Equal(t, tt.expr.Status, saved.Status)
			assert.False(t, saved.CreatedAt.IsZero())
			assert.False(t, saved.UpdatedAt.IsZero())
		})
	}
}

func TestStorage_GetExpression(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	// Test non-existent expression
	_, err := store.GetExpression("non-existent")
	assert.Error(t, err)

	// Save and get expression
	expr := &models.Expression{
		ID:         "test-id-1",
		Expression: "2+2",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(expr))

	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, expr.Expression, saved.Expression)
}

func TestStorage_ListExpressions(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	// Initial list should be empty
	expressions := store.ListExpressions()
	assert.Empty(t, expressions)

	// Add some expressions
	exprs := []*models.Expression{
		{
			ID:         "test-id-1",
			Expression: "2+2",
			Status:     models.StatusPending,
		},
		{
			ID:         "test-id-2",
			Expression: "3*3",
			Status:     models.StatusProgress,
		},
	}

	for _, expr := range exprs {
		require.NoError(t, store.SaveExpression(expr))
	}

	// Check list
	list := store.ListExpressions()
	assert.Len(t, list, len(exprs))
}

func TestStorage_SaveAndGetTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	task := &models.Task{
		ID:            "task-1",
		Arg1:          2.0,
		Arg2:          3.0,
		Operation:     "+",
		OperationTime: 100,
		ExpressionID:  "expr-1",
	}

	// Test save
	err := store.SaveTask(task)
	require.NoError(t, err)

	// Test get
	saved, err := store.GetTask(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.ID, saved.ID)
	assert.Equal(t, task.Arg1, saved.Arg1)
	assert.Equal(t, task.Arg2, saved.Arg2)
	assert.Equal(t, task.Operation, saved.Operation)
	assert.Equal(t, task.OperationTime, saved.OperationTime)
	assert.Equal(t, task.ExpressionID, saved.ExpressionID)
}

func TestStorage_UpdateTaskResult(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	// Try to update non-existent task
	err := store.UpdateTaskResult("non-existent", 42.0)
	assert.Error(t, err)

	// Create and update task
	task := &models.Task{
		ID:            "task-1",
		Arg1:          2.0,
		Arg2:          3.0,
		Operation:     "+",
		OperationTime: 100,
		ExpressionID:  "expr-1",
	}
	require.NoError(t, store.SaveTask(task))

	result := 5.0
	require.NoError(t, store.UpdateTaskResult(task.ID, result))

	// Verify update
	saved, err := store.GetTask(task.ID)
	require.NoError(t, err)
	assert.NotNil(t, saved.Result)
	assert.Equal(t, result, *saved.Result)
}

func TestStorage_GetNextTask(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	// Initially should return error
	_, err := store.GetNextTask()
	assert.Error(t, err)

	// Add tasks
	tasks := []*models.Task{
		{
			ID:            "task-1",
			Arg1:          2.0,
			Arg2:          3.0,
			Operation:     "+",
			OperationTime: 100,
			ExpressionID:  "expr-1",
		},
		{
			ID:            "task-2",
			Arg1:          4.0,
			Arg2:          5.0,
			Operation:     "*",
			OperationTime: 200,
			ExpressionID:  "expr-1",
		},
	}

	for _, task := range tasks {
		require.NoError(t, store.SaveTask(task))
	}

	// Get tasks in order
	for _, expected := range tasks {
		task, getErr := store.GetNextTask()
		require.NoError(t, getErr)
		assert.Equal(t, expected.ID, task.ID)
		assert.Equal(t, expected.Operation, task.Operation)
	}

	_, err = store.GetNextTask()
	assert.Error(t, err)
}

func TestStorage_UpdateExpressionStatus(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	err := store.UpdateExpressionStatus("non-existent", models.StatusProgress)
	assert.Error(t, err)

	expr := &models.Expression{
		ID:         "test-id-1",
		Expression: "2+2",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(expr))

	newStatus := models.StatusProgress
	require.NoError(t, store.UpdateExpressionStatus(expr.ID, newStatus))

	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, newStatus, saved.Status)
	assert.True(t, saved.UpdatedAt.After(saved.CreatedAt))
}

func TestStorage_UpdateExpressionResult(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	err := store.UpdateExpressionResult("non-existent", 42.0)
	assert.Error(t, err)

	expr := &models.Expression{
		ID:         "test-id-1",
		Expression: "2+2",
		Status:     models.StatusProgress,
	}
	require.NoError(t, store.SaveExpression(expr))

	result := 4.0
	require.NoError(t, store.UpdateExpressionResult(expr.ID, result))

	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.NotNil(t, saved.Result)
	assert.Equal(t, result, *saved.Result)
	assert.Equal(t, models.StatusComplete, saved.Status)
}

func TestStorage_UpdateExpressionError(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	err := store.UpdateExpressionError("non-existent", "error message")
	assert.Error(t, err)

	expr := &models.Expression{
		ID:         "test-id-1",
		Expression: "2+2",
		Status:     models.StatusProgress,
	}
	require.NoError(t, store.SaveExpression(expr))

	errMsg := "test error"
	require.NoError(t, store.UpdateExpressionError(expr.ID, errMsg))

	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, errMsg, saved.Error)
	assert.Equal(t, models.StatusError, saved.Status)
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)
	done := make(chan bool)
	const goroutines = 10

	for i := 0; i < goroutines; i++ {
		go func(id string) {
			expr := &models.Expression{
				ID:         id,
				Expression: "2+2",
				Status:     models.StatusPending,
			}
			_ = store.SaveExpression(expr)
			_ = store.UpdateExpressionStatus(id, models.StatusProgress)
			_ = store.UpdateExpressionResult(id, 4.0)
			done <- true
		}(fmt.Sprintf("expr-%d", i))
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

	expressions := store.ListExpressions()
	assert.Len(t, expressions, goroutines)
}

func TestStorage_SaveTask_Validation(t *testing.T) {
	t.Parallel()
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	tests := []struct {
		name    string
		task    *models.Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &models.Task{
				ID:            "task-1",
				Arg1:          2.0,
				Arg2:          3.0,
				Operation:     "+",
				OperationTime: 100,
				ExpressionID:  "expr-1",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			task: &models.Task{
				Arg1:          2.0,
				Arg2:          3.0,
				Operation:     "+",
				OperationTime: 100,
				ExpressionID:  "expr-1",
			},
			wantErr: true,
		},
		{
			name: "invalid operation",
			task: &models.Task{
				ID:            "task-2",
				Arg1:          2.0,
				Arg2:          3.0,
				Operation:     "%",
				OperationTime: 100,
				ExpressionID:  "expr-1",
			},
			wantErr: false,
		},
		{
			name: "zero operation time",
			task: &models.Task{
				ID:            "task-3",
				Arg1:          2.0,
				Arg2:          3.0,
				Operation:     "+",
				OperationTime: 0,
				ExpressionID:  "expr-1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := store.SaveTask(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			saved, err := store.GetTask(tt.task.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.task.Operation, saved.Operation)
			assert.Equal(t, tt.task.OperationTime, saved.OperationTime)
		})
	}
}

func TestStorage_TaskQueue_Behavior(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	tasks := []*models.Task{
		{
			ID:            "task-1",
			Arg1:          2.0,
			Arg2:          3.0,
			Operation:     "+",
			OperationTime: 100,
			ExpressionID:  "expr-1",
		},
		{
			ID:            "task-2",
			Arg1:          4.0,
			Arg2:          5.0,
			Operation:     "*",
			OperationTime: 200,
			ExpressionID:  "expr-1",
		},
	}

	for i := 0; i < len(tasks); i++ {
		require.NoError(t, store.SaveTask(tasks[i]))
	}

	for i := 0; i < len(tasks); i++ {
		task, err := store.GetNextTask()
		require.NoError(t, err)
		assert.Equal(t, tasks[i].ID, task.ID, "Tasks should be returned in FIFO order")
	}

	_, err := store.GetNextTask()
	assert.Error(t, err, "Queue should be empty after processing all tasks")

	for _, expectedTask := range tasks {

		require.NoError(t, store.SaveTask(expectedTask))

		task, err := store.GetNextTask()
		require.NoError(t, err)
		assert.Equal(t, expectedTask.ID, task.ID, "Task should be retrieved in FIFO order")
	}
}

func TestStorage_ExpressionLifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	expr := &models.Expression{
		ID:         "test-lifecycle",
		Expression: "2+2*3",
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	require.NoError(t, store.SaveExpression(expr))

	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusPending, saved.Status)
	assert.Nil(t, saved.Result)

	require.NoError(t, store.UpdateExpressionStatus(expr.ID, models.StatusProgress))
	saved, err = store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusProgress, saved.Status)

	errMsg := "division by zero"
	require.NoError(t, store.UpdateExpressionError(expr.ID, errMsg))
	saved, err = store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusError, saved.Status)
	assert.Equal(t, errMsg, saved.Error)

	result := 42.0
	require.NoError(t, store.UpdateExpressionResult(expr.ID, result))
	saved, err = store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusComplete, saved.Status)
	assert.Equal(t, &result, saved.Result)
}

func TestStorage_ConcurrentTaskProcessing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)
	done := make(chan bool)
	const workers = 5
	const tasksPerWorker = 10

	for i := 0; i < workers*tasksPerWorker; i++ {
		task := &models.Task{
			ID:            fmt.Sprintf("task-%d", i),
			Arg1:          float64(i),
			Arg2:          float64(i + 1),
			Operation:     "+",
			OperationTime: 100,
			ExpressionID:  "expr-1",
		}
		require.NoError(t, store.SaveTask(task))
	}

	processedTasks := make(map[string]bool)
	var mu sync.Mutex

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for {
				task, err := store.GetNextTask()
				if err != nil {

					break
				}

				time.Sleep(time.Millisecond)
				result := task.Arg1 + task.Arg2
				require.NoError(t, store.UpdateTaskResult(task.ID, result))

				mu.Lock()
				processedTasks[task.ID] = true
				mu.Unlock()
			}
			done <- true
		}(i)
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	mu.Lock()
	assert.Equal(t, workers*tasksPerWorker, len(processedTasks))
	mu.Unlock()
}

func TestStorage_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	store := storage.New(logger)

	expr := &models.Expression{
		ID:         "large-numbers",
		Expression: "1e308 + 1e308",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(expr))

	expr = &models.Expression{
		ID:         "small-numbers",
		Expression: "1e-308 * 1e-308",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(expr))

	longExpr := &models.Expression{
		ID:         "long-expression",
		Expression: strings.Repeat("1+", 1000) + "1",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(longExpr))

	expr = &models.Expression{
		ID:         "concurrent-updates",
		Expression: "1+1",
		Status:     models.StatusPending,
	}
	require.NoError(t, store.SaveExpression(expr))

	// Test concurrent updates with valid status transitions
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// First update: PENDING -> IN_PROGRESS
			_ = store.UpdateExpressionStatus(expr.ID, models.StatusProgress)
		}()
	}
	wg.Wait()

	// Verify intermediate status
	saved, err := store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusProgress, saved.Status)

	// Final update: IN_PROGRESS -> COMPLETE
	require.NoError(t, store.UpdateExpressionResult(expr.ID, 2.0))

	// Verify final status
	saved, err = store.GetExpression(expr.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusComplete, saved.Status)
}
