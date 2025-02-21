package worker

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"

	"go.uber.org/zap"
)

// Agent represents a computation agent.
type Agent struct {
	config     *configs.WorkerConfig
	logger     *logger.Logger
	httpClient *http.Client
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// New creates a new agent.
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

// Start launches the agent.
func (a *Agent) Start() error {
	a.logger.Info("Starting agent",
		zap.Int(common.FieldComputingPower, a.config.ComputingPower),
		zap.String(common.FieldOrchestratorURL, a.config.OrchestratorURL))

	for i := 0; i < a.config.ComputingPower; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}

	return nil
}

// Stop stops the agent.
func (a *Agent) Stop() {
	a.cancel()
	a.wg.Wait()
	a.logger.Info("Agent stopped")
}
