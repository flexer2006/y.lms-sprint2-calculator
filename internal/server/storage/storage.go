package storage

import (
	"sync"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/models"

	"go.uber.org/zap"
)

type Storage struct {
	expressions sync.Map
	tasks       sync.Map
	taskQueue   []models.Task // Изменяем на слайс для гарантированного FIFO порядка
	mu          sync.Mutex
	logger      *zap.Logger
}

func New(logger *zap.Logger) *Storage {
	return &Storage{
		taskQueue: make([]models.Task, 0),
		logger:    logger,
	}
}
