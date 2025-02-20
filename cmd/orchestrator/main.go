package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"go.uber.org/zap"
)

func main() {

	opts := logger.DefaultOptions()
	opts.LogDir = "logs/orchestrator"

	log, err := logger.New(opts)
	if err != nil {

		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Starting orchestrator service...")

	if err := initializeOrchestrator(); err != nil {
		log.Fatal("Failed to initialize orchestrator",
			zap.Error(err),
			zap.String("service", "orchestrator"),
			zap.Time("startup_time", time.Now()),
		)
	}

	<-ctx.Done()
	log.Info("Shutting down orchestrator service gracefully")
}

func initializeOrchestrator() error {

	return nil
}
