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
	opts.LogDir = "logs/agent"

	log, err := logger.New(opts)
	if err != nil {

		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Starting agent service...")

	if err := initializeAgent(); err != nil {
		log.Fatal("Failed to initialize agent",
			zap.Error(err),
			zap.String("service", "agent"),
			zap.Time("startup_time", time.Now()),
		)
	}

	<-ctx.Done()
	log.Info("Shutting down agent service gracefully")
}

func initializeAgent() error {

	return nil
}
