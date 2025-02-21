// Package storage provides data storage functionalities for expressions and tasks.
package storage

import (
	"sync"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

// Storage manages the storage of expressions and tasks.
type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Changed to a slice to ensure FIFO order
	mu          sync.Mutex
	logger      *zap.Logger
}

// New creates a new Storage instance with the provided logger.
func New(logger *zap.Logger) *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
		logger:    logger,
	}
}
