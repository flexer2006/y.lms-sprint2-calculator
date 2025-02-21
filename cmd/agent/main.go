package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/worker"

	"go.uber.org/zap"
)

func main() {
	// Инициализируем логгер
	opts := logger.DefaultOptions()
	opts.LogDir = "logs/agent"

	log, err := logger.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() // Используем метод Sync() вместо logger.Close()

	// Создаем контекст с отменой
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инициализируем конфигурацию агента
	cfg, err := configs.NewWorkerConfig()
	if err != nil {
		log.Fatal("Failed to initialize config", zap.Error(err))
	}

	// Создаем и запускаем агента
	agent := worker.New(cfg, log)
	if err := agent.Start(); err != nil {
		log.Fatal("Failed to start agent", zap.Error(err))
	}

	log.Info("Agent service started successfully")

	// Ожидаем сигнала завершения
	<-ctx.Done()

	agent.Stop()
	log.Info("Agent service stopped gracefully")
}
