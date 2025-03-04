// Package server обеспечивает настройку и управление HTTP-сервером.
package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/flexer2006/y.lms-sprint2-calculator/common"
	"github.com/flexer2006/y.lms-sprint2-calculator/configs"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/logger"
	"github.com/flexer2006/y.lms-sprint2-calculator/internal/server/storage"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server представляет собой HTTP-сервер с его конфигурацией, хранилищем и регистратором.
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

	web := router.PathPrefix("/web").Subrouter()
	web.HandleFunc("/calculate", s.handleWebCalculatePage)
	web.HandleFunc("/expressions", s.handleWebExpressionsPage)
	web.HandleFunc("/expressions/{id}", s.handleWebExpressionDetailPage)

	fs := http.FileServer(http.Dir("./web"))
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", fs))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/calculate", http.StatusMovedPermanently)
	})

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

func (s *Server) handleWebCalculatePage(w http.ResponseWriter, r *http.Request) {
	serveHTMLFile(w, r, "./web/html/calculate.html")
}

func (s *Server) handleWebExpressionsPage(w http.ResponseWriter, r *http.Request) {
	serveHTMLFile(w, r, "./web/html/expressions.html")
}

func (s *Server) handleWebExpressionDetailPage(w http.ResponseWriter, r *http.Request) {
	serveHTMLFile(w, r, "./web/html/expression-detail.html")
}

// Вспомогательная функция для обслуживания HTML-файлов
func serveHTMLFile(w http.ResponseWriter, r *http.Request, path string) {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, path)
}
