package worker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"go.uber.org/zap"
)

// Agent представляет собой агента-вычислителя
type Agent struct {
	config     *configs.WorkerConfig
	logger     *logger.Logger
	httpClient *http.Client
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// New создает нового агента
func New(cfg *configs.WorkerConfig, log *logger.Logger) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	return &Agent{
		config: cfg,
		logger: log,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start запускает агента
func (a *Agent) Start() error {
	a.logger.Info("Starting agent",
		zap.Int("computing_power", a.config.ComputingPower),
		zap.String("orchestrator_url", a.config.OrchestratorURL))

	for i := 0; i < a.config.ComputingPower; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}

	return nil
}

// Stop останавливает агента
func (a *Agent) Stop() {
	a.cancel()
	a.wg.Wait()
	a.logger.Info("Agent stopped")
}
