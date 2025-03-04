// Package main is the entry point for the orchestrator application.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server"

	"go.uber.org/zap"
)

// main инициализирует регистратор, конфигурацию и запускает сервер.
// Он также обрабатывает graceful shutdown при получении сигналов о завершении работы.
func main() {

	opts := logger.DefaultOptions()
	opts.LogDir = "logs/orchestrator"

	log, err := logger.New(opts)
	if err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, common.ErrFailedInitLogger, err)
		if printErr != nil {
			_, writeErr := fmt.Fprintln(os.Stderr, "Failed to write to stderr:", printErr)
			if writeErr != nil {
				os.Exit(2)
			}
			os.Exit(2)
		}
		os.Exit(1)
	}
	defer func() {
		if syncErr := log.Sync(); syncErr != nil {
			_, printErr := fmt.Fprintf(os.Stderr, common.ErrFailedSyncLogger, syncErr)
			if printErr != nil {
				_, writeErr := fmt.Fprintln(os.Stderr, "Failed to write to stderr:", printErr)
				if writeErr != nil {
					os.Exit(2)
				}
				os.Exit(2)
			}
		}
	}()

	cfg, err := configs.NewServerConfig()
	if err != nil {
		log.Fatal(common.ErrFailedInitConfig, zap.Error(err))
	}

	srv := server.New(cfg, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal(common.ErrFailedStartServer, zap.Error(err))
		}
	}()

	log.Info(common.LogOrchestratorStarted)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(common.ErrServerShutdownFailed, zap.Error(err))
	}

	log.Info(common.LogOrchestratorStoppedGrace)
}
