// Package server provides the HTTP server setup and management.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/storage"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server represents the HTTP server with its configuration, storage, and logger.
type Server struct {
	config  *configs.ServerConfig
	storage *storage.Storage
	logger  *logger.Logger
	server  *http.Server
}

// New creates a new Server instance with the provided configuration and logger.
func New(cfg *configs.ServerConfig, log *logger.Logger) *Server {
	s := &Server{
		config:  cfg,
		storage: storage.New(log.Logger),
		logger:  log,
	}

	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/calculate", s.handleCalculate).Methods(http.MethodPost)
	api.HandleFunc("/expressions", s.handleListExpressions).Methods(http.MethodGet)
	api.HandleFunc("/expressions/{id}", s.handleGetExpression).Methods(http.MethodGet)

	internal := router.PathPrefix("/internal").Subrouter()
	internal.HandleFunc(common.PathTask, s.handleGetTask).Methods(http.MethodGet)
	internal.HandleFunc(common.PathTask, s.handleSubmitTaskResult).Methods(http.MethodPost)

	s.server = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("Server initialized",
		zap.String(common.FieldPort, cfg.Port),
		zap.Int64("timeAdditionMS", cfg.TimeAdditionMS),
		zap.Int64("timeSubtractionMS", cfg.TimeSubtractionMS),
		zap.Int64("timeMultiplyMS", cfg.TimeMultiplyMS),
		zap.Int64("timeDivisionMS", cfg.TimeDivisionMS))

	return s
}

// GetHandler returns the HTTP handler for the server.
func (s *Server) GetHandler() http.Handler {
	return s.server.Handler
}

// Start begins listening on the configured port and serves HTTP requests.
func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String(common.FieldPort, s.config.Port))
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
