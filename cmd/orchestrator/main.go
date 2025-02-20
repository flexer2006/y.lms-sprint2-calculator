package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server"

	"go.uber.org/zap"
)

func main() {
	// Инициализируем логгер
	opts := logger.DefaultOptions()
	opts.LogDir = "logs/orchestrator"

	log, err := logger.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() // Используем метод Sync() вместо logger.Close()

	// Загружаем конфигурацию сервера
	cfg, err := configs.NewServerConfig()
	if err != nil {
		log.Fatal("Failed to initialize config", zap.Error(err))
	}

	// Создаем и запускаем сервер
	srv := server.New(cfg, log)

	// Создаем контекст с отменой
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	log.Info("Orchestrator service started successfully")

	// Ожидаем сигнала завершения
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown failed", zap.Error(err))
	}

	log.Info("Orchestrator service stopped gracefully")
}
