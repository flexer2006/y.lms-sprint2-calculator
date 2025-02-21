// Package main is the entry point for the agent application.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/worker"

	"go.uber.org/zap"
)

// main initializes the logger, configuration, and starts the agent.
// It also handles graceful shutdown on receiving termination signals.
func main() {

	opts := logger.DefaultOptions()
	opts.LogDir = "logs/agent"

	log, err := logger.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, common.ErrFailedInitLogger, err)
		os.Exit(1)
	}
	defer func() {
		if syncErr := log.Sync(); syncErr != nil {
			fmt.Fprintf(os.Stderr, common.ErrFailedSyncLogger, syncErr)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := configs.NewWorkerConfig()
	if err != nil {
		log.Fatal(common.ErrFailedInitConfig, zap.Error(err))
	}

	agent := worker.New(cfg, log)
	if err := agent.Start(); err != nil {
		log.Fatal(common.ErrFailedStartAgent, zap.Error(err))
	}

	log.Info(common.LogAgentStarted)

	<-ctx.Done()

	agent.Stop()
	log.Info(common.LogAgentStoppedGrace)
}
